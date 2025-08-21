package services

import (
	"fmt"
	"time"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/models"
	"otp-auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo  repository.UserRepository
	otpRepo   repository.OTPRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, otpRepo repository.OTPRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) RequestOTP(phoneNumber string) (*models.RequestOTPResponse, error) {
	// Generate OTP
	_, err := s.otpRepo.GenerateOTP(phoneNumber)
	if err != nil {
		return nil, err
	}

	return &models.RequestOTPResponse{
		Message:     "OTP sent successfully",
		PhoneNumber: phoneNumber,
	}, nil
}

func (s *AuthService) VerifyOTP(phoneNumber, otp string) (*models.VerifyOTPResponse, error) {
	// Verify OTP
	isValid, err := s.otpRepo.VerifyOTP(phoneNumber, otp)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errors.ErrInvalidOTP
	}

	// Check if user exists
	user, err := s.userRepo.GetByPhoneNumber(phoneNumber)
	isNewUser := false

	if err != nil {
		// User doesn't exist, create new user
		user = models.NewUser(phoneNumber)
		err = s.userRepo.Create(user)
		if err != nil {
			return nil, err
		}
		isNewUser = true
	} else {
		// Update last login time
		user.LastLoginAt = time.Now()
		err = s.userRepo.Update(user)
		if err != nil {
			return nil, err
		}
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID, user.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &models.VerifyOTPResponse{
		Message:   "Authentication successful",
		Token:     token,
		User:      user.ToResponse(),
		IsNewUser: isNewUser,
	}, nil
}

func (s *AuthService) generateJWT(userID, phoneNumber string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      userID,
		"phone_number": phoneNumber,
		"exp":          time.Now().Add(24 * time.Hour).Unix(), // 24 hours expiry
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.ErrInvalidToken.WithDetails(err.Error())
	}

	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return token, nil
}
