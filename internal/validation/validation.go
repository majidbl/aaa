package validation

import (
	"fmt"
	"regexp"
	"strings"

	"otp-auth-service/internal/errors"
)

// PhoneNumberRegex is a regex pattern for international phone numbers
var PhoneNumberRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// OTPRegex is a regex pattern for 6-digit OTP codes
var OTPRegex = regexp.MustCompile(`^\d{6}$`)

// UUIDRegex is a regex pattern for UUID validation
var UUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ValidatePhoneNumber validates phone number format
func ValidatePhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return errors.ErrMissingRequiredField.WithDetails("phone_number is required")
	}

	phoneNumber = strings.TrimSpace(phoneNumber)
	
	if !PhoneNumberRegex.MatchString(phoneNumber) {
		return errors.ErrInvalidPhoneNumber.WithDetails(
			fmt.Sprintf("phone number '%s' must be in international format (e.g., +1234567890)", phoneNumber),
		)
	}

	// Additional validation: check length
	if len(phoneNumber) < 10 || len(phoneNumber) > 16 {
		return errors.ErrInvalidPhoneNumber.WithDetails(
			fmt.Sprintf("phone number length must be between 10 and 16 characters, got %d", len(phoneNumber)),
		)
	}

	return nil
}

// ValidateOTP validates OTP format
func ValidateOTP(otp string) error {
	if otp == "" {
		return errors.ErrMissingRequiredField.WithDetails("otp is required")
	}

	otp = strings.TrimSpace(otp)
	
	if !OTPRegex.MatchString(otp) {
		return errors.ErrInvalidOTPFormat.WithDetails(
			fmt.Sprintf("OTP must be exactly 6 digits, got '%s'", otp),
		)
	}

	return nil
}

// ValidateUUID validates UUID format
func ValidateUUID(uuid string) error {
	if uuid == "" {
		return errors.ErrMissingRequiredField.WithDetails("user ID is required")
	}

	uuid = strings.TrimSpace(uuid)
	
	if !UUIDRegex.MatchString(uuid) {
		return errors.ErrInvalidUUID.WithDetails(
			fmt.Sprintf("invalid UUID format: '%s'", uuid),
		)
	}

	return nil
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, limit int) error {
	if page < 1 {
		return errors.ErrInvalidRequest.WithDetails("page must be greater than 0")
	}

	if limit < 1 || limit > 100 {
		return errors.ErrInvalidRequest.WithDetails("limit must be between 1 and 100")
	}

	return nil
}

// ValidateSearchQuery validates search query parameters
func ValidateSearchQuery(query string) error {
	if query == "" {
		return nil // Empty search is valid
	}

	query = strings.TrimSpace(query)
	
	// Check minimum length
	if len(query) < 3 {
		return errors.ErrInvalidRequest.WithDetails("search query must be at least 3 characters long")
	}

	// Check maximum length
	if len(query) > 50 {
		return errors.ErrInvalidRequest.WithDetails("search query must be less than 50 characters")
	}

	// Check for valid characters (alphanumeric, spaces, +, -, _)
	validQueryRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\+\-_]+$`)
	if !validQueryRegex.MatchString(query) {
		return errors.ErrInvalidRequest.WithDetails("search query contains invalid characters")
	}

	return nil
}

// ValidateRequestOTP validates RequestOTP request
func ValidateRequestOTP(phoneNumber string) error {
	return ValidatePhoneNumber(phoneNumber)
}

// ValidateVerifyOTP validates VerifyOTP request
func ValidateVerifyOTP(phoneNumber, otp string) error {
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return err
	}
	
	if err := ValidateOTP(otp); err != nil {
		return err
	}

	return nil
}

// ValidateGetUsers validates GetUsers request parameters
func ValidateGetUsers(pageStr, limitStr, search string) error {
	// Parse and validate page
	page := 1
	if pageStr != "" {
		if p, err := parsePositiveInt(pageStr); err != nil {
			return errors.ErrInvalidRequest.WithDetails("invalid page number")
		} else {
			page = p
		}
	}

	// Parse and validate limit
	limit := 10
	if limitStr != "" {
		if l, err := parsePositiveInt(limitStr); err != nil {
			return errors.ErrInvalidRequest.WithDetails("invalid limit number")
		} else {
			limit = l
		}
	}

	// Validate pagination
	if err := ValidatePagination(page, limit); err != nil {
		return errors.ErrInvalidPagination.WithDetails(err.Error())
	}

	// Validate search query
	if err := ValidateSearchQuery(search); err != nil {
		return errors.ErrInvalidSearchQuery.WithDetails(err.Error())
	}

	return nil
}

// ValidateGetUser validates GetUser request parameters
func ValidateGetUser(userID string) error {
	return ValidateUUID(userID)
}

// parsePositiveInt parses a string to a positive integer
func parsePositiveInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, err
	}
	if result <= 0 {
		return 0, fmt.Errorf("must be positive")
	}
	return result, nil
}
