package repository

import (
	"context" // Pastikan package context ini sudah di-import
	"errors"
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type mekanikRepository struct {
	db *gorm.DB
}

func NewMekanikRepository(db *gorm.DB) domain.MekanikRepository {
	return &mekanikRepository{db: db}
}

func (r *mekanikRepository) Create(ctx context.Context, m *domain.Mekanik) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *mekanikRepository) Fetch(ctx context.Context) ([]domain.Mekanik, error) {
	var listMekanik []domain.Mekanik
	err := r.db.WithContext(ctx).Find(&listMekanik).Error
	return listMekanik, err
}

func (r *mekanikRepository) GetByID(ctx context.Context, id int) (domain.Mekanik, error) {
	var m domain.Mekanik
	err := r.db.WithContext(ctx).First(&m, id).Error
	return m, err
}

func (r *mekanikRepository) Update(ctx context.Context, m *domain.Mekanik) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingMekanik domain.Mekanik
		if err := tx.First(&existingMekanik, m.IDMekanik).Error; err != nil {
			return err
		}

		err := tx.Model(&existingMekanik).Updates(map[string]interface{}{
			"nama_mekanik": m.NamaMekanik,
			"no_hp":        m.NoHp,
			"status_aktif": m.StatusAktif,
		}).Error
		return err
	})
}

func (r *mekanikRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m domain.Mekanik
		if err := tx.First(&m, id).Error; err != nil {
			return err
		}

		// Hard delete data mekanik
		// Jika ID mekanik ini sudah terikat di tabel transaksi, MySQL otomatis menolak (Foreign Key Constraint)
		result := tx.Delete(&m)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("tidak ada data mekanik yang terhapus")
		}
		return nil
	})
}

func (r *mekanikRepository) Search(ctx context.Context, keyword string) ([]domain.Mekanik, error) {
	var listMekanik []domain.Mekanik
	err := r.db.WithContext(ctx).
		Where("nama_mekanik LIKE ?", "%"+keyword+"%").
		Find(&listMekanik).Error
	return listMekanik, err
}
