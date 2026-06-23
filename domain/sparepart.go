package domain

import (
	"context"
	"time"
)

type MasterSparepart struct {
	KodeSparepart string    `json:"kode_sparepart" gorm:"primaryKey;type:varchar(50);not null"`
	NamaSparepart string    `json:"nama_sparepart" gorm:"type:varchar(150);not null"`
	StokSekarang  int       `json:"stok_sekarang" gorm:"type:int;default:0"`
	StokMinimum   int       `json:"stok_minimum" gorm:"type:int;default:5"`
	HargaBeliHpp  float64   `json:"harga_beli_hpp" gorm:"type:decimal(10,2);not null"`
	HargaJual     float64   `json:"harga_jual" gorm:"type:decimal(10,2);not null"`
	LokasiRak     string    `json:"lokasi_rak" gorm:"type:varchar(50);null"`
	CreatedAt     time.Time `json:"created_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy     string    `json:"created_by" gorm:"type:varchar(100);not null;column:created_by"`
}

type SparepartRepository interface {
	Create(ctx context.Context, sp *MasterSparepart) error
	Fetch(ctx context.Context) ([]MasterSparepart, error)
	GetByKode(ctx context.Context, kode string) (MasterSparepart, error)
	Update(ctx context.Context, sp *MasterSparepart) error
	Delete(ctx context.Context, kode string) error
	Search(ctx context.Context, keyword string) ([]MasterSparepart, error)
}

type SparepartUsecase interface {
	Create(ctx context.Context, sp *MasterSparepart) error
	Fetch(ctx context.Context) ([]MasterSparepart, error)
	GetByKode(ctx context.Context, kode string) (MasterSparepart, error)
	Update(ctx context.Context, sp *MasterSparepart) error
	Delete(ctx context.Context, kode string) error
	Search(ctx context.Context, keyword string) ([]MasterSparepart, error)
}
