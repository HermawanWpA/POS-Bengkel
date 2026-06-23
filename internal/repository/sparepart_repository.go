package repository

import (
	"context" // Pastikan package context ini sudah di-import
	"errors"
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type sparepartRepository struct {
	db *gorm.DB
}

// NewSparepartRepository berfungsi untuk menginisialisasi repository sparepart
func NewSparepartRepository(db *gorm.DB) domain.SparepartRepository {
	return &sparepartRepository{
		db: db,
	}
}

// 1. CREATE: Menambahkan sparepart baru ke database
func (r *sparepartRepository) Create(ctx context.Context, sp *domain.MasterSparepart) error {
	return r.db.WithContext(ctx).Create(sp).Error
}

// 2. FETCH: Mengambil semua data sparepart dari database
func (r *sparepartRepository) Fetch(ctx context.Context) ([]domain.MasterSparepart, error) {
	var listSparepart []domain.MasterSparepart

	err := r.db.WithContext(ctx).Find(&listSparepart).Error
	return listSparepart, err
}

// 3. GET BY KODE: Mengambil detail satu data sparepart berdasarkan Kode SKU
func (r *sparepartRepository) GetByKode(ctx context.Context, kode string) (domain.MasterSparepart, error) {
	var sp domain.MasterSparepart

	// Menggunakan kode_sparepart sebagai pencarian utama karena bertindak sebagai Primary Key
	err := r.db.WithContext(ctx).Where("kode_sparepart = ?", kode).First(&sp).Error
	return sp, err
}

// 4. UPDATE: Memperbarui data sparepart secara selektif
func (r *sparepartRepository) Update(ctx context.Context, sp *domain.MasterSparepart) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingSp domain.MasterSparepart

		// Validasi fisik: Pastikan sparepart dengan kode tersebut memang ada di gudang
		if err := tx.Where("kode_sparepart = ?", sp.KodeSparepart).First(&existingSp).Error; err != nil {
			return err // Otomatis melempar error "record not found" jika tidak ada
		}

		// Update kolom-kolom yang diizinkan untuk diubah oleh admin/gudang
		err := tx.Model(&existingSp).Updates(map[string]interface{}{
			"nama_sparepart": sp.NamaSparepart,
			"stok_sekarang":  sp.StokSekarang,
			"stok_minimum":   sp.StokMinimum,
			"harga_beli_hpp": sp.HargaBeliHpp,
			"harga_jual":     sp.HargaJual,
			"lokasi_rak":     sp.LokasiRak,
		}).Error

		return err
	})
}

// 5. DELETE: Menghapus sparepart dari database
func (r *sparepartRepository) Delete(ctx context.Context, kode string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var sp domain.MasterSparepart

		// Validasi fisik sebelum menghapus
		if err := tx.Where("kode_sparepart = ?", kode).First(&sp).Error; err != nil {
			return err
		}

		// Eksekusi Hard Delete dari tabel master_sparepart
		// Catatan: Jika kode_sparepart sudah pernah dipakai di tabel detail_transaksi_sparepart,
		// MySQL akan otomatis menolak proses hapus ini demi menjaga integritas data keuangan (Safe Constraint).
		result := tx.Where("kode_sparepart = ?", kode).Delete(&domain.MasterSparepart{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("tidak ada data sparepart yang terhapus")
		}

		return nil
	})
}

func (r *sparepartRepository) Search(ctx context.Context, keyword string) ([]domain.MasterSparepart, error) {
	var listSparepart []domain.MasterSparepart

	// Mencari yang cocok berdasarkan kode_sparepart ATAU nama_sparepart
	err := r.db.WithContext(ctx).
		Where("kode_sparepart LIKE ?", "%"+keyword+"%").
		Or("nama_sparepart LIKE ?", "%"+keyword+"%").
		Find(&listSparepart).Error

	return listSparepart, err
}
