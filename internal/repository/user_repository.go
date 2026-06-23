package repository

import (
	"context"
	"pos-echo-app/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	return u, err
}

func (r *userRepository) GetByID(ctx context.Context, id int) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	return u, err
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingUser domain.User
		if err := tx.First(&existingUser, u.IDUser).Error; err != nil {
			return err
		}

		// Siapkan map untuk menampung field yang akan di-update
		updateData := map[string]interface{}{
			"username":  u.Username,
			"nama_user": u.NamaUser,
			"role":      u.Role,
		}

		// LOGIKA PENTING: Hanya update password jika admin mengisi kolom password baru
		if u.Password != "" {
			updateData["password"] = u.Password
		}

		return tx.Model(&existingUser).Updates(updateData).Error
	})
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}
