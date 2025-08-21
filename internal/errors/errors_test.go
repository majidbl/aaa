package errors

import (
	"errors"
	"testing"
)

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *DomainError
		want    string
	}{
		{
			name: "error without details",
			err: &DomainError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
			},
			want: "TEST_ERROR: Test error message",
		},
		{
			name: "error with details",
			err: &DomainError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
				Details: "Additional details",
			},
			want: "TEST_ERROR: Test error message (Additional details)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("DomainError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := New("TEST_CODE", "Test message", 400)
	
	if err.Code != "TEST_CODE" {
		t.Errorf("Expected Code to be 'TEST_CODE', got %s", err.Code)
	}
	
	if err.Message != "Test message" {
		t.Errorf("Expected Message to be 'Test message', got %s", err.Message)
	}
	
	if err.HTTPStatus != 400 {
		t.Errorf("Expected HTTPStatus to be 400, got %d", err.HTTPStatus)
	}
}

func TestDomainError_WithDetails(t *testing.T) {
	err := New("TEST_CODE", "Test message", 400)
	errWithDetails := err.WithDetails("Additional details")
	
	if errWithDetails.Details != "Additional details" {
		t.Errorf("Expected Details to be 'Additional details', got %s", errWithDetails.Details)
	}
	
	// Original error should not be modified
	if err.Details != "" {
		t.Errorf("Original error should not be modified, got Details: %s", err.Details)
	}
}

func TestIsDomainError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "domain error",
			err:  ErrInvalidRequest,
			want: true,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDomainError(tt.err); got != tt.want {
				t.Errorf("IsDomainError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDomainError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "domain error",
			err:  ErrInvalidRequest,
			want: "INVALID_REQUEST",
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: "INTERNAL_SERVER_ERROR",
		},
		{
			name: "nil error",
			err:  nil,
			want: "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDomainError(tt.err)
			if got.Code != tt.want {
				t.Errorf("GetDomainError().Code = %v, want %v", got.Code, tt.want)
			}
		})
	}
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "domain error",
			err:  ErrInvalidRequest,
			want: 400,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: 500,
		},
		{
			name: "nil error",
			err:  nil,
			want: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHTTPStatus(tt.err); got != tt.want {
				t.Errorf("GetHTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Test that all predefined errors have the correct structure
	errorTests := []struct {
		name string
		err  *DomainError
	}{
		{"ErrInvalidOTP", ErrInvalidOTP},
		{"ErrOTPExpired", ErrOTPExpired},
		{"ErrOTPNotFound", ErrOTPNotFound},
		{"ErrTooManyAttempts", ErrTooManyAttempts},
		{"ErrRateLimitExceeded", ErrRateLimitExceeded},
		{"ErrInvalidToken", ErrInvalidToken},
		{"ErrMissingAuthHeader", ErrMissingAuthHeader},
		{"ErrInvalidAuthFormat", ErrInvalidAuthFormat},
		{"ErrUserNotFound", ErrUserNotFound},
		{"ErrUserAlreadyExists", ErrUserAlreadyExists},
		{"ErrInvalidUserID", ErrInvalidUserID},
		{"ErrInvalidRequest", ErrInvalidRequest},
		{"ErrInvalidPhoneNumber", ErrInvalidPhoneNumber},
		{"ErrMissingRequiredField", ErrMissingRequiredField},
		{"ErrInvalidOTPFormat", ErrInvalidOTPFormat},
		{"ErrInvalidUUID", ErrInvalidUUID},
		{"ErrInvalidPagination", ErrInvalidPagination},
		{"ErrInvalidSearchQuery", ErrInvalidSearchQuery},
		{"ErrInternalServer", ErrInternalServer},
		{"ErrDatabaseError", ErrDatabaseError},
		{"ErrRedisError", ErrRedisError},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code == "" {
				t.Errorf("%s has empty Code", tt.name)
			}
			if tt.err.Message == "" {
				t.Errorf("%s has empty Message", tt.name)
			}
			if tt.err.HTTPStatus == 0 {
				t.Errorf("%s has zero HTTPStatus", tt.name)
			}
		})
	}
}

func TestErrorWithDetails(t *testing.T) {
	// Test that WithDetails creates a new error instance
	originalErr := ErrInvalidRequest
	errWithDetails := originalErr.WithDetails("test details")
	
	if originalErr.Details != "" {
		t.Error("Original error should not be modified")
	}
	
	if errWithDetails.Details != "test details" {
		t.Errorf("Expected details 'test details', got '%s'", errWithDetails.Details)
	}
	
	// Test that the new error is a copy
	if originalErr == errWithDetails {
		t.Error("WithDetails should return a new error instance")
	}
}
