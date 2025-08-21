package repository

import (
	"sync"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(user *models.User) error
	GetAll(page, limit int, search string) ([]*models.User, int, error)
}

type InMemoryUserRepository struct {
	users      map[string]*models.User
	phoneIndex map[string]string // phone number -> user ID
	mutex      sync.RWMutex
}

func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users:      make(map[string]*models.User),
		phoneIndex: make(map[string]string),
	}
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if phone number already exists
	if _, exists := r.phoneIndex[user.PhoneNumber]; exists {
		return errors.ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	r.phoneIndex[user.PhoneNumber] = user.ID
	return nil
}

func (r *InMemoryUserRepository) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, exists := r.phoneIndex[phoneNumber]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetAll(page, limit int, search string) ([]*models.User, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var filteredUsers []*models.User

	// Filter users based on search criteria
	for _, user := range r.users {
		if search == "" || user.PhoneNumber == search {
			filteredUsers = append(filteredUsers, user)
		}
	}

	total := len(filteredUsers)

	// Calculate pagination
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []*models.User{}, total, nil
	}

	if end > total {
		end = total
	}

	return filteredUsers[start:end], total, nil
}
