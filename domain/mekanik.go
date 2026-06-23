package domain

import (
	"context"
	"time"
)


type Mekanik struct {
	IDMekanik   int       `json:"id_mekanik" gorm:"primaryKey;autoIncrement;column:id_mekanik"`
	NamaMekanik string    `json:"nama_mekanik" gorm:"type:varchar(100);not null"`
	NoHp        string    `json:"no_hp" gorm:"type:varchar(15);null"`
	StatusAktif string    `json:"status_aktif" gorm:"type:enum('aktif', 'nonaktif');default:aktif"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
}

func (Mekanik) TableName() string {
	return "mekanik"
}

// Kontrak fungsi yang harus diimplementasikan

type MekanikRepository interface {
	Create(ctx context.Context, m *Mekanik) error
	Fetch(ctx context.Context) ([]Mekanik, error)
	GetByID(ctx context.Context, id int) (Mekanik, error)
	Update(ctx context.Context, m *Mekanik) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, keyword string) ([]Mekanik, error)
}

type MekanikUsecase interface {
	Create(ctx context.Context, m *Mekanik) error
	Fetch(ctx context.Context) ([]Mekanik, error)
	GetByID(ctx context.Context, id int) (Mekanik, error)
	Update(ctx context.Context, m *Mekanik) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, keyword string) ([]Mekanik, error)
}
