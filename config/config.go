package config

import (
	"fmt"
	"log"
	"os"
	"pos-echo-app/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ConnectDB digunakan untuk menginisialisasi koneksi ke database MySQL menggunakan GORM
func ConnectDB() *gorm.DB {
	// 1. Load file .env jika ada
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// 2. Ambil konfigurasi dari environment variable
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Jika env kosong, gunakan default value (fallback)
	if dbUser == "" {
		dbUser = "root"
	}
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "pos_bengkel1" // Tambahkan fallback nama database Anda di sini
	}

	// 3. Susun Data Source Name (DSN) untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// 4. Buka koneksi dengan NamingStrategy SingularTable
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Memaksa GORM menggunakan nama struct asli sebagai nama tabel (tanpa "s")
		},
		DisableAutomaticPing: true, // Kemarin kita tambahkan ini untuk mematikan fitur RETURNING di MySQL lama
	})
	if err != nil {
		log.Fatalf("Gagal terkoneksi ke database: %v", err)
	}

	log.Println("Berhasil terkoneksi ke database MySQL")

	// =========================================================================
	// PERBAIKAN DI SINI: Daftarkan &domain.Kendaraan{} ke dalam AutoMigrate
	// =========================================================================
	// Daftarkan semua struct domain Anda di sini dengan urutan yang logis
	err = db.AutoMigrate(
		// 1. Tabel Master Mandiri (Wajib Paling Atas)
		&domain.User{},
		&domain.Pelanggan{},
		&domain.Mekanik{},
		&domain.MasterJasa{},
		&domain.MasterSparepart{},

		// 2. Tabel yang Memiliki Foreign Key ke Tabel Master
		&domain.Kendaraan{}, // Merujuk ke Pelanggan
		&domain.Transaksi{}, // Merujuk ke Kendaraan dan Mekanik

		// 3. Tabel Detail Transaksi (Wajib Paling Bawah)
		&domain.DetailTransaksiJasa{},      // Merujuk ke Transaksi dan MasterJasa
		&domain.DetailTransaksiSparepart{}, // Merujuk ke Transaksi dan MasterSparepart

		// 4. Tabel dashboard
		&domain.PeriodStats{},      // Merujuk ke Transaksi dan MasterSparepart
		&domain.ServiceDetailRow{}, // Merujuk ke Transaksi dan MasterSparepart
	)

	if err != nil {
		log.Fatalf("Gagal melakukan AutoMigrate: %v", err)
	}

	return db
}
