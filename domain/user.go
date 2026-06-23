package domain

import (
	"context"
	"time"
)

type User struct {
	IDUser    int       `json:"id_user" gorm:"primaryKey;autoIncrement;column:id_user"`
	Username  string    `json:"username" gorm:"type:varchar(50);unique;not null"`
	Password  string    `json:"password,omitempty" gorm:"type:varchar(255);not null"` // omitempty agar password hash tidak ikut tecatat di respon JSON
	NamaUser  string    `json:"nama_user" gorm:"type:varchar(100);not null"`
	Role      string    `json:"role" gorm:"type:enum('admin','kasir');default:kasir"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "user"
}

// DTO untuk Request dari Postman / Frontend
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	NamaUser string `json:"nama_user"`
	Role     string `json:"role"` // 'admin' atau 'kasir'
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Contracts Interface
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByID(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int) error
}

type UserUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, req *LoginRequest) (LoginResponse, error)
	Update(ctx context.Context, id int, req *RegisterRequest) error // Menggunakan RegisterRequest karena field-nya sama
	Delete(ctx context.Context, id int) error
}
