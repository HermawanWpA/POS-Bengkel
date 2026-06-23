package repository

import (
	"context"
	"errors"
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (domain.DashboardStatsResponse, error) {
	var stats domain.DashboardStatsResponse

	// 1. Agregasi Jumlah Kendaraan Servis (PERBAIKAN: Hanya yang berstatus selesai)
	queryKendaraan := `
		SELECT 
			COUNT(CASE WHEN DATE(tanggal_transaksi) = CURDATE() THEN 1 END) as hari_ini,
			COUNT(CASE WHEN YEARWEEK(tanggal_transaksi, 1) = YEARWEEK(CURDATE(), 1) THEN 1 END) as minggu_ini,
			COUNT(CASE WHEN MONTH(tanggal_transaksi) = MONTH(CURDATE()) AND YEAR(tanggal_transaksi) = YEAR(CURDATE()) THEN 1 END) as bulan_ini,
			COUNT(CASE WHEN YEAR(tanggal_transaksi) = YEAR(CURDATE()) THEN 1 END) as tahun_ini
		FROM transaksi 
		WHERE status_pengerjaan = 'selesai'` // <-- Diubah dari IN ('proses', 'selesai') menjadi hanya 'selesai'

	if err := r.db.WithContext(ctx).Raw(queryKendaraan).Scan(&stats.JumlahKendaraanServis).Error; err != nil {
		return stats, err
	}

	// 2. Agregasi Total Laba Kotor (Pendapatan dari Transaksi Selesai)
	queryLabaKotor := `
		SELECT 
			COALESCE(SUM(CASE WHEN DATE(tanggal_transaksi) = CURDATE() THEN total_bayar END), 0) as hari_ini,
			COALESCE(SUM(CASE WHEN YEARWEEK(tanggal_transaksi, 1) = YEARWEEK(CURDATE(), 1) THEN total_bayar END), 0) as minggu_ini,
			COALESCE(SUM(CASE WHEN MONTH(tanggal_transaksi) = MONTH(CURDATE()) AND YEAR(tanggal_transaksi) = YEAR(CURDATE()) THEN total_bayar END), 0) as bulan_ini,
			COALESCE(SUM(CASE WHEN YEAR(tanggal_transaksi) = YEAR(CURDATE()) THEN total_bayar END), 0) as tahun_ini
		FROM transaksi 
		WHERE status_pengerjaan = 'selesai'`

	if err := r.db.WithContext(ctx).Raw(queryLabaKotor).Scan(&stats.LabaKotor).Error; err != nil {
		return stats, err
	}

	// 3. Agregasi Total Laba Bersih (Total Bayar - HPP Sparepart)
	queryLabaBersih := `
		SELECT 
			COALESCE(SUM(CASE WHEN DATE(t.tanggal_transaksi) = CURDATE() THEN (t.total_bayar - tx_cost.total_hpp) END), 0) as hari_ini,
			COALESCE(SUM(CASE WHEN YEARWEEK(t.tanggal_transaksi, 1) = YEARWEEK(CURDATE(), 1) THEN (t.total_bayar - tx_cost.total_hpp) END), 0) as minggu_ini,
			COALESCE(SUM(CASE WHEN MONTH(t.tanggal_transaksi) = MONTH(CURDATE()) AND YEAR(t.tanggal_transaksi) = YEAR(CURDATE()) THEN (t.total_bayar - tx_cost.total_hpp) END), 0) as bulan_ini,
			COALESCE(SUM(CASE WHEN YEAR(t.tanggal_transaksi) = YEAR(CURDATE()) THEN (t.total_bayar - tx_cost.total_hpp) END), 0) as tahun_ini
		FROM transaksi t
		JOIN (
			SELECT 
				id_transaksi,
				COALESCE((SELECT SUM(dt_sp.qty * sp.harga_beli_hpp) FROM detail_transaksi_sparepart dt_sp JOIN master_sparepart sp ON sp.kode_sparepart = dt_sp.kode_sparepart WHERE dt_sp.id_transaksi = transaksi.id_transaksi), 0) as total_hpp
			FROM transaksi
		) tx_cost ON tx_cost.id_transaksi = t.id_transaksi
		WHERE t.status_pengerjaan = 'selesai'`

	if err := r.db.WithContext(ctx).Raw(queryLabaBersih).Scan(&stats.LabaBersih).Error; err != nil {
		return stats, err
	}

	// 4. Agregasi Total Pelanggan (Berdasarkan waktu pendaftaran/pembuatan data pelanggan)
	queryPelanggan := `
		SELECT 
			COUNT(CASE WHEN DATE(created_at) = CURDATE() THEN 1 END) as hari_ini,
			COUNT(CASE WHEN YEARWEEK(created_at, 1) = YEARWEEK(CURDATE(), 1) THEN 1 END) as minggu_ini,
			COUNT(CASE WHEN MONTH(created_at) = MONTH(CURDATE()) AND YEAR(created_at) = YEAR(CURDATE()) THEN 1 END) as bulan_ini,
			COUNT(CASE WHEN YEAR(created_at) = YEAR(CURDATE()) THEN 1 END) as tahun_ini
		FROM pelanggan`

	if err := r.db.WithContext(ctx).Raw(queryPelanggan).Scan(&stats.TotalPelanggan).Error; err != nil {
		return stats, err
	}

	return stats, nil
}

func (r *dashboardRepository) GetDetailsByPeriod(ctx context.Context, period string) ([]domain.ServiceDetailRow, error) {
	var list []domain.ServiceDetailRow

	var dateFilterSQL string
	switch period {
	case "hari":
		dateFilterSQL = "DATE(tanggal_transaksi) = CURDATE() AND status_pengerjaan = 'selesai'"
	case "minggu":
		dateFilterSQL = "YEARWEEK(tanggal_transaksi, 1) = YEARWEEK(CURDATE(), 1) AND status_pengerjaan = 'selesai'"
	case "bulan":
		dateFilterSQL = "MONTH(tanggal_transaksi) = MONTH(CURDATE()) AND YEAR(tanggal_transaksi) = YEAR(CURDATE()) AND status_pengerjaan = 'selesai'"
	case "tahun":
		dateFilterSQL = "YEAR(tanggal_transaksi) = YEAR(CURDATE()) AND status_pengerjaan = 'selesai'"
	default:
		return nil, errors.New("periode filter tidak valid (pilih: hari, minggu, bulan, tahun)")
	}

	// Murni mengambil dari tabel transaksi saja tanpa JOIN ke tabel lain
	err := r.db.WithContext(ctx).Table("transaksi").
		Select("id_transaksi, " +
			"no_polisi, " +
			"'-' as nama_pelanggan, " + // Isi sementara dengan strip
			"'-' as merek_model, " + // Isi sementara dengan strip
			"'-' as nama_mekanik, " + // Isi sementara dengan strip
			"status_pengerjaan, " +
			"total_bayar, " +
			"tanggal_transaksi").
		Where(dateFilterSQL).
		Order("tanggal_transaksi DESC").
		Scan(&list).Error

	return list, err
}
