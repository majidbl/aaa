package middleware

import (
	"strings"

	"otp-auth-service/internal/errors"
	"otp-auth-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(errors.ErrMissingAuthHeader.HTTPStatus, gin.H{
				"error": errors.ErrMissingAuthHeader,
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(errors.ErrInvalidAuthFormat.HTTPStatus, gin.H{
				"error": errors.ErrInvalidAuthFormat,
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		authService := services.NewAuthService(nil, nil, jwtSecret)
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			domainErr := errors.GetDomainError(err)
			c.JSON(domainErr.HTTPStatus, gin.H{
				"error": domainErr,
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(errors.ErrInvalidToken.HTTPStatus, gin.H{
				"error": errors.ErrInvalidToken,
			})
			c.Abort()
			return
		}

		// Set user information in context
		userID, _ := claims["user_id"].(string)
		phoneNumber, _ := claims["phone_number"].(string)

		c.Set("user_id", userID)
		c.Set("phone_number", phoneNumber)

		c.Next()
	}
}
