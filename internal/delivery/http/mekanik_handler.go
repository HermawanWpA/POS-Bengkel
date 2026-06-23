package http

import (
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils" // Pastikan import utils Anda sudah benar
	"strconv"

	"github.com/labstack/echo/v4"
)

type MekanikHandler struct {
	UMekanik domain.MekanikUsecase
}

// func NewMekanikHandler(e *echo.Echo, um domain.MekanikUsecase) {
// 	handler := &MekanikHandler{UMekanik: um}

// 	e.POST("/api/v1/mekanik", handler.Create)
// 	e.GET("/api/v1/mekanik", handler.Fetch)
// 	e.GET("/api/v1/mekanik/search", handler.Search)
// 	e.GET("/api/v1/mekanik/:id", handler.GetByID)
// 	e.PUT("/api/v1/mekanik/:id", handler.Update)
// 	e.DELETE("/api/v1/mekanik/:id", handler.Delete)
// }

func NewMekanikHandler(g *echo.Group, up domain.MekanikUsecase) {
	handler := &MekanikHandler{UMekanik: up}

	// Gunakan variabel 'g' bukan 'e'
	g.POST("/mekanik", handler.Create)
	g.GET("/mekanik", handler.Fetch)
	g.GET("/mekanik/search", handler.Search)
	g.GET("/mekanik/:id", handler.GetByID)
	g.PUT("/mekanik/:id", handler.Update)
	g.DELETE("/mekanik/:id", handler.Delete)
}

func (h *MekanikHandler) Create(c echo.Context) error {
	var m domain.Mekanik
	if err := c.Bind(&m); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UMekanik.Create(ctx, &m); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "Berhasil menambahkan mekanik baru", m)
}

func (h *MekanikHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	listMekanik, err := h.UMekanik.Fetch(ctx)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil data mekanik", err)
	}

	if len(listMekanik) == 0 {
		return utils.JSONSuccess(c, http.StatusOK, "Data mekanik masih kosong", []domain.Mekanik{})
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan seluruh data mekanik", listMekanik)
}

func (h *MekanikHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	ctx := c.Request().Context()
	m, err := h.UMekanik.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return utils.JSONError(c, http.StatusNotFound, "Mekanik tidak ditemukan", err)
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil detail mekanik", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan detail mekanik", m)
}

func (h *MekanikHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	var m domain.Mekanik
	if err := c.Bind(&m); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}
	m.IDMekanik = id

	ctx := c.Request().Context()
	if err := h.UMekanik.Update(ctx, &m); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil memperbarui data mekanik", m)
}

func (h *MekanikHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UMekanik.Delete(ctx, id); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal menghapus data mekanik", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil menghapus data mekanik", map[string]int{"id_mekanik": id})
}

func (h *MekanikHandler) Search(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	ctx := c.Request().Context()
	result, err := h.UMekanik.Search(ctx, keyword)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mencari data mekanik", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Hasil pencarian mekanik", result)
}
