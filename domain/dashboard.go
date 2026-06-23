package domain

import (
	"context"
	"time"
)

// Periode waktu untuk statistik dashboard
type PeriodStats struct {
	HariIni   float64 `json:"hari_ini"`
	MingguIni float64 `json:"minggu_ini"`
	BulanIni  float64 `json:"bulan_ini"`
	TahunIni  float64 `json:"tahun_ini"`
}

type DashboardStatsResponse struct {
	JumlahKendaraanServis PeriodStats `json:"jumlah_kendaraan_servis"`
	TotalPelanggan        PeriodStats `json:"total_pelanggan"`
	LabaKotor             PeriodStats `json:"laba_kotor"`
	LabaBersih            PeriodStats `json:"laba_bersih"`
}

// Struct untuk detail list kendaraan yang diservis
type ServiceDetailRow struct {
	IDTransaksi string `json:"id_transaksi"`
	NoPolisi    string `json:"no_polisi"`
	// NamaPelanggan    string    `json:"nama_pelanggan"`
	// MerekModel       string    `json:"merek_model"`
	// NamaMekanik      string    `json:"nama_mekanik"`
	StatusPengerjaan string    `json:"status_pengerjaan"`
	TotalBayar       float64   `json:"total_bayar"`
	TanggalTransaksi time.Time `json:"tanggal_transaksi"`
}

type DashboardRepository interface {
	GetStats(ctx context.Context) (DashboardStatsResponse, error)
	GetDetailsByPeriod(ctx context.Context, period string) ([]ServiceDetailRow, error)
}

type DashboardUsecase interface {
	GetStats(ctx context.Context) (DashboardStatsResponse, error)
	GetDetailsByPeriod(ctx context.Context, period string) ([]ServiceDetailRow, error)
}
