package services

import (
	"strconv"

	"otp-auth-service/internal/models"
	"otp-auth-service/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(id string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *UserService) GetUsers(pageStr, limitStr, search string) ([]*models.UserResponse, int, error) {
	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	users, total, err := s.userRepo.GetAll(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []*models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, total, nil
}
