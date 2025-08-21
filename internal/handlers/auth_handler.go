package handlers

import (
	"net/http"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/models"
	"otp-auth-service/internal/services"
	"otp-auth-service/internal/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RequestOTP godoc
// @Summary Request OTP for authentication
// @Description Generate and send OTP to the provided phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RequestOTPRequest true "Phone number"
// @Success 200 {object} models.RequestOTPResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Router /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req models.RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errors.ErrInvalidRequest.HTTPStatus, gin.H{
			"error": errors.ErrInvalidRequest.WithDetails(err.Error()),
		})
		return
	}

	// Validate phone number
	if err := validation.ValidateRequestOTP(req.PhoneNumber); err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	response, err := h.authService.RequestOTP(req.PhoneNumber)
	if err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTP godoc
// @Summary Verify OTP and authenticate user
// @Description Verify OTP and return JWT token for authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Phone number and OTP"
// @Success 200 {object} models.VerifyOTPResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			errors.ErrInvalidRequest.HTTPStatus,
			gin.H{
				"error": errors.ErrInvalidRequest.WithDetails(err.Error()),
			})
		return
	}

	// Validate phone number and OTP
	if err := validation.ValidateVerifyOTP(req.PhoneNumber, req.OTP); err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	response, err := h.authService.VerifyOTP(req.PhoneNumber, req.OTP)
	if err != nil {
		domainErr := errors.GetDomainError(err)
		c.JSON(domainErr.HTTPStatus, gin.H{
			"error": domainErr,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
