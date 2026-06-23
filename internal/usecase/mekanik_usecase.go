package usecase

import (
	"context"
	"errors"
	"pos-echo-app/domain"
	"strings"
)

type mekanikUsecase struct {
	mekanikRepo domain.MekanikRepository
}

func NewMekanikUsecase(repo domain.MekanikRepository) domain.MekanikUsecase {
	return &mekanikUsecase{mekanikRepo: repo}
}

func (u *mekanikUsecase) Create(ctx context.Context, m *domain.Mekanik) error {
	m.NamaMekanik = strings.TrimSpace(strings.Title(strings.ToLower(m.NamaMekanik)))
	m.NoHp = strings.TrimSpace(m.NoHp)

	if m.NamaMekanik == "" {
		return errors.New("nama mekanik wajib diisi")
	}

	// Set default status jika kosong
	if m.StatusAktif == "" {
		m.StatusAktif = "aktif"
	}

	return u.mekanikRepo.Create(ctx, m)
}

func (u *mekanikUsecase) Fetch(ctx context.Context) ([]domain.Mekanik, error) {
	return u.mekanikRepo.Fetch(ctx)
}

func (u *mekanikUsecase) GetByID(ctx context.Context, id int) (domain.Mekanik, error) {
	if id <= 0 {
		return domain.Mekanik{}, errors.New("ID mekanik tidak valid")
	}
	return u.mekanikRepo.GetByID(ctx, id)
}

func (u *mekanikUsecase) Update(ctx context.Context, m *domain.Mekanik) error {
	m.NamaMekanik = strings.TrimSpace(strings.Title(strings.ToLower(m.NamaMekanik)))
	m.NoHp = strings.TrimSpace(m.NoHp)
	m.StatusAktif = strings.ToLower(strings.TrimSpace(m.StatusAktif))

	if m.IDMekanik <= 0 {
		return errors.New("ID mekanik wajib disertakan untuk update")
	}
	if m.NamaMekanik == "" {
		return errors.New("nama mekanik tidak boleh dikosongkan")
	}
	if m.StatusAktif != "aktif" && m.StatusAktif != "nonaktif" {
		return errors.New("status harus bernilai 'aktif' atau 'nonaktif'")
	}

	return u.mekanikRepo.Update(ctx, m)
}

func (u *mekanikUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID mekanik tidak valid")
	}
	return u.mekanikRepo.Delete(ctx, id)
}

func (u *mekanikUsecase) Search(ctx context.Context, keyword string) ([]domain.Mekanik, error) {
	return u.mekanikRepo.Search(ctx, strings.TrimSpace(keyword))
}
