package usecase

import (
	"context"
	"errors"
	"fmt"
	"pos-echo-app/domain"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type transaksiUsecase struct {
	transaksiRepo domain.TransaksiRepository
	jasaRepo      domain.JasaRepository
	sparepartRepo domain.SparepartRepository
}

func NewTransaksiUsecase(tr domain.TransaksiRepository, jr domain.JasaRepository, sr domain.SparepartRepository) domain.TransaksiUsecase {
	return &transaksiUsecase{
		transaksiRepo: tr,
		jasaRepo:      jr,
		sparepartRepo: sr,
	}
}

func (u *transaksiUsecase) UpdateStatus(ctx context.Context, id string, statusBaru string, _ string) error {
	// 1. Ambil data dari database untuk mengecek status yang ada saat ini
	headerLama, _, _, err := u.transaksiRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("NOT_FOUND: transaksi dengan ID %s tidak ditemukan", id)
		}
		return err
	}

	// 2. KUNCI UTAMA: Hanya izinkan jika status lama 'proses' DAN status baru 'selesai'
	if headerLama.StatusPengerjaan != "proses" || statusBaru != "selesai" {
		// Berikan pesan yang spesifik agar kasir/mekanik tahu aturannya
		return fmt.Errorf("BAD_REQUEST: Fitur ini hanya untuk mengubah status dari 'proses' menjadi 'selesai'. Status saat ini: '%s', Status yang Anda minta: '%s'",
			headerLama.StatusPengerjaan, statusBaru)
	}

	// 3. Jika kondisi di atas terpenuhi, eksekusi perubahan ke database via Repository
	return u.transaksiRepo.UpdateStatus(ctx, id, statusBaru, headerLama.MetodePembayaran)
}

func (u *transaksiUsecase) Create(ctx context.Context, req *domain.CreateTransaksiRequest) (string, error) {
	if req.NoPolisi == "" || req.IDMekanik <= 0 {
		return "", errors.New("nomor polisi dan mekanik wajib diisi")
	}

	if len(req.ListJasa) == 0 && len(req.ListSparepart) == 0 {
		return "", errors.New("transaksi harus memiliki minimal 1 jasa atau 1 sparepart")
	}

	var generatedID string

	// Memulai blok transaksi database yang aman (ACID)
	err := u.transaksiRepo.CreateTransactionTx(ctx, func(txRepo domain.TransaksiRepository) error {
		now := time.Now()
		dateStr := now.Format("20060102") // Format: YYYYMMDD

		// 1. Pembuatan Kode Nota Otomatis (Format: TX-20260607-0001)
		lastID, err := txRepo.GetLastIDByDate(ctx, dateStr)
		if err != nil {
			return err
		}

		counter := 1
		if lastID != "" {
			// Potong 4 digit terakhir nota lama ("0001"), ubah ke int, lalu tambah 1
			lastCounterStr := lastID[len(lastID)-4:]
			lastCounterInt, _ := strconv.Atoi(lastCounterStr)
			counter = lastCounterInt + 1
		}
		generatedID = fmt.Sprintf("TX-%s-%04d", dateStr, counter)

		var totalJasa, totalSparepart float64
		var detailJasas []domain.DetailTransaksiJasa
		var detailSpareparts []domain.DetailTransaksiSparepart

		// 2. Loop Validasi Jasa & Hitung Komisi Otomatis
		for _, jReq := range req.ListJasa {
			masterJasa, err := u.jasaRepo.GetByID(ctx, jReq.IDJasa)
			if err != nil {
				return fmt.Errorf("id jasa %d tidak valid", jReq.IDJasa)
			}

			// HITUNG OTOMATIS: Tarif * Persentase (misal: 100.000 * 0.30 = 30.000)
			nominalKomisi := masterJasa.TarifJasa * masterJasa.PersentaseKomisi

			totalJasa += masterJasa.TarifJasa
			detailJasas = append(detailJasas, domain.DetailTransaksiJasa{
				IDTransaksi:     generatedID,
				IDJasa:          masterJasa.IDJasa,
				HargaPenerapan:  masterJasa.TarifJasa,
				KomisiPenerapan: nominalKomisi, // Simpan hasil perhitungan (Snapshot)
			})
		}

		// 3. Loop Validasi Sparepart, Ambil Harga Asli, & Potong Stok
		for _, spReq := range req.ListSparepart {
			if spReq.Qty <= 0 {
				return errors.New("kuantitas sparepart harus lebih besar dari 0")
			}

			masterSp, err := u.sparepartRepo.GetByKode(ctx, spReq.KodeSparepart)
			if err != nil {
				return fmt.Errorf("kode sparepart %s tidak ditemukan", spReq.KodeSparepart)
			}

			// Eksekusi potong stok di DB aman (jika kurang, fungsi ini otomatis melempar error)
			if err := txRepo.MinusSparepartStock(masterSp.KodeSparepart, spReq.Qty); err != nil {
				return err
			}

			totalSparepart += (masterSp.HargaJual * float64(spReq.Qty))
			detailSpareparts = append(detailSpareparts, domain.DetailTransaksiSparepart{
				IDTransaksi:   generatedID,
				KodeSparepart: masterSp.KodeSparepart,
				Qty:           spReq.Qty,
				HargaSatuan:   masterSp.HargaJual,
			})
		}

		// 4. Kalkulasi Total Final Nota
		totalBayar := (totalJasa + totalSparepart) - req.Diskon
		if totalBayar < 0 {
			totalBayar = 0 // Proteksi nilai mata uang agar tidak minus
		}

		// 5. Simpan ke Database: Simpan Header Terlebih Dahulu
		header := domain.Transaksi{
			IDTransaksi:      generatedID,
			NoPolisi:         req.NoPolisi,
			IDMekanik:        req.IDMekanik,
			IDUser:           req.IDUser, // <-- TAMBAHKAN BARIS INI
			TanggalTransaksi: now,
			TotalJasa:        totalJasa,
			TotalSparepart:   totalSparepart,
			Diskon:           req.Diskon,
			TotalBayar:       totalBayar,
			MetodePembayaran: req.MetodePembayaran,
			StatusPengerjaan: "proses",
			Catatan:          req.Catatan,
		}

		if err := txRepo.InsertHeader(&header); err != nil {
			return err
		}

		// 6. Simpan Seluruh Detail Jasa ke DB
		for _, dj := range detailJasas {
			if err := txRepo.InsertDetailJasa(&dj); err != nil {
				return err
			}
		}

		// 7. Simpan Seluruh Detail Sparepart ke DB
		for _, dsp := range detailSpareparts {
			if err := txRepo.InsertDetailSparepart(&dsp); err != nil {
				return err
			}
		}

		return nil // Commit Transaksi Sukses!
	})

	if err != nil {
		return "", err
	}

	return generatedID, nil
}

func (u *transaksiUsecase) Fetch(ctx context.Context) ([]domain.Transaksi, error) {
	return u.transaksiRepo.Fetch(ctx)
}

func (u *transaksiUsecase) GetByID(ctx context.Context, id string) (domain.Transaksi, []domain.DetailTransaksiJasa, []domain.DetailTransaksiSparepart, error) {
	if id == "" {
		return domain.Transaksi{}, nil, nil, errors.New("ID transaksi tidak boleh kosong")
	}
	return u.transaksiRepo.GetByID(ctx, id)
}
