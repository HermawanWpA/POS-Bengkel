package repository

import (
	"context" // Pastikan package context ini sudah di-import
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type pelangganRepository struct {
	db *gorm.DB
}

func NewPelangganRepository(db *gorm.DB) domain.PelangganRepository {
	return &pelangganRepository{db: db}
}

// 1. Tambahkan parameter ctx context.Context
func (r *pelangganRepository) Create(ctx context.Context, pelanggan *domain.Pelanggan) error {
	// Gunakan .WithContext(ctx) sebelum .Create agar query aman dan mendukung timeout/cancellation
	return r.db.WithContext(ctx).Create(pelanggan).Error
}

// 2. Tambahkan parameter ctx context.Context juga di sini
// func (r *pelangganRepository) FetchWithVehicles(ctx context.Context) ([]domain.Pelanggan, error) {
// 	var pelangganList []domain.Pelanggan

// 	// Gunakan .WithContext(ctx) sebelum .Preload
// 	err := r.db.WithContext(ctx).Preload("Kendaraan").Find(&pelangganList).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pelangganList, nil
// }

func (r *pelangganRepository) GetByID(ctx context.Context, id int) (domain.Pelanggan, error) {
	var pelanggan domain.Pelanggan

	// Mencari satu data berdasarkan ID primary key dan mengikutkan data kendaraannya
	err := r.db.WithContext(ctx).
		Preload("Kendaraan").
		First(&pelanggan, id).Error

	return pelanggan, err
}

func (r *pelangganRepository) Delete(ctx context.Context, id int) error {
	var pelanggan domain.Pelanggan

	// 1. Cari data pelanggannya terlebih dahulu
	err := r.db.WithContext(ctx).First(&pelanggan, id).Error
	if err != nil {
		return err // Ini akan otomatis mengembalikan error "record not found" jika ID tidak ada
	}

	// 2. Hapus data pelanggan tersebut (GORM akan otomatis menghapus kendaraan yang terikat)
	return r.db.WithContext(ctx).Delete(&pelanggan).Error
}
func (r *pelangganRepository) Update(ctx context.Context, pelanggan *domain.Pelanggan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingPelanggan domain.Pelanggan

		// 1. Validasi: Cari tahu apakah pelanggan dengan ID ini benar-benar ada
		if err := tx.First(&existingPelanggan, pelanggan.ID).Error; err != nil {
			return err // Jika tidak ada, ini akan otomatis melempar error "record not found"
		}

		// 2. Jika ada, update data utama profil pelanggan
		if err := tx.Model(&existingPelanggan).Updates(map[string]interface{}{
			"nama_pelanggan": pelanggan.NamaPelanggan,
			"no_hp":          pelanggan.NoHp,
			"alamat":         pelanggan.Alamat,
		}).Error; err != nil {
			return err
		}

		// 3. Sinkronisasi Data Kendaraan jika disertakan di JSON
		if len(pelanggan.Kendaraan) > 0 {
			for _, k := range pelanggan.Kendaraan {
				k.IdPelanggan = pelanggan.ID

				// tx.Save akan melakukan UPDATE jika no_polisi sudah ada, atau INSERT jika belum ada
				if err := tx.Save(&k).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *pelangganRepository) Search(ctx context.Context, keyword string) ([]domain.Pelanggan, error) {
	var listPelanggan []domain.Pelanggan

	// Mencari data yang cocok dengan nama_pelanggan ATAU no_hp
	err := r.db.WithContext(ctx).
		Preload("Kendaraan").
		Where("nama_pelanggan LIKE ?", "%"+keyword+"%").
		Or("no_hp LIKE ?", "%"+keyword+"%").
		Find(&listPelanggan).Error

	return listPelanggan, err
}

func (r *pelangganRepository) FetchWithPagination(ctx context.Context, param domain.PaginationParam) ([]domain.Pelanggan, int64, error) {
	var listPelanggan []domain.Pelanggan
	var totalData int64

	// 1. Inisialisasi query dasar mendeteksi Model Pelanggan
	baseQuery := r.db.WithContext(ctx).Model(&domain.Pelanggan{})

	// 2. HITUNG TOTAL DATA (Sebelum diberikan Limit & Offset untuk kebutuhan Pagination frontend)
	if err := baseQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// 3. JALANKAN PAGINASI DENGAN PRELOAD REKOMENDASI (Menggunakan "Kendaraan" sesuai nama field struct Go)
	offset := (param.Page - 1) * param.Limit
	err := baseQuery.Preload("Kendaraan").Limit(param.Limit).Offset(offset).Find(&listPelanggan).Error
	if err != nil {
		return nil, 0, err
	}

	// Mengembalikan 3 parameter: slice data, total baris, dan error (jika nil)
	return listPelanggan, totalData, nil
}

func (r *pelangganRepository) FetchWithVehicles(ctx context.Context) ([]domain.Pelanggan, error) {
	var listPelanggan []domain.Pelanggan

	// Perbaikan dari Preload("Vehicles") -> Preload("Kendaraan")
	err := r.db.WithContext(ctx).Preload("Kendaraan").Find(&listPelanggan).Error
	if err != nil {
		return nil, err
	}
	return listPelanggan, nil
}
