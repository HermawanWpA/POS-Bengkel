package usecase

import (
	"context"
	"errors"
	"pos-echo-app/domain"
	"strings"
)

type jasaUsecase struct {
	jasaRepo domain.JasaRepository
}

func NewJasaUsecase(repo domain.JasaRepository) domain.JasaUsecase {
	return &jasaUsecase{jasaRepo: repo}
}

func (u *jasaUsecase) Create(ctx context.Context, jasa *domain.MasterJasa) error {
	jasa.NamaJasa = strings.TrimSpace(strings.Title(strings.ToLower(jasa.NamaJasa)))

	if jasa.NamaJasa == "" || jasa.TarifJasa <= 0 {
		return errors.New("nama jasa dan tarif jasa wajib diisi dengan benar")
	}

	// Validasi Bisnis: Komisi mekanik tidak boleh lebih besar dari tarif jasa
	if jasa.PersentaseKomisi <= 0 || jasa.PersentaseKomisi > 100 {
		return errors.New("persentase komisi mekanik harus antara 1% hingga 100%")
	}

	// Kalkulasi otomatis nominal komisi mekanik berdasarkan persentase
	jasa.KomisiMekanik = (jasa.PersentaseKomisi / 100.0) * jasa.TarifJasa

	// validasi created_by
	if strings.TrimSpace(jasa.CreatedBy) == "" {
		return errors.New("nama penginput (created_by) wajib diisi")
	}

	return u.jasaRepo.Create(ctx, jasa)
}

func (u *jasaUsecase) Fetch(ctx context.Context) ([]domain.MasterJasa, error) {
	return u.jasaRepo.Fetch(ctx)
}

func (u *jasaUsecase) GetByID(ctx context.Context, id int) (domain.MasterJasa, error) {
	if id <= 0 {
		return domain.MasterJasa{}, errors.New("ID jasa tidak valid")
	}
	return u.jasaRepo.GetByID(ctx, id)
}

func (u *jasaUsecase) Update(ctx context.Context, jasa *domain.MasterJasa) error {
	jasa.NamaJasa = strings.TrimSpace(strings.Title(strings.ToLower(jasa.NamaJasa)))

	if jasa.IDJasa <= 0 {
		return errors.New("ID jasa wajib disertakan untuk update")
	}
	if jasa.KomisiMekanik > jasa.TarifJasa {
		return errors.New("komisi mekanik tidak boleh melebihi tarif jasa servis")
	}

	return u.jasaRepo.Update(ctx, jasa)
}

func (u *jasaUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID jasa tidak valid")
	}
	return u.jasaRepo.Delete(ctx, id)
}

func (u *jasaUsecase) Search(ctx context.Context, keyword string) ([]domain.MasterJasa, error) {
	return u.jasaRepo.Search(ctx, strings.TrimSpace(keyword))
}
