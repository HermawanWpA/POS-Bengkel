package usecase

import (
	"context"
	"errors"
	"time"

	"pos-echo-app/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) Register(ctx context.Context, req *domain.RegisterRequest) error {
	if req.Username == "" || req.Password == "" || req.NamaUser == "" {
		return errors.New("semua field (username, password, nama) wajib diisi")
	}

	if req.Role != "admin" && req.Role != "kasir" {
		return errors.New("role harus berupa 'admin' atau 'kasir'")
	}

	// 1. Enkripsi Password menggunakan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userModel := domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		NamaUser: req.NamaUser,
		Role:     req.Role,
	}

	return u.userRepo.Create(ctx, &userModel)
}

func (u *userUsecase) Login(ctx context.Context, req *domain.LoginRequest) (domain.LoginResponse, error) {
	var res domain.LoginResponse

	// 1. Cari user berdasarkan username
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return res, errors.New("username atau password salah")
	}

	// 2. Cocokkan password asli dengan password hash di DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, errors.New("username atau password salah")
	}

	// 3. Jika cocok, generate JWT Token (Masa aktif 24 Jam)
	claims := jwt.MapClaims{
		"id_user":  user.IDUser,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token menggunakan JWT Secret Key
	jwtSecret := []byte("KUNCI_RAHASIA_BENGKEL_2026")
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return res, err
	}

	res.Token = tokenString
	res.Username = user.Username
	res.Role = user.Role

	return res, nil
}

func (u *userUsecase) Update(ctx context.Context, id int, req *domain.RegisterRequest) error {
	if req.Username == "" || req.NamaUser == "" {
		return errors.New("username dan nama user wajib diisi")
	}

	if req.Role != "admin" && req.Role != "kasir" {
		return errors.New("role harus berupa 'admin' atau 'kasir'")
	}

	// 1. Ambil data user yang lama dari DB
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	existingUser.Username = req.Username
	existingUser.NamaUser = req.NamaUser
	existingUser.Role = req.Role

	// 2. Cek apakah admin ingin mengganti password karyawan ini
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		existingUser.Password = string(hashedPassword)
	} else {
		// Jika password kiriman frontend kosong, set string kosong agar repo tahu tidak perlu di-update
		existingUser.Password = ""
	}

	return u.userRepo.Update(ctx, &existingUser)
}

func (u *userUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID user tidak valid")
	}

	// Validasi opsional: Mencegah admin menghapus dirinya sendiri secara tidak sengaja
	_, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user memang tidak ditemukan di database")
	}

	return u.userRepo.Delete(ctx, id)
}
