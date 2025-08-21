package models

import "time"

type OTP struct {
	PhoneNumber string    `json:"phone_number"`
	Code        string    `json:"code"`
	ExpiresAt   time.Time `json:"expires_at"`
	Attempts    int       `json:"attempts"`
}

type RequestOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type RequestOTPResponse struct {
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	Message   string        `json:"message"`
	Token     string        `json:"token"`
	User      *UserResponse `json:"user"`
	IsNewUser bool          `json:"is_new_user"`
}

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
