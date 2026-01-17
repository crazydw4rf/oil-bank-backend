package dto

import "github.com/crazydw4rf/oil-bank-backend/internal/entity"

// TODO: tambahkan validator nanti

type UserCreateRequest struct {
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	UserType entity.UserType `json:"user_type"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
