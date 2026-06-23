package usecase

import (
	"context"
	"errors"
	"pos-echo-app/domain"
	"strings"
)

type sparepartUsecase struct {
	sparepartRepo domain.SparepartRepository
}

// NewSparepartUsecase berfungsi untuk menginisialisasi usecase sparepart
func NewSparepartUsecase(repo domain.SparepartRepository) domain.SparepartUsecase {
	return &sparepartUsecase{
		sparepartRepo: repo,
	}
}

// 1. LOGIKA BISNIS: TAMBAH SPAREPART
func (u *sparepartUsecase) Create(ctx context.Context, sp *domain.MasterSparepart) error {
	// Trim spasi berlebih pada kode dan nama barang
	sp.KodeSparepart = strings.TrimSpace(strings.ToUpper(sp.KodeSparepart))
	sp.NamaSparepart = strings.TrimSpace(strings.Title(strings.ToLower(sp.NamaSparepart)))

	// Validasi Bisnis: Harga jual tidak boleh di bawah harga modal (HPP) agar bengkel tidak rugi
	if sp.HargaJual < sp.HargaBeliHpp {
		return errors.New("harga jual tidak boleh lebih rendah dari harga beli (HPP)")
	}

	// Validasi Bisnis: Stok awal tidak boleh minus
	if sp.StokSekarang < 0 {
		return errors.New("stok awal tidak boleh bernilai minus")
	}

	// Cek apakah kode sparepart sudah terdaftar sebelumnya agar tidak duplikat primary key
	existingSp, _ := u.sparepartRepo.GetByKode(ctx, sp.KodeSparepart)
	if existingSp.KodeSparepart != "" {
		return errors.New("kode sparepart sudah terdaftar di sistem")
	}

	// kondisi created_by
	if strings.TrimSpace(sp.CreatedBy) == "" {
		return errors.New("nama penginput (created_by) wajib diisi")
	}

	return u.sparepartRepo.Create(ctx, sp)
}

// 2. LOGIKA BISNIS: AMBIL SEMUA DATA SPAREPART
func (u *sparepartUsecase) Fetch(ctx context.Context) ([]domain.MasterSparepart, error) {
	return u.sparepartRepo.Fetch(ctx)
}

// 3. LOGIKA BISNIS: DETAIL SPAREPART BY KODE
func (u *sparepartUsecase) GetByKode(ctx context.Context, kode string) (domain.MasterSparepart, error) {
	kode = strings.TrimSpace(strings.ToUpper(kode))
	if kode == "" {
		return domain.MasterSparepart{}, errors.New("kode sparepart kosong")
	}

	return u.sparepartRepo.GetByKode(ctx, kode)
}

// 4. LOGIKA BISNIS: UPDATE DATA SPAREPART
func (u *sparepartUsecase) Update(ctx context.Context, sp *domain.MasterSparepart) error {
	sp.KodeSparepart = strings.TrimSpace(strings.ToUpper(sp.KodeSparepart))
	sp.NamaSparepart = strings.TrimSpace(strings.Title(strings.ToLower(sp.NamaSparepart)))

	if sp.KodeSparepart == "" {
		return errors.New("kode sparepart wajib disertakan")
	}

	// Validasi Bisnis: Cek kembali margin keuntungan saat update harga
	if sp.HargaJual < sp.HargaBeliHpp {
		return errors.New("harga jual tidak boleh lebih rendah dari harga beli (HPP)")
	}

	if sp.StokSekarang < 0 {
		return errors.New("stok sekarang tidak boleh bernilai minus")
	}

	return u.sparepartRepo.Update(ctx, sp)
}

// 5. LOGIKA BISNIS: HAPUS SPAREPART
func (u *sparepartUsecase) Delete(ctx context.Context, kode string) error {
	kode = strings.TrimSpace(strings.ToUpper(kode))
	if kode == "" {
		return errors.New("kode sparepart tidak valid")
	}

	return u.sparepartRepo.Delete(ctx, kode)
}

// 6. Search by kode dan nama sparepart

func (u *sparepartUsecase) Search(ctx context.Context, keyword string) ([]domain.MasterSparepart, error) {
	keyword = strings.TrimSpace(keyword)

	// Jika kata kunci kosong, langsung tampilkan semua sparepart
	if keyword == "" {
		return u.sparepartRepo.Fetch(ctx)
	}

	return u.sparepartRepo.Search(ctx, keyword)
}
