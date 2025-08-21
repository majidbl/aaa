package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/models"

	"github.com/redis/go-redis/v9"
)

type OTPRepository interface {
	GenerateOTP(phoneNumber string) (string, error)
	VerifyOTP(phoneNumber, otp string) (bool, error)
	IsRateLimited(phoneNumber string) (bool, error)
	GetOTP(phoneNumber string) (*models.OTP, error)
}

type RedisOTPRepository struct {
	client *redis.Client
}

func NewOTPRepository(client *redis.Client) OTPRepository {
	return &RedisOTPRepository{
		client: client,
	}
}

func (r *RedisOTPRepository) GenerateOTP(phoneNumber string) (string, error) {
	ctx := context.Background()

	// Check rate limiting
	isLimited, err := r.IsRateLimited(phoneNumber)
	if err != nil {
		return "", err
	}
	if isLimited {
		return "", errors.ErrRateLimitExceeded.WithDetails(fmt.Sprintf("phone number: %s", phoneNumber))
	}

	// Generate 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Create OTP object
	otpData := &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        otp,
		ExpiresAt:   time.Now().Add(2 * time.Minute), // 2 minutes expiry
		Attempts:    0,
	}

	// Store OTP in Redis with 2 minutes expiry
	otpKey := fmt.Sprintf("otp:%s", phoneNumber)
	otpJSON, _ := json.Marshal(otpData)

	err = r.client.Set(ctx, otpKey, otpJSON, 2*time.Minute).Err()
	if err != nil {
		return "", err
	}

	// Update rate limiting counter
	rateLimitKey := fmt.Sprintf("rate_limit:%s", phoneNumber)
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, rateLimitKey)
	pipe.Expire(ctx, rateLimitKey, 10*time.Minute) // 10 minutes window
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", err
	}

	// Print OTP to console (for development)
	fmt.Printf("OTP for %s: %s\n", phoneNumber, otp)

	return otp, nil
}

func (r *RedisOTPRepository) VerifyOTP(phoneNumber, otp string) (bool, error) {
	ctx := context.Background()

	// Get OTP from Redis
	otpKey := fmt.Sprintf("otp:%s", phoneNumber)
	otpJSON, err := r.client.Get(ctx, otpKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, errors.ErrOTPNotFound
		}
		return false, err
	}

	var otpData models.OTP
	err = json.Unmarshal([]byte(otpJSON), &otpData)
	if err != nil {
		return false, err
	}

	// Check if OTP is expired
	if time.Now().After(otpData.ExpiresAt) {
		// Delete expired OTP
		r.client.Del(ctx, otpKey)
		return false, errors.ErrOTPExpired
	}

	// Check if OTP matches
	if otpData.Code != otp {
				// Increment attempts
		otpData.Attempts++
		if otpData.Attempts >= 3 {
			// Delete OTP after 3 failed attempts
			r.client.Del(ctx, otpKey)
			return false, errors.ErrTooManyAttempts
		}
		
		// Update OTP data
		otpJSON, _ := json.Marshal(otpData)
		r.client.Set(ctx, otpKey, otpJSON, 2*time.Minute)
		return false, errors.ErrInvalidOTP
	}

	// OTP is valid, delete it
	r.client.Del(ctx, otpKey)
	return true, nil
}

func (r *RedisOTPRepository) IsRateLimited(phoneNumber string) (bool, error) {
	ctx := context.Background()

	rateLimitKey := fmt.Sprintf("rate_limit:%s", phoneNumber)
	count, err := r.client.Get(ctx, rateLimitKey).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// Allow max 3 requests per 10 minutes
	return count >= 3, nil
}

func (r *RedisOTPRepository) GetOTP(phoneNumber string) (*models.OTP, error) {
	ctx := context.Background()

	otpKey := fmt.Sprintf("otp:%s", phoneNumber)
	otpJSON, err := r.client.Get(ctx, otpKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("OTP not found")
		}
		return nil, err
	}

	var otpData models.OTP
	err = json.Unmarshal([]byte(otpJSON), &otpData)
	if err != nil {
		return nil, err
	}

	return &otpData, nil
}
