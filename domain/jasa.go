package domain

import (
	"context"
	"time"
)

type MasterJasa struct {
	IDJasa           int       `json:"id_jasa" gorm:"primaryKey;autoIncrement;column:id_jasa"`
	NamaJasa         string    `json:"nama_jasa" gorm:"type:varchar(150);not null"`
	TarifJasa        float64   `json:"tarif_jasa" gorm:"type:decimal(10,2);not null"`       // Harga yang dibayar konsumen
	PersentaseKomisi float64   `json:"persentase_komisi" gorm:"type:decimal(5,2);not null"` // Contoh: 30.00
	KomisiMekanik    float64   `json:"komisi_mekanik" gorm:"type:decimal(10,2);not null"`   // Dihitung otomatis (PersentaseKomisi/100 * TarifJasa)
	CreatedAt        time.Time `json:"created_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy        string    `json:"created_by" gorm:"type:varchar(100);not null"`
}

// Menentukan nama tabel secara eksplisit agar sesuai dengan DDL SQL Anda
func (MasterJasa) TableName() string {
	return "master_jasa"
}

type JasaRepository interface {
	Create(ctx context.Context, jasa *MasterJasa) error
	Fetch(ctx context.Context) ([]MasterJasa, error)
	GetByID(ctx context.Context, id int) (MasterJasa, error)
	Update(ctx context.Context, jasa *MasterJasa) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, keyword string) ([]MasterJasa, error)
}

type JasaUsecase interface {
	Create(ctx context.Context, jasa *MasterJasa) error
	Fetch(ctx context.Context) ([]MasterJasa, error)
	GetByID(ctx context.Context, id int) (MasterJasa, error)
	Update(ctx context.Context, jasa *MasterJasa) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, keyword string) ([]MasterJasa, error)
}
