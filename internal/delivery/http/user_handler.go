package http

import (
	"net/http"
	"pos-echo-app/domain"
	"pos-echo-app/pkg/utils" // Pastikan import utils Anda sudah benar
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UUser domain.UserUsecase
}

// func NewUserHandler(e *echo.Echo, uu domain.UserUsecase) {
// 	handler := &UserHandler{UUser: uu}

// 	// Route tanpa proteksi auth (Public)
// 	e.POST("/api/v1/auth/register", handler.Register)
// 	e.POST("/api/v1/auth/login", handler.Login)
// 	e.PUT("/api/v1/users/:id", handler.Update)
// 	e.DELETE("/api/v1/users/:id", handler.Delete)
// }

func NewUserHandler(e *echo.Echo, g *echo.Group, uu domain.UserUsecase) {
	handler := &UserHandler{UUser: uu}

	// 1. RUTE PUBLIC: Didaftarkan ke objek 'e' (Bisa ditembak langsung)
	e.POST("/api/v1/auth/register", handler.Register)
	e.POST("/api/v1/auth/login", handler.Login)

	// 2. RUTE PROTECTED: Didaftarkan ke objek 'g' (Wajib bawa token JWT)
	// Karena 'g' sudah punya prefix "/api/v1", rute di bawah ini otomatis menjadi:
	// PUT /api/v1/users/:id  dan  DELETE /api/v1/users/:id
	g.PUT("/users/:id", handler.Update)
	g.DELETE("/users/:id", handler.Delete)
}

func (h *UserHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UUser.Register(ctx, &req); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusCreated, "Registrasi user baru berhasil", nil)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}

	ctx := c.Request().Context()
	res, err := h.UUser.Login(ctx, &req)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Login sukses", res)
}

func (h *UserHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format JSON salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UUser.Update(ctx, id, &req); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil memperbarui data user/karyawan", nil)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Format ID salah", err)
	}

	ctx := c.Request().Context()
	if err := h.UUser.Delete(ctx, id); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error(), err)
	}

	return utils.JSONSuccess(c, http.StatusOK, "Berhasil menghapus akun user dari sistem", map[string]int{"id_user": id})
}
