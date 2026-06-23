package repository

import (
	"context"
	"errors"
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type jasaRepository struct {
	db *gorm.DB
}

func NewJasaRepository(db *gorm.DB) domain.JasaRepository {
	return &jasaRepository{db: db}
}

func (r *jasaRepository) Create(ctx context.Context, jasa *domain.MasterJasa) error {
	return r.db.WithContext(ctx).Create(jasa).Error
}

func (r *jasaRepository) Fetch(ctx context.Context) ([]domain.MasterJasa, error) {
	var listJasa []domain.MasterJasa
	err := r.db.WithContext(ctx).Find(&listJasa).Error
	return listJasa, err
}

func (r *jasaRepository) GetByID(ctx context.Context, id int) (domain.MasterJasa, error) {
	var jasa domain.MasterJasa
	err := r.db.WithContext(ctx).First(&jasa, id).Error
	return jasa, err
}

func (r *jasaRepository) Update(ctx context.Context, jasa *domain.MasterJasa) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingJasa domain.MasterJasa
		if err := tx.First(&existingJasa, jasa.IDJasa).Error; err != nil {
			return err
		}

		err := tx.Model(&existingJasa).Updates(map[string]interface{}{
			"nama_jasa":      jasa.NamaJasa,
			"tarif_jasa":     jasa.TarifJasa,
			"komisi_mekanik": jasa.KomisiMekanik,
		}).Error
		return err
	})
}

func (r *jasaRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var jasa domain.MasterJasa
		if err := tx.First(&jasa, id).Error; err != nil {
			return err
		}

		result := tx.Delete(&jasa)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("tidak ada data jasa yang terhapus")
		}
		return nil
	})
}

func (r *jasaRepository) Search(ctx context.Context, keyword string) ([]domain.MasterJasa, error) {
	var listJasa []domain.MasterJasa
	err := r.db.WithContext(ctx).
		Where("nama_jasa LIKE ?", "%"+keyword+"%").
		Find(&listJasa).Error
	return listJasa, err
}
