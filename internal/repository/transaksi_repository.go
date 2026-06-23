package repository

import (
	"context"
	"errors"
	"fmt"
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepository(db *gorm.DB) domain.TransaksiRepository {
	return &transaksiRepository{db: db}
}

func (r *transaksiRepository) CreateTransactionTx(ctx context.Context, txFunc func(repo domain.TransaksiRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &transaksiRepository{db: tx}
		return txFunc(txRepo)
	})
}

func (r *transaksiRepository) Fetch(ctx context.Context) ([]domain.Transaksi, error) {
	var listTransaksi []domain.Transaksi
	err := r.db.WithContext(ctx).Order("tanggal_transaksi DESC").Find(&listTransaksi).Error
	return listTransaksi, err
}

func (r *transaksiRepository) GetByID(ctx context.Context, id string) (domain.Transaksi, []domain.DetailTransaksiJasa, []domain.DetailTransaksiSparepart, error) {
	var header domain.Transaksi
	var jasas []domain.DetailTransaksiJasa
	var spareparts []domain.DetailTransaksiSparepart

	if err := r.db.WithContext(ctx).Where("id_transaksi = ?", id).First(&header).Error; err != nil {
		return header, nil, nil, err
	}
	if err := r.db.WithContext(ctx).Where("id_transaksi = ?", id).Find(&jasas).Error; err != nil {
		return header, nil, nil, err
	}
	if err := r.db.WithContext(ctx).Where("id_transaksi = ?", id).Find(&spareparts).Error; err != nil {
		return header, nil, nil, err
	}
	return header, jasas, spareparts, nil
}

// ====== FUNGSI UPDATE DATA REALTIME DI DATABASE ======
func (r *transaksiRepository) UpdateStatus(ctx context.Context, id string, statusPengerjaan string, metodeBayar string) error {
	return r.db.WithContext(ctx).Model(&domain.Transaksi{}).
		Where("id_transaksi = ?", id).
		Updates(map[string]interface{}{
			"status_pengerjaan": statusPengerjaan,
			"metode_pembayaran": metodeBayar,
		}).Error
}

func (r *transaksiRepository) GetLastIDByDate(ctx context.Context, dateStr string) (string, error) {
	var lastID string
	query := "SELECT id_transaksi FROM transaksi WHERE id_transaksi LIKE ? ORDER BY id_transaksi DESC LIMIT 1"
	err := r.db.WithContext(ctx).Raw(query, "TX-"+dateStr+"-%").Row().Scan(&lastID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || lastID == "" {
			return "", nil
		}
		return "", err
	}
	return lastID, nil
}

func (r *transaksiRepository) InsertHeader(t *domain.Transaksi) error {
	return r.db.Create(t).Error
}

func (r *transaksiRepository) InsertDetailJasa(dj *domain.DetailTransaksiJasa) error {
	return r.db.Create(dj).Error
}

func (r *transaksiRepository) InsertDetailSparepart(dsp *domain.DetailTransaksiSparepart) error {
	return r.db.Create(dsp).Error
}

func (r *transaksiRepository) MinusSparepartStock(kode string, qty int) error {
	var sp domain.MasterSparepart
	if err := r.db.Where("kode_sparepart = ?", kode).First(&sp).Error; err != nil {
		return fmt.Errorf("sparepart %s tidak ditemukan", kode)
	}
	if sp.StokSekarang < qty {
		return fmt.Errorf("stok %s tidak mencukupi (Sisa: %d, Diminta: %d)", sp.NamaSparepart, sp.StokSekarang, qty)
	}
	return r.db.Model(&sp).Update("stok_sekarang", sp.StokSekarang-qty).Error
}
