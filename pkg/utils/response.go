package utils

import (
	"log" // PENGEMBANGAN: Untuk mencetak log error di terminal server

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ResponseFormat adalah struktur standar untuk response API yang sukses
type ResponseFormat struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponseFormat adalah struktur standar untuk response API yang gagal/error
type ErrorResponseFormat struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// Pagination meta data untuk response API
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalRows   int64 `json:"total_rows"`
}

// ResponseWithPagination standar untuk data list
type ResponseWithPagination struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Pagination  `json:"meta"`
}

// JSONSuccess mengirimkan HTTP response sukses dengan format yang seragam
func JSONSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, ResponseFormat{
		Message: message,
		Data:    data,
	})
}

// JSONError mengirimkan HTTP response gagal dengan format yang seragam
func JSONError(c echo.Context, statusCode int, message string, err error) error {
	var errMessage string
	if err != nil {
		errMessage = err.Error()

		// =========================================================================
		// PENGEMBANGAN: Cetak error asli ke terminal untuk mempermudah debugging
		// =========================================================================
		log.Printf("[API ERROR] Path: %s | Message: %s | Internal Error: %v", c.Path(), message, err)
	}

	return c.JSON(statusCode, ErrorResponseFormat{
		Message: message,
		Error:   errMessage,
	})
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
