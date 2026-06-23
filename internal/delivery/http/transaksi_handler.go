package http

import (
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils" // Pastikan import utils Anda sudah benar
	"strings"

	"github.com/labstack/echo/v4"
)

type TransaksiHandler struct {
	UTransaksi domain.TransaksiUsecase
}

// func NewTransaksiHandler(e *echo.Echo, ut domain.TransaksiUsecase) {
// 	// PERBAIKAN: Ubah uTransaksi menjadi UTransaksi
// 	handler := &TransaksiHandler{UTransaksi: ut}

// 	// Route Utama Pembuatan Work Order / Nota Transaksi Bengkel
// 	e.POST("/api/v1/transaksi", handler.Create)
// 	e.GET("/api/v1/transaksi", handler.Fetch)
// 	e.GET("/api/v1/transaksi/:id", handler.GetByID)
// }

func NewTransaksiHandler(g *echo.Group, up domain.TransaksiUsecase) {
	handler := &TransaksiHandler{UTransaksi: up}

	// Gunakan variabel 'g' bukan 'e'
	g.POST("/transaksi", handler.Create)
	g.GET("/transaksi", handler.Fetch)
	g.GET("/transaksi/:id", handler.GetByID)
	g.PATCH("/transaksi/:id/status", handler.UpdateStatusPengerjaan)

}

func (h *TransaksiHandler) Create(c echo.Context) error {
	var req domain.CreateTransaksiRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format payload data salah", err)
	}

	ctx := c.Request().Context()
	idTransaksi, err := h.UTransaksi.Create(ctx, &req)
	if err != nil {
		// Menangkap segala jenis error bisnis (seperti stok habis, id salah, dll)
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	// Mengembalikan ID Transaksi yang sukses digenerate ke kasir
	return utils.JSONSuccess(c, http.StatusCreated, "Transaksi berhasil disimpan", map[string]string{
		"id_transaksi": idTransaksi,
	})
}

func (h *TransaksiHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	result, err := h.UTransaksi.Fetch(ctx)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil riwayat transaksi", err)
	}
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mengambil riwayat transaksi", result)
}

func (h *TransaksiHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	header, jasas, spareparts, err := h.UTransaksi.GetByID(ctx, id)
	if err != nil {
		return utils.JSONError(c, http.StatusNotFound, "Transaksi tidak ditemukan", err)
	}

	// Membungkus data header dan detail ke dalam satu response JSON yang rapi untuk struk cetak
	responseData := map[string]interface{}{
		"nota_header":      header,
		"detail_jasa":      jasas,
		"detail_sparepart": spareparts,
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil memuat detail nota transaksi", responseData)
}

type UpdateStatusPengerjaanRequest struct {
	StatusPengerjaan string `json:"status_pengerjaan"`
}

func (h *TransaksiHandler) UpdateStatusPengerjaan(c echo.Context) error {
	id := c.Param("id")

	var req UpdateStatusPengerjaanRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}

	ctx := c.Request().Context()

	err := h.UTransaksi.UpdateStatus(ctx, id, req.StatusPengerjaan, "")
	if err != nil {
		// 1. Tangani jika data Transaksi tidak ditemukan (404)
		if strings.Contains(err.Error(), "NOT_FOUND") {
			cleanMessage := strings.Replace(err.Error(), "NOT_FOUND: ", "", 1)
			return utils.JSONError(c, http.StatusNotFound, cleanMessage, err)
		}

		// 2. Tangani jika melanggar aturan alur status 'proses -> selesai' (400)
		if strings.Contains(err.Error(), "BAD_REQUEST") {
			cleanMessage := strings.Replace(err.Error(), "BAD_REQUEST: ", "", 1)
			return utils.JSONError(c, http.StatusBadRequest, cleanMessage, err)
		}

		// 3. Tangani jika ada kendala internal sistem/database (500)
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Status pengerjaan berhasil diperbarui menjadi "+req.StatusPengerjaan, nil)
}
