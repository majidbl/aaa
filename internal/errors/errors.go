package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	HTTPStatus int    `json:"-"`
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// New creates a new domain error
func New(code, message string, httpStatus int) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WithDetails adds details to the error
func (e *DomainError) WithDetails(details string) *DomainError {
	return &DomainError{
		Code:       e.Code,
		Message:    e.Message,
		Details:    details,
		HTTPStatus: e.HTTPStatus,
	}
}

// Common domain errors
var (
	// Authentication errors
	ErrInvalidOTP        = New("INVALID_OTP", "Invalid OTP provided", http.StatusUnauthorized)
	ErrOTPExpired        = New("OTP_EXPIRED", "OTP has expired", http.StatusUnauthorized)
	ErrOTPNotFound       = New("OTP_NOT_FOUND", "OTP not found", http.StatusUnauthorized)
	ErrTooManyAttempts   = New("TOO_MANY_ATTEMPTS", "Too many failed attempts", http.StatusUnauthorized)
	ErrRateLimitExceeded = New("RATE_LIMIT_EXCEEDED", "Rate limit exceeded", http.StatusTooManyRequests)
	ErrInvalidToken      = New("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	ErrMissingAuthHeader = New("MISSING_AUTH_HEADER", "Authorization header is required", http.StatusUnauthorized)
	ErrInvalidAuthFormat = New("INVALID_AUTH_FORMAT", "Invalid authorization header format", http.StatusUnauthorized)

	// User errors
	ErrUserNotFound      = New("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("USER_ALREADY_EXISTS", "User with this phone number already exists", http.StatusConflict)
	ErrInvalidUserID     = New("INVALID_USER_ID", "Invalid user ID", http.StatusBadRequest)

	// Validation errors
	ErrInvalidRequest       = New("INVALID_REQUEST", "Invalid request body", http.StatusBadRequest)
	ErrInvalidPhoneNumber   = New("INVALID_PHONE_NUMBER", "Invalid phone number format", http.StatusBadRequest)
	ErrMissingRequiredField = New("MISSING_REQUIRED_FIELD", "Required field is missing", http.StatusBadRequest)
	ErrInvalidOTPFormat     = New("INVALID_OTP_FORMAT", "Invalid OTP format", http.StatusBadRequest)
	ErrInvalidUUID          = New("INVALID_UUID", "Invalid UUID format", http.StatusBadRequest)
	ErrInvalidPagination    = New("INVALID_PAGINATION", "Invalid pagination parameters", http.StatusBadRequest)
	ErrInvalidSearchQuery    = New("INVALID_SEARCH_QUERY", "Invalid search query", http.StatusBadRequest)

	// Internal errors
	ErrInternalServer = New("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrDatabaseError  = New("DATABASE_ERROR", "Database operation failed", http.StatusInternalServerError)
	ErrRedisError     = New("REDIS_ERROR", "Redis operation failed", http.StatusInternalServerError)
)

// IsDomainError checks if an error is a domain error
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// GetDomainError extracts domain error from an error
func GetDomainError(err error) *DomainError {
	if err == nil {
		return ErrInternalServer
	}
	
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr
	}
	return ErrInternalServer.WithDetails(err.Error())
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
