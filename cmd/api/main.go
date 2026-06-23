package main

import (
	"pos-echo-app/config"
	"pos-echo-app/internal/delivery/http"
	"pos-echo-app/internal/repository"
	"pos-echo-app/internal/usecase"

	"github.com/labstack/echo/v4"
	// ====== PASTIKAN ADA /v4/ DI TENGAH SEPERTI INI ======
	"github.com/labstack/echo/v4/middleware" // Menggunakan modul resmi Echo
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	// 1. Init Koneksi Database GORM
	db := config.ConnectDB()

	// ==========================================
	// INISIALISASI REPOSITORY & USECASE (REPOT & UU)
	// ==========================================

	// Inisialisasi Modul Pelanggan
	pelangganRepo := repository.NewPelangganRepository(db)
	pelangganUsecase := usecase.NewPelangganUsecase(pelangganRepo)

	// Inisialisasi Modul Sparepart
	sparepartRepo := repository.NewSparepartRepository(db)
	sparepartUU := usecase.NewSparepartUsecase(sparepartRepo)

	// Inisialisasi Modul master jasa
	jasaRepo := repository.NewJasaRepository(db)
	jasaUU := usecase.NewJasaUsecase(jasaRepo)

	// Inisialisasi Modul Mekanik
	mekanikRepo := repository.NewMekanikRepository(db)
	mekanikUU := usecase.NewMekanikUsecase(mekanikRepo)

	// Inisialisasi Modul Transaksi
	transaksiRepo := repository.NewTransaksiRepository(db)
	transaksiUU := usecase.NewTransaksiUsecase(transaksiRepo, jasaRepo, sparepartRepo)

	// Inisialisasi Modul Auth / User
	userRepo := repository.NewUserRepository(db)
	userUU := usecase.NewUserUsecase(userRepo)

	// 2. Setup Layer Dashboard (Repo -> Usecase)
	dashRepo := repository.NewDashboardRepository(db)
	dashUsecase := usecase.NewDashboardUsecase(dashRepo)

	// ==========================================
	// REKAYASA ROUTING & PEMASANGAN MIDDLEWARE
	// ==========================================

	// Jalur 1: JALUR BEBAS / PUBLIC (Tanpa Token JWT)
	// Login dan Register ditaruh di sini supaya user bisa masuk sistem terlebih dahulu

	// Jalur 2: JALUR DIKUNCI / PROTECTED (Wajib Pakai Token JWT)
	// Kita buat group rute baru dengan prefix /api/v1
	apiGroup := e.Group("/api/v1")

	// Pasang satpam middleware di pintu gerbang group ini
	apiGroup.Use(http.AuthMiddleware)

	// Daftarkan semua handler operasional ke dalam 'apiGroup' (Bukan ke 'e' lagi)
	http.NewPelangganHandler(apiGroup, pelangganUsecase)
	http.NewSparepartHandler(apiGroup, sparepartUU)
	http.NewJasaHandler(apiGroup, jasaUU)
	http.NewMekanikHandler(apiGroup, mekanikUU)
	http.NewTransaksiHandler(apiGroup, transaksiUU)
	http.NewDashboardHandler(apiGroup, dashUsecase)

	http.NewUserHandler(e, apiGroup, userUU)
	e.Logger.Fatal(e.Start(":8080"))
}
