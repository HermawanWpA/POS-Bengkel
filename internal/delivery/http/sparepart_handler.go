package http

import (
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils" // Pastikan import utils Anda sudah benar

	"github.com/labstack/echo/v4"
)

type SparepartHandler struct {
	USparepart domain.SparepartUsecase
}

// func NewSparepartHandler(e *echo.Echo, us domain.SparepartUsecase) {
// 	handler := &SparepartHandler{
// 		USparepart: us,
// 	}

// 	// Routing API Master Sparepart
// 	e.POST("/api/v1/sparepart", handler.Create)
// 	e.GET("/api/v1/sparepart", handler.Fetch)
// 	e.GET("/api/v1/sparepart/search", handler.Search)
// 	e.GET("/api/v1/sparepart/:kode", handler.GetByKode)
// 	e.PUT("/api/v1/sparepart/:kode", handler.Update)
// 	e.DELETE("/api/v1/sparepart/:kode", handler.Delete)
// }

func NewSparepartHandler(g *echo.Group, up domain.SparepartUsecase) {
	handler := &SparepartHandler{USparepart: up}

	// Gunakan variabel 'g' bukan 'e'
	g.POST("/sparepart", handler.Create)
	g.GET("/sparepart", handler.Fetch)
	g.GET("/sparepart/search", handler.Search)
	g.GET("/sparepart/:id", handler.GetByKode)
	g.PUT("/sparepart/:id", handler.Update)
	g.DELETE("/sparepart/:id", handler.Delete)
}

// 1. CREATE SPAREPART
func (h *SparepartHandler) Create(c echo.Context) error {
	var sp domain.MasterSparepart

	// Bind JSON dari Postman ke struct
	if err := c.Bind(&sp); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format data tidak sesuai atau JSON rusak", err)
	}

	// Validasi input wajib sesuai blueprint database
	if sp.KodeSparepart == "" || sp.NamaSparepart == "" || sp.HargaJual <= 0 || sp.HargaBeliHpp <= 0 {
		return utils.JSONError(c, http.StatusBadRequest, "Kode, Nama, HPP, dan Harga Jual wajib diisi dengan benar", nil)
	}

	ctx := c.Request().Context()
	if err := h.USparepart.Create(ctx, &sp); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal menambahkan sparepart", err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "Berhasil menambahkan sparepart baru", sp)
}

// 2. FETCH ALL / GET ALL SPAREPART
func (h *SparepartHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	listSparepart, err := h.USparepart.Fetch(ctx)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil data sparepart", err)
	}

	// Supaya frontend/Postman menerima array kosong [] bukan null jika data di DB kosong
	if len(listSparepart) == 0 {
		return utils.JSONSuccess(c, http.StatusOK, "Data sparepart masih kosong", []domain.MasterSparepart{})
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan seluruh data sparepart", listSparepart)
}

// 3. GET DETAIL BY KODE SPAREPART
func (h *SparepartHandler) GetByKode(c echo.Context) error {
	kode := c.Param("kode")
	if kode == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Kode sparepart tidak boleh kosong", nil)
	}

	ctx := c.Request().Context()
	sp, err := h.USparepart.GetByKode(ctx, kode)
	if err != nil {
		if err.Error() == "record not found" {
			return utils.JSONError(c, http.StatusNotFound, "Sparepart tidak ditemukan", err)
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil detail sparepart", err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan detail sparepart", sp)
}

// 4. UPDATE SPAREPART
func (h *SparepartHandler) Update(c echo.Context) error {
	kode := c.Param("kode")
	if kode == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Parameter kode sparepart tidak valid", nil)
	}

	var sp domain.MasterSparepart
	if err := c.Bind(&sp); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah atau tidak sesuai", err)
	}

	// Paksa kode_sparepart di struct mengikuti parameter URL agar konsisten
	sp.KodeSparepart = kode

	ctx := c.Request().Context()
	if err := h.USparepart.Update(ctx, &sp); err != nil {
		if err.Error() == "record not found" {
			return utils.JSONError(c, http.StatusNotFound, "Data sparepart tidak ditemukan untuk diupdate", err)
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal memperbarui data sparepart", err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil memperbarui data sparepart", sp)
}

// 5. DELETE SPAREPART
func (h *SparepartHandler) Delete(c echo.Context) error {
	kode := c.Param("kode")
	if kode == "" {
		return utils.JSONError(c, http.StatusBadRequest, "Parameter kode sparepart tidak boleh kosong", nil)
	}

	ctx := c.Request().Context()
	if err := h.USparepart.Delete(ctx, kode); err != nil {
		if err.Error() == "record not found" {
			return utils.JSONError(c, http.StatusNotFound, "Data sparepart tidak ditemukan", err)
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal menghapus sparepart", err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil menghapus sparepart", map[string]string{"kode_sparepart": kode})
}

// 6. Search SPAREPART

func (h *SparepartHandler) Search(c echo.Context) error {
	// Mengambil value dari query param (?keyword=...)
	keyword := c.QueryParam("keyword")

	ctx := c.Request().Context()
	result, err := h.USparepart.Search(ctx, keyword)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal melakukan pencarian sparepart", err)
	}

	// Jika tidak ditemukan, kembalikan array kosong [] agar frontend tidak crash
	if len(result) == 0 {
		return utils.JSONSuccess(c, http.StatusOK, "Data sparepart tidak ditemukan", []domain.MasterSparepart{})
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan hasil pencarian", result)
}
