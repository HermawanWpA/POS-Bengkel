package http

import (
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils" // Pastikan import utils Anda sudah benar
	"strconv"

	"github.com/labstack/echo/v4"
)

type JasaHandler struct {
	UJasa domain.JasaUsecase
}

// func NewJasaHandler(e *echo.Echo, us domain.JasaUsecase) {
// 	handler := &JasaHandler{UJasa: us}

// 	e.POST("/api/v1/jasa", handler.Create)
// 	e.GET("/api/v1/jasa", handler.Fetch)
// 	e.GET("/api/v1/jasa/search", handler.Search)
// 	e.GET("/api/v1/jasa/:id", handler.GetByID)
// 	e.PUT("/api/v1/jasa/:id", handler.Update)
// 	e.DELETE("/api/v1/jasa/:id", handler.Delete)
// }

func NewJasaHandler(g *echo.Group, up domain.JasaUsecase) {
	handler := &JasaHandler{UJasa: up}

	// Gunakan variabel 'g' bukan 'e'
	g.POST("/jasa", handler.Create)
	g.GET("/jasa", handler.Fetch)
	g.GET("/jasa/search", handler.Search)
	g.GET("/jasa/:id", handler.GetByID)
	g.PUT("/jasa/:id", handler.Update)
	g.DELETE("/jasa/:id", handler.Delete)
}

func (h *JasaHandler) Create(c echo.Context) error {
	var jasa domain.MasterJasa
	if err := c.Bind(&jasa); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UJasa.Create(ctx, &jasa); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "Berhasil menambahkan data jasa baru", jasa)
}

func (h *JasaHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	listJasa, err := h.UJasa.Fetch(ctx)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil data jasa", err)
	}

	if len(listJasa) == 0 {
		return utils.JSONSuccess(c, http.StatusOK, "Data jasa masih kosong", []domain.MasterJasa{})
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan seluruh data jasa", listJasa)
}

func (h *JasaHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	ctx := c.Request().Context()
	jasa, err := h.UJasa.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return utils.JSONError(c, http.StatusNotFound, "Data jasa tidak ditemukan", err)
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil detail jasa", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan detail jasa", jasa)
}

func (h *JasaHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	var jasa domain.MasterJasa
	if err := c.Bind(&jasa); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}
	jasa.IDJasa = id

	ctx := c.Request().Context()
	if err := h.UJasa.Update(ctx, &jasa); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil memperbarui data jasa", jasa)
}

func (h *JasaHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UJasa.Delete(ctx, id); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal menghapus data jasa", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil menghapus data jasa", map[string]int{"id_jasa": id})
}

func (h *JasaHandler) Search(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	ctx := c.Request().Context()
	result, err := h.UJasa.Search(ctx, keyword)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mencari data jasa", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Hasil pencarian jasa", result)
}
