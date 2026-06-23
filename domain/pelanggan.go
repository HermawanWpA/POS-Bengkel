package domain

import (
	"context"
	"time"
)

type Pelanggan struct {
	ID            int         `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaPelanggan string      `json:"nama_pelanggan" gorm:"type:varchar(100);not null"`
	NoHp          string      `json:"no_hp" gorm:"type:varchar(15);not null"`
	Alamat        string      `json:"alamat" gorm:"type:text;null"`
	CreatedAt     time.Time   `json:"created_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
	Kendaraan     []Kendaraan `json:"kendaraan,omitempty" gorm:"foreignKey:IdPelanggan;constraint:OnDelete:CASCADE;"`
	CreatedBy     string      `json:"created_by" gorm:"type:varchar(100);not null;column:created_by"`
}

type Kendaraan struct {
	NoPolisi    string    `json:"no_polisi" gorm:"primaryKey;type:varchar(15);not null"`
	IdPelanggan int       `json:"id_pelanggan" gorm:"type:int;not null;index"` // Tipe int harus sama dengan ID milik Pelanggan
	MerekTipe   string    `json:"merek_tipe" gorm:"type:varchar(100);not null"`
	NoRangka    string    `json:"no_rangka" gorm:"type:varchar(50);null"`
	NoMesin     string    `json:"no_mesin" gorm:"type:varchar(50);null"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
}

type PaginationParam struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PelangganRepository interface {
	Create(ctx context.Context, pelanggan *Pelanggan) error     // Create
	FetchWithVehicles(ctx context.Context) ([]Pelanggan, error) // Read
	Update(ctx context.Context, pelanggan *Pelanggan) error     // Update (Sudah Ada)
	Delete(ctx context.Context, id int) error                   // Delete (Sudah Ada)
	Search(ctx context.Context, keyword string) ([]Pelanggan, error)
	GetByID(ctx context.Context, id int) (Pelanggan, error)
	FetchWithPagination(ctx context.Context, param PaginationParam) ([]Pelanggan, int64, error)
}

type PelangganUsecase interface {
	Create(ctx context.Context, pelanggan *Pelanggan) error
	FetchWithVehicles(ctx context.Context) ([]Pelanggan, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, pelanggan *Pelanggan) error
	Search(ctx context.Context, keyword string) ([]Pelanggan, error)
	GetByID(ctx context.Context, id int) (Pelanggan, error)
	GetAllWithPagination(ctx context.Context, param PaginationParam) ([]Pelanggan, int64, error)
}
