package handlers

import (
	"net/http"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/services"
	"otp-auth-service/internal/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get list of users
// @Description Retrieve paginated list of users with optional search
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param search query string false "Search by phone number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	search := c.Query("search")

	// Validate query parameters
	if err := validation.ValidateGetUsers(page, limit, search); err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	users, total, err := h.userService.GetUsers(page, limit, search)
	if err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	// Validate user ID
	if err := validation.ValidateGetUser(userID); err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
