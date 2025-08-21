package repository

import (
	"testing"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/models"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository()
	if repo == nil {
		t.Error("Expected repository to be created")
	}
}

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()
	user := models.NewUser("+1234567890")

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test duplicate phone number
	err = repo.Create(models.NewUser("+1234567890"))
	if err == nil {
		t.Error("Expected error for duplicate phone number")
	}

	// Verify error type
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected DomainError")
	}
	if domainErr.Code != "USER_ALREADY_EXISTS" {
		t.Errorf("Expected error code USER_ALREADY_EXISTS, got %s", domainErr.Code)
	}
}

func TestInMemoryUserRepository_GetByPhoneNumber(t *testing.T) {
	repo := NewUserRepository()
	user := models.NewUser("+1234567890")

	// Create user
	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user by phone number
	retrievedUser, err := repo.GetByPhoneNumber("+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}

	// Test non-existent user
	_, err = repo.GetByPhoneNumber("+9999999999")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected DomainError")
	}
	if domainErr.Code != "USER_NOT_FOUND" {
		t.Errorf("Expected error code USER_NOT_FOUND, got %s", domainErr.Code)
	}
}

func TestInMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewUserRepository()
	user := models.NewUser("+1234567890")

	// Create user
	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user by ID
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}

	// Test non-existent user ID
	_, err = repo.GetByID("non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent user ID")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected DomainError")
	}
	if domainErr.Code != "USER_NOT_FOUND" {
		t.Errorf("Expected error code USER_NOT_FOUND, got %s", domainErr.Code)
	}
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	repo := NewUserRepository()
	user := models.NewUser("+1234567890")

	// Create user
	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.PhoneNumber = "+9876543210"
	err = repo.Update(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if retrievedUser.PhoneNumber != "+9876543210" {
		t.Errorf("Expected updated phone number +9876543210, got %s", retrievedUser.PhoneNumber)
	}

	// Test update non-existent user
	nonExistentUser := models.NewUser("+1111111111")
	nonExistentUser.ID = "non-existent-id"
	err = repo.Update(nonExistentUser)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected DomainError")
	}
	if domainErr.Code != "USER_NOT_FOUND" {
		t.Errorf("Expected error code USER_NOT_FOUND, got %s", domainErr.Code)
	}
}

func TestInMemoryUserRepository_GetAll(t *testing.T) {
	repo := NewUserRepository()

	// Create multiple users
	users := []*models.User{
		models.NewUser("+1111111111"),
		models.NewUser("+2222222222"),
		models.NewUser("+3333333333"),
		models.NewUser("+4444444444"),
		models.NewUser("+5555555555"),
	}

	for _, user := range users {
		err := repo.Create(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	// Test get all users
	allUsers, total, err := repo.GetAll(1, 10, "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(allUsers) != 5 {
		t.Errorf("Expected 5 users, got %d", len(allUsers))
	}

	// Test pagination
	paginatedUsers, total, err := repo.GetAll(1, 2, "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(paginatedUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(paginatedUsers))
	}

	// Test second page
	secondPageUsers, total, err := repo.GetAll(2, 2, "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(secondPageUsers) != 2 {
		t.Errorf("Expected 2 users on second page, got %d", len(secondPageUsers))
	}

	// Test search
	searchResults, total, err := repo.GetAll(1, 10, "+1111111111")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if total != 1 {
		t.Errorf("Expected total 1 for search, got %d", total)
	}

	if len(searchResults) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(searchResults))
	}

	if searchResults[0].PhoneNumber != "+1111111111" {
		t.Errorf("Expected phone number +1111111111, got %s", searchResults[0].PhoneNumber)
	}

	// Test empty page
	emptyPageUsers, total, err := repo.GetAll(10, 10, "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(emptyPageUsers) != 0 {
		t.Errorf("Expected 0 users on empty page, got %d", len(emptyPageUsers))
	}
}

func TestInMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewUserRepository()
	done := make(chan bool, 10)

	// Test concurrent reads
	for i := 0; i < 5; i++ {
		go func() {
			_, _, _ = repo.GetAll(1, 10, "")
			done <- true
		}()
	}

	// Test concurrent writes
	for i := 0; i < 5; i++ {
		go func(index int) {
			user := models.NewUser("+1234567890")
			user.PhoneNumber = "+1234567890"
			_ = repo.Create(user)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify no panic occurred
	users, total, err := repo.GetAll(1, 100, "")
	if err != nil {
		t.Errorf("Expected no error after concurrent access, got %v", err)
	}

	if total < 0 {
		t.Error("Expected non-negative total after concurrent access")
	}

	if len(users) < 0 {
		t.Error("Expected non-negative user count after concurrent access")
	}
}

func TestInMemoryUserRepository_EdgeCases(t *testing.T) {
	repo := NewUserRepository()

	// Test with empty phone number
	user := models.NewUser("")
	user.PhoneNumber = ""
	err := repo.Create(user)
	if err != nil {
		t.Errorf("Expected no error for empty phone number, got %v", err)
	}

	// Test with very long phone number
	longPhoneUser := models.NewUser("+12345678901234567890")
	err = repo.Create(longPhoneUser)
	if err != nil {
		t.Errorf("Expected no error for long phone number, got %v", err)
	}

	// Test pagination edge cases - repository uses defaults for invalid values
	_, total, err := repo.GetAll(0, 0, "")
	if err != nil {
		t.Errorf("Expected no error for invalid pagination (should use defaults), got %v", err)
	}
	if total != 2 { // We created 2 users in edge cases
		t.Errorf("Expected total 2, got %d", total)
	}

	_, total, err = repo.GetAll(-1, -1, "")
	if err != nil {
		t.Errorf("Expected no error for negative pagination (should use defaults), got %v", err)
	}
	if total != 2 { // We created 2 users in edge cases
		t.Errorf("Expected total 2, got %d", total)
	}
}
