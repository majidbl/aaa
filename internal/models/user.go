package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `json:"id"`
	PhoneNumber  string    `json:"phone_number"`
	RegisteredAt time.Time `json:"registered_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
	IsActive     bool      `json:"is_active"`
}

type CreateUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UserResponse struct {
	ID           string    `json:"id"`
	PhoneNumber  string    `json:"phone_number"`
	RegisteredAt time.Time `json:"registered_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
	IsActive     bool      `json:"is_active"`
}

func NewUser(phoneNumber string) *User {
	return &User{
		ID:           uuid.New().String(),
		PhoneNumber:  phoneNumber,
		RegisteredAt: time.Now(),
		LastLoginAt:  time.Now(),
		IsActive:     true,
	}
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		PhoneNumber:  u.PhoneNumber,
		RegisteredAt: u.RegisteredAt,
		LastLoginAt:  u.LastLoginAt,
		IsActive:     u.IsActive,
	}
}
