package http

import (
	"math"
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils" // Pastikan import utils Anda sudah benar
	"strconv"

	"github.com/labstack/echo/v4"
)

type PelangganHandler struct {
	UUsecase domain.PelangganUsecase
}

// func NewPelangganHandler(e *echo.Echo, us domain.PelangganUsecase) {
// 	handler := &PelangganHandler{
// 		UUsecase: us,
// 	}

//		// Gunakan langsung instance 'e' tanpa group bertumpuk agar jalurnya absolut dan pasti
//		e.POST("/api/v1/pelanggan", handler.Create)
//		e.GET("/api/v1/pelanggan", handler.FetchWithVehicles)
//		e.GET("/api/v1/pelanggan/search", handler.Search)
//		e.GET("/api/v1/pelanggan/:id", handler.GetByID)
//		e.PUT("/api/v1/pelanggan/:id", handler.Update)
//		e.DELETE("/api/v1/pelanggan/:id", handler.Delete) // Pastikan polanya murni /:id di ujung
//	}
func NewPelangganHandler(g *echo.Group, up domain.PelangganUsecase) {
	handler := &PelangganHandler{UUsecase: up}

	// Gunakan variabel 'g' bukan 'e'
	g.POST("/pelanggan", handler.Create)
	g.GET("/pelanggan", handler.FetchWithVehicles)
	g.GET("/pelanggan/search", handler.Search)
	g.GET("/pelanggan/:id", handler.GetByID)
	g.PUT("/pelanggan/:id", handler.Update)
	g.DELETE("/pelanggan/:id", handler.Delete)
	g.GET("/pelanggan", handler.GetPelanggan)
}

func (h *PelangganHandler) Create(c echo.Context) error {
	var pelanggan domain.Pelanggan

	if err := c.Bind(&pelanggan); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format data salah", err)
	}

	if err := h.UUsecase.Create(c.Request().Context(), &pelanggan); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal menyimpan data pelanggan", err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "Berhasil menambahkan pelanggan baru", pelanggan)
}

func (h *PelangganHandler) FetchWithVehicles(c echo.Context) error {
	// Ambil data dari usecase dengan melewatkan request context
	pelangganList, err := h.UUsecase.FetchWithVehicles(c.Request().Context())
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil data pelanggan", err)
	}

	// Kembalikan response sukses menggunakan helper utils Anda
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mengambil data pelanggan beserta kendaraan", pelangganList)
}

func (h *PelangganHandler) GetByID(c echo.Context) error {
	// Ambil id dari parameter URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	// Panggil usecase
	pelanggan, err := h.UUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		// Jika GORM mengembalikan error record not found, ubah status menjadi 404
		if err.Error() == "record not found" {
			return utils.JSONError(c, http.StatusNotFound, "Data pelanggan tidak ditemukan", err)
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal mengambil detail pelanggan", err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan detail pelanggan", pelanggan)
}

func (h *PelangganHandler) Search(c echo.Context) error {
	// Mengambil value dari query param (?keyword=...)
	keyword := c.QueryParam("keyword")

	result, err := h.UUsecase.Search(c.Request().Context(), keyword)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal melakukan pencarian pelanggan", err)
	}

	if len(result) == 0 {
		return utils.JSONSuccess(c, http.StatusOK, "Data pelanggan tidak ditemukan", []domain.Pelanggan{})
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil mendapatkan hasil pencarian", result)
}

func (h *PelangganHandler) Delete(c echo.Context) error {
	// Ambil id dari parameter URL (misal: /api/v1/pelanggan/1)
	idStr := c.Param("id")

	// Konversi id dari string ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	// Panggil usecase untuk menghapus data
	if err := h.UUsecase.Delete(c.Request().Context(), id); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal menghapus data pelanggan", err)
	}

	// Kembalikan response sukses tanpa membawa data (cukup kirim pesan)
	return utils.JSONSuccess(c, http.StatusOK, "Berhasil menghapus data pelanggan beserta kendaraannya", nil)
}

func (h *PelangganHandler) Update(c echo.Context) error {
	// Ambil ID dari parameter URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	var pelanggan domain.Pelanggan
	// Bind JSON body dari Postman ke struct pelanggan
	if err := c.Bind(&pelanggan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":      "JSON Anda Rusak / Salah Tipe Data",
			"error_detail": err.Error(), // Ini akan memberi tahu kolom mana yang error
		})
	}

	pelanggan.ID = id

	if err := h.UUsecase.Update(c.Request().Context(), &pelanggan); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal memperbarui data pelanggan", err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil memperbarui data pelanggan", pelanggan)
}

func (h *PelangganHandler) GetPelanggan(c echo.Context) error {
	// 1. Ambil Query Param menggunakan c.QueryParam() khas Echo
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	// 2. Konversi ke Integer
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	param := domain.PaginationParam{
		Page:  page,
		Limit: limit,
	}

	// 3. Panggil UseCase
	customers, totalRows, err := h.UUsecase.GetAllWithPagination(c.Request().Context(), param)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Gagal memuat data paginasi", err)
	}

	// 4. Hitung total halaman
	activeLimit := param.Limit
	if activeLimit <= 0 {
		activeLimit = 10
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(activeLimit)))

	// 5. Bungkus dengan struktur response pagination yang seragam
	response := utils.ResponseWithPagination{
		Status:  "success",
		Message: "Berhasil mengambil data pelanggan dengan paginasi",
		Data:    customers,
		Meta: utils.Pagination{
			CurrentPage: page,
			PageSize:    activeLimit,
			TotalPages:  totalPages,
			TotalRows:   totalRows,
		},
	}

	// Mengembalikan JSON standar Echo
	return c.JSON(http.StatusOK, response)
}
