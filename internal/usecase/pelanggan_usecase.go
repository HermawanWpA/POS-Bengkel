package usecase

import (
	"context"
	"errors"
	"pos-echo-app/domain"
	"strings"
)

type pelangganUsecase struct {
	pelangganRepo domain.PelangganRepository
}

// NewPelangganUsecase untuk menginisialisasi usecase pelanggan
func NewPelangganUsecase(repo domain.PelangganRepository) domain.PelangganUsecase {
	return &pelangganUsecase{
		pelangganRepo: repo,
	}
}

// Create memproses bisnis logik sebelum dikirim ke repository
func (u *pelangganUsecase) Create(ctx context.Context, pelanggan *domain.Pelanggan) error {
	// 1. Validasi pointer nil untuk mencegah panic error
	if pelanggan == nil {
		return errors.New("data pelanggan tidak boleh kosong")
	}

	// 2. Validasi aturan bisnis dasar Pelanggan
	if pelanggan.NamaPelanggan == "" {
		return errors.New("nama pelanggan wajib diisi")
	}

	if len(pelanggan.NoHp) < 10 {
		return errors.New("nomor HP tidak valid, minimal 10 digit")
	}

	// =========================================================================
	// PENGEMBANGAN: Validasi nested data kendaraan (jika ada kendaraan yang dikirim)
	// =========================================================================
	if len(pelanggan.Kendaraan) > 0 {
		for _, v := range pelanggan.Kendaraan {
			if v.NoPolisi == "" {
				return errors.New("nomor polisi kendaraan tidak boleh kosong")
			}
			if v.MerekTipe == "" {
				return errors.New("merek/tipe kendaraan tidak boleh kosong")
			}
		}
	}

	// validasi created_by
	if strings.TrimSpace(pelanggan.CreatedBy) == "" {
		return errors.New("nama penginput (created_by) wajib diisi")
	}
	// 3. Lempar data ke layer repository dengan membawa ctx
	return u.pelangganRepo.Create(ctx, pelanggan)
}

// FetchWithVehicles mengambil semua data pelanggan dengan membawa ctx
func (u *pelangganUsecase) FetchWithVehicles(ctx context.Context) ([]domain.Pelanggan, error) {
	// Anda juga bisa menambahkan logika manipulasi data di sini jika dibutuhkan ke depan
	return u.pelangganRepo.FetchWithVehicles(ctx)
}
func (u *pelangganUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID pelanggan tidak valid")
	}

	// PENGEMBANGAN: Ambil data pelanggan dulu untuk memastikan datanya eksis
	list, err := u.pelangganRepo.FetchWithVehicles(ctx)
	if err != nil {
		return err
	}

	idEksis := false
	for _, p := range list {
		if p.ID == id {
			idEksis = true
			break
		}
	}

	if !idEksis {
		return errors.New("data pelanggan tidak ditemukan di database")
	}

	return u.pelangganRepo.Delete(ctx, id)
}

func (u *pelangganUsecase) GetByID(ctx context.Context, id int) (domain.Pelanggan, error) {
	// Validasi dasar: pastikan ID rasional
	if id <= 0 {
		return domain.Pelanggan{}, errors.New("ID pelanggan tidak valid")
	}

	return u.pelangganRepo.GetByID(ctx, id)
}

func (u *pelangganUsecase) Update(ctx context.Context, pelanggan *domain.Pelanggan) error {
	// 1. Validasi data kosong
	if pelanggan == nil {
		return errors.New("data pelanggan tidak boleh kosong")
	}
	if pelanggan.ID <= 0 {
		return errors.New("ID pelanggan harus disertakan untuk update")
	}
	if pelanggan.NamaPelanggan == "" {
		return errors.New("nama pelanggan wajib diisi")
	}
	if len(pelanggan.NoHp) < 10 {
		return errors.New("nomor HP tidak valid, minimal 10 digit")
	}

	// 2. Validasi data kendaraan jika ada
	if len(pelanggan.Kendaraan) > 0 {
		for _, v := range pelanggan.Kendaraan {
			if v.NoPolisi == "" {
				return errors.New("nomor polisi kendaraan tidak boleh kosong")
			}
			if v.MerekTipe == "" {
				return errors.New("merek/tipe kendaraan tidak boleh kosong")
			}
		}
	}

	// 3. Kirim ke repository
	return u.pelangganRepo.Update(ctx, pelanggan)
}

func (u *pelangganUsecase) Search(ctx context.Context, keyword string) ([]domain.Pelanggan, error) {
	// Jika kasir tidak memasukkan kata kunci apapun, kembalikan semua data
	if keyword == "" {
		return u.pelangganRepo.FetchWithVehicles(ctx)
	}

	return u.pelangganRepo.Search(ctx, keyword)
}
func (u *pelangganUsecase) GetAllWithPagination(ctx context.Context, param domain.PaginationParam) ([]domain.Pelanggan, int64, error) {
	// PEMBARUAN 1: Validasi Halaman (Page)
	if param.Page <= 0 {
		param.Page = 1
	}

	// PEMBARUAN 2: Validasi Jumlah Data per Halaman (Limit)
	if param.Limit <= 0 {
		param.Limit = 10 // Default jika tidak diisi atau minus
	} else if param.Limit > 100 {
		param.Limit = 100 // Batasi maksimal hanya 100 data per request demi keamanan server
	}

	// Teruskan parameter yang sudah divalidasi ke repository
	return u.pelangganRepo.FetchWithPagination(ctx, param)
}
