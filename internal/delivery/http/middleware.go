package http

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware akan memvalidasi JWT Token yang dikirim dari Frontend / Postman
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Ambil header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  false,
				"message": "Akses ditolak: Anda belum login (Token tidak ditemukan)",
			})
		}

		// 2. Format token biasanya: "Bearer <token_string>", kita ambil tokennya saja
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Parsing dan validasi token menggunakan Kunci Rahasia yang sama saat login
		jwtSecret := []byte("KUNCI_RAHASIA_BENGKEL_2026")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Metode enkripsi token tidak valid")
			}
			return jwtSecret, nil
		})

		// 4. Jika token rusak atau kedaluwarsa (Expired)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  false,
				"message": "Sesi login Anda telah habis atau token tidak sah. Silakan login kembali",
			})
		}

		// 5. Jika valid, simpan data claims (id, username, role) ke dalam context Echo
		// Supaya handler lain bisa tahu siapa nama user yang sedang mengakses rute ini
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("user_id", claims["id_user"])
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
		}

		// Lanjutkan ke halaman/endpoint yang dituju
		return next(c)
	}
}
