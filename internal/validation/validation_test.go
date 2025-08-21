package validation

import (
	"testing"

	"otp-auth-service/internal/errors"
)

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		wantErr     bool
		errCode     string
	}{
		{
			name:        "valid phone number with country code",
			phoneNumber: "+1234567890",
			wantErr:     false,
		},
		{
			name:        "valid phone number with longer country code",
			phoneNumber: "+44123456789",
			wantErr:     false,
		},
		{
			name:        "valid phone number with maximum length",
			phoneNumber: "+123456789012345",
			wantErr:     false,
		},
		{
			name:        "empty phone number",
			phoneNumber: "",
			wantErr:     true,
			errCode:     "MISSING_REQUIRED_FIELD",
		},
		{
			name:        "phone number without plus",
			phoneNumber: "1234567890",
			wantErr:     true,
			errCode:     "INVALID_PHONE_NUMBER",
		},
		{
			name:        "phone number with invalid characters",
			phoneNumber: "+123-456-7890",
			wantErr:     true,
			errCode:     "INVALID_PHONE_NUMBER",
		},
		{
			name:        "phone number too short",
			phoneNumber: "+123456",
			wantErr:     true,
			errCode:     "INVALID_PHONE_NUMBER",
		},
		{
			name:        "phone number too long",
			phoneNumber: "+12345678901234567",
			wantErr:     true,
			errCode:     "INVALID_PHONE_NUMBER",
		},
		{
			name:        "phone number with spaces",
			phoneNumber: " +1234567890 ",
			wantErr:     false,
		},
		{
			name:        "phone number starting with 0",
			phoneNumber: "+01234567890",
			wantErr:     true,
			errCode:     "INVALID_PHONE_NUMBER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				domainErr, ok := err.(*errors.DomainError)
				if !ok {
					t.Errorf("Expected DomainError, got %T", err)
					return
				}
				if domainErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, domainErr.Code)
				}
			}
		})
	}
}

func TestValidateOTP(t *testing.T) {
	tests := []struct {
		name    string
		otp     string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid 6-digit OTP",
			otp:     "123456",
			wantErr: false,
		},
		{
			name:    "valid OTP with spaces",
			otp:     " 123456 ",
			wantErr: false,
		},
		{
			name:    "empty OTP",
			otp:     "",
			wantErr: true,
			errCode: "MISSING_REQUIRED_FIELD",
		},
		{
			name:    "OTP too short",
			otp:     "12345",
			wantErr: true,
			errCode: "INVALID_OTP_FORMAT",
		},
		{
			name:    "OTP too long",
			otp:     "1234567",
			wantErr: true,
			errCode: "INVALID_OTP_FORMAT",
		},
		{
			name:    "OTP with letters",
			otp:     "12345a",
			wantErr: true,
			errCode: "INVALID_OTP_FORMAT",
		},
		{
			name:    "OTP with special characters",
			otp:     "123-45",
			wantErr: true,
			errCode: "INVALID_OTP_FORMAT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOTP(tt.otp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				domainErr, ok := err.(*errors.DomainError)
				if !ok {
					t.Errorf("Expected DomainError, got %T", err)
					return
				}
				if domainErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, domainErr.Code)
				}
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid UUID",
			uuid:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID with spaces",
			uuid:    " 550e8400-e29b-41d4-a716-446655440000 ",
			wantErr: false,
		},
		{
			name:    "empty UUID",
			uuid:    "",
			wantErr: true,
			errCode: "MISSING_REQUIRED_FIELD",
		},
		{
			name:    "invalid UUID format",
			uuid:    "invalid-uuid",
			wantErr: true,
			errCode: "INVALID_UUID",
		},
		{
			name:    "UUID without hyphens",
			uuid:    "550e8400e29b41d4a716446655440000",
			wantErr: true,
			errCode: "INVALID_UUID",
		},
		{
			name:    "UUID with wrong length",
			uuid:    "550e8400-e29b-41d4-a716-44665544000",
			wantErr: true,
			errCode: "INVALID_UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				domainErr, ok := err.(*errors.DomainError)
				if !ok {
					t.Errorf("Expected DomainError, got %T", err)
					return
				}
				if domainErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, domainErr.Code)
				}
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		page    int
		limit   int
		wantErr bool
	}{
		{
			name:    "valid pagination",
			page:    1,
			limit:   10,
			wantErr: false,
		},
		{
			name:    "valid pagination with max limit",
			page:    5,
			limit:   100,
			wantErr: false,
		},
		{
			name:    "page too small",
			page:    0,
			limit:   10,
			wantErr: true,
		},
		{
			name:    "page negative",
			page:    -1,
			limit:   10,
			wantErr: true,
		},
		{
			name:    "limit too small",
			page:    1,
			limit:   0,
			wantErr: true,
		},
		{
			name:    "limit too large",
			page:    1,
			limit:   101,
			wantErr: true,
		},
		{
			name:    "limit negative",
			page:    1,
			limit:   -5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePagination(tt.page, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePagination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "empty search query",
			query:   "",
			wantErr: false,
		},
		{
			name:    "valid search query",
			query:   "john doe",
			wantErr: false,
		},
		{
			name:    "valid search query with numbers",
			query:   "user123",
			wantErr: false,
		},
		{
			name:    "valid search query with special chars",
			query:   "user-name_123",
			wantErr: false,
		},
		{
			name:    "search query too short",
			query:   "ab",
			wantErr: true,
		},
		{
			name:    "search query too long",
			query:   "this is a very long search query that exceeds the maximum allowed length of fifty characters",
			wantErr: true,
		},
		{
			name:    "search query with invalid characters",
			query:   "user@domain.com",
			wantErr: true,
		},
		{
			name:    "search query with spaces",
			query:   "  john doe  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSearchQuery(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequestOTP(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		wantErr     bool
	}{
		{
			name:        "valid phone number",
			phoneNumber: "+1234567890",
			wantErr:     false,
		},
		{
			name:        "invalid phone number",
			phoneNumber: "invalid",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequestOTP(tt.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequestOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVerifyOTP(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		otp         string
		wantErr     bool
	}{
		{
			name:        "valid phone number and OTP",
			phoneNumber: "+1234567890",
			otp:         "123456",
			wantErr:     false,
		},
		{
			name:        "invalid phone number",
			phoneNumber: "invalid",
			otp:         "123456",
			wantErr:     true,
		},
		{
			name:        "invalid OTP",
			phoneNumber: "+1234567890",
			otp:         "123",
			wantErr:     true,
		},
		{
			name:        "both invalid",
			phoneNumber: "invalid",
			otp:         "123",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVerifyOTP(tt.phoneNumber, tt.otp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUsers(t *testing.T) {
	tests := []struct {
		name     string
		pageStr  string
		limitStr string
		search   string
		wantErr  bool
	}{
		{
			name:     "valid parameters",
			pageStr:  "1",
			limitStr: "10",
			search:   "john",
			wantErr:  false,
		},
		{
			name:     "empty parameters",
			pageStr:  "",
			limitStr: "",
			search:   "",
			wantErr:  false,
		},
		{
			name:     "invalid page",
			pageStr:  "0",
			limitStr: "10",
			search:   "",
			wantErr:  true,
		},
		{
			name:     "invalid limit",
			pageStr:  "1",
			limitStr: "101",
			search:   "",
			wantErr:  true,
		},
		{
			name:     "invalid search query",
			pageStr:  "1",
			limitStr: "10",
			search:   "ab",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUsers(tt.pageStr, tt.limitStr, tt.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			userID:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			userID:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUser(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParsePositiveInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid positive integer",
			input:   "123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "zero",
			input:   "0",
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative integer",
			input:   "-123",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid string",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePositiveInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePositiveInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePositiveInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
