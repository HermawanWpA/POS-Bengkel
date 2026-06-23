package http

import (
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	UDash domain.DashboardUsecase
}

func NewDashboardHandler(g *echo.Group, ud domain.DashboardUsecase) {
	handler := &DashboardHandler{UDash: ud}

	// Endpoint untuk memuat ringkasan angka statistik dashboard
	g.GET("/dashboard/stats", handler.GetDashboardStats)

	// Endpoint untuk memuat list tabel detail saat kartu/tombol periode diklik
	g.GET("/dashboard/details", handler.GetDashboardDetails)
}

func (h *DashboardHandler) GetDashboardStats(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.UDash.GetStats(ctx)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal memuat statistik dashboard", err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Statistik dashboard berhasil dimuat", data)
}

func (h *DashboardHandler) GetDashboardDetails(c echo.Context) error {
	ctx := c.Request().Context()

	// Tangkap query param, contoh: /api/v1/dashboard/details?period=hari
	period := c.QueryParam("period")
	if period == "" {
		period = "hari" // Default jika parameter kosong
	}

	data, err := h.UDash.GetDetailsByPeriod(ctx, period)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Detail data servis berhasil dimuat", data)
}
