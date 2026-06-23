package usecase

import (
	"context"
	"pos-echo-app/domain"
)

type dashboardUsecase struct {
	dashRepo domain.DashboardRepository
}

func NewDashboardUsecase(dr domain.DashboardRepository) domain.DashboardUsecase {
	return &dashboardUsecase{dashRepo: dr}
}

func (u *dashboardUsecase) GetStats(ctx context.Context) (domain.DashboardStatsResponse, error) {
	return u.dashRepo.GetStats(ctx)
}

func (u *dashboardUsecase) GetDetailsByPeriod(ctx context.Context, period string) ([]domain.ServiceDetailRow, error) {
	return u.dashRepo.GetDetailsByPeriod(ctx, period)
}
