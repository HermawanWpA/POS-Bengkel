package domain

import (
	"context"
	"time"
)

// 1. Struct untuk Tabel Utama (Header)
type Transaksi struct {
	IDTransaksi      string    `json:"id_transaksi" gorm:"primaryKey;type:varchar(20);column:id_transaksi"`
	NoPolisi         string    `json:"no_polisi" gorm:"type:varchar(12);not null"`
	IDMekanik        int       `json:"id_mekanik" gorm:"type:int;not null"`
	IDUser           int       `json:"id_user" gorm:"type:int;not null;column:id_user"`
	TanggalTransaksi time.Time `json:"tanggal_transaksi" gorm:"default:CURRENT_TIMESTAMP"`
	TotalJasa        float64   `json:"total_jasa" gorm:"type:decimal(10,2);default:0.00"`
	TotalSparepart   float64   `json:"total_sparepart" gorm:"type:decimal(10,2);default:0.00"`
	Diskon           float64   `json:"diskon" gorm:"type:decimal(10,2);default:0.00"`
	TotalBayar       float64   `json:"total_bayar" gorm:"type:decimal(10,2);default:0.00"`
	MetodePembayaran string    `json:"metode_pembayaran" gorm:"type:enum('tunai','transfer','qris','belum_bayar');default:belum_bayar"`
	StatusPengerjaan string    `json:"status_pengerjaan" gorm:"type:enum('antri','proses','selesai');default:antri"`
	Catatan          string    `json:"catatan" gorm:"type:text;null"`
}

func (Transaksi) TableName() string { return "transaksi" }

// 2. Struct untuk Tabel Detail Jasa
type DetailTransaksiJasa struct {
	IDDetailJasa    int     `json:"id_detail_jasa" gorm:"primaryKey;autoIncrement;column:id_detail_jasa"`
	IDTransaksi     string  `json:"id_transaksi" gorm:"type:varchar(20);not null"`
	IDJasa          int     `json:"id_jasa" gorm:"type:int;not null"`
	HargaPenerapan  float64 `json:"harga_penerapan" gorm:"type:decimal(10,2);not null"`
	KomisiPenerapan float64 `json:"komisi_penerapan" gorm:"type:decimal(10,2);not null"`
}

func (DetailTransaksiJasa) TableName() string { return "detail_transaksi_jasa" }

// 3. Struct untuk Tabel Detail Sparepart
type DetailTransaksiSparepart struct {
	IDDetailSparepart int     `json:"id_detail_sparepart" gorm:"primaryKey;autoIncrement;column:id_detail_sparepart"`
	IDTransaksi       string  `json:"id_transaksi" gorm:"type:varchar(20);not null"`
	KodeSparepart     string  `json:"kode_sparepart" gorm:"type:varchar(50);not null"`
	Qty               int     `json:"qty" gorm:"type:int;not null"`
	HargaSatuan       float64 `json:"harga_satuan" gorm:"type:decimal(10,2);not null"`
	Subtotal          float64 `json:"subtotal" gorm:"->;column:subtotal"` // Read-only karena Generated Column MySQL
}

func (DetailTransaksiSparepart) TableName() string { return "detail_transaksi_sparepart" }

// ========================================================
// DTO (Data Transfer Object) UNTUK REQUEST KASIR (POSTMAN)
// ========================================================

type JasaRequest struct {
	IDJasa int `json:"id_jasa"`
}

type SparepartRequest struct {
	KodeSparepart string `json:"kode_sparepart"`
	Qty           int    `json:"qty"`
}

type CreateTransaksiRequest struct {
	NoPolisi         string             `json:"no_polisi"`
	IDMekanik        int                `json:"id_mekanik"`
	IDUser           int                `json:"id_user"` // <-- TAMBAHKAN INI
	Diskon           float64            `json:"diskon"`
	MetodePembayaran string             `json:"metode_pembayaran"`
	Catatan          string             `json:"catatan"`
	ListJasa         []JasaRequest      `json:"list_jasa"`
	ListSparepart    []SparepartRequest `json:"list_sparepart"`
}

// ========================================================
// INTERFACE CONTRACTS
// ========================================================

type TransaksiRepository interface {
	GetLastIDByDate(ctx context.Context, dateStr string) (string, error)
	CreateTransactionTx(ctx context.Context, txFunc func(repo TransaksiRepository) error) error
	InsertHeader(t *Transaksi) error
	InsertDetailJasa(dj *DetailTransaksiJasa) error
	InsertDetailSparepart(dsp *DetailTransaksiSparepart) error
	MinusSparepartStock(kode string, qty int) error
	Fetch(ctx context.Context) ([]Transaksi, error)
	GetByID(ctx context.Context, id string) (Transaksi, []DetailTransaksiJasa, []DetailTransaksiSparepart, error)
	UpdateStatus(ctx context.Context, id string, statusPengerjaan string, metodeBayar string) error
}

type TransaksiUsecase interface {
	Create(ctx context.Context, req *CreateTransaksiRequest) (string, error) // Mengembalikan ID Transaksi
	Fetch(ctx context.Context) ([]Transaksi, error)
	GetByID(ctx context.Context, id string) (Transaksi, []DetailTransaksiJasa, []DetailTransaksiSparepart, error)
	UpdateStatus(ctx context.Context, id string, statusBaru string, metodeBayar string) error
}
