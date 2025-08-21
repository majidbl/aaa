package models

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	phoneNumber := "+1234567890"
	user := NewUser(phoneNumber)

	if user.PhoneNumber != phoneNumber {
		t.Errorf("Expected PhoneNumber to be %s, got %s", phoneNumber, user.PhoneNumber)
	}

	if user.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if user.RegisteredAt.IsZero() {
		t.Error("Expected RegisteredAt to be set")
	}

	if user.LastLoginAt.IsZero() {
		t.Error("Expected LastLoginAt to be set")
	}

	if !user.IsActive {
		t.Error("Expected IsActive to be true")
	}

	// Check that timestamps are recent (within 1 second)
	if time.Since(user.RegisteredAt) > time.Second {
		t.Error("Expected RegisteredAt to be recent")
	}
	if time.Since(user.LastLoginAt) > time.Second {
		t.Error("Expected LastLoginAt to be recent")
	}
}

func TestUser_ToResponse(t *testing.T) {
	user := &User{
		ID:           "test-id",
		PhoneNumber:  "+1234567890",
		RegisteredAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		LastLoginAt:  time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
		IsActive:     true,
	}

	response := user.ToResponse()

	if response.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, response.ID)
	}

	if response.PhoneNumber != user.PhoneNumber {
		t.Errorf("Expected PhoneNumber %s, got %s", user.PhoneNumber, response.PhoneNumber)
	}

	if !response.RegisteredAt.Equal(user.RegisteredAt) {
		t.Errorf("Expected RegisteredAt %v, got %v", user.RegisteredAt, response.RegisteredAt)
	}

	if !response.LastLoginAt.Equal(user.LastLoginAt) {
		t.Errorf("Expected LastLoginAt %v, got %v", user.LastLoginAt, response.LastLoginAt)
	}

	if response.IsActive != user.IsActive {
		t.Errorf("Expected IsActive %v, got %v", user.IsActive, response.IsActive)
	}
}

func TestUser_ToResponse_EdgeCases(t *testing.T) {
	// Test with zero time
	user := &User{
		ID:           "test-id",
		PhoneNumber:  "+1234567890",
		RegisteredAt: time.Time{},
		LastLoginAt:  time.Time{},
		IsActive:     false,
	}

	response := user.ToResponse()

	if response.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, response.ID)
	}

	if response.IsActive != user.IsActive {
		t.Errorf("Expected IsActive %v, got %v", user.IsActive, response.IsActive)
	}
}

func TestNewUser_UniqueIDs(t *testing.T) {
	// Test that multiple users get different IDs
	user1 := NewUser("+1111111111")
	user2 := NewUser("+2222222222")

	if user1.ID == user2.ID {
		t.Error("Expected different IDs for different users")
	}
}

func TestNewUser_PhoneNumber(t *testing.T) {
	phoneNumber := "+1234567890"
	user := NewUser(phoneNumber)

	if user.PhoneNumber != phoneNumber {
		t.Errorf("Expected PhoneNumber to be %s, got %s", phoneNumber, user.PhoneNumber)
	}
}

func TestUser_ImmutableFields(t *testing.T) {
	user := NewUser("+1234567890")
	originalID := user.ID
	originalRegisteredAt := user.RegisteredAt

	// Try to modify fields
	user.ID = "modified-id"
	user.RegisteredAt = time.Now().Add(-24 * time.Hour)

	// Fields should be modified in the instance
	if user.ID != "modified-id" {
		t.Error("Expected ID to be modifiable")
	}

	if user.RegisteredAt.Equal(originalRegisteredAt) {
		t.Error("Expected RegisteredAt to be modifiable")
	}

	// But this doesn't affect the original values
	if originalID == "modified-id" {
		t.Error("Original ID should not be affected")
	}
}
