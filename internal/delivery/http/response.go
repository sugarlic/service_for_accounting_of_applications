package http

import (
	"application-service/internal/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	OK bool `json:"ok"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// Success sends a success response with data.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func MapError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrNameRequired):
		return http.StatusBadRequest, "NAME_REQUIRED", "Name is required"

	case errors.Is(err, domain.ErrPhoneRequired):
		return http.StatusBadRequest, "PHONE_REQUIRED", "Phone is required"

	case errors.Is(err, domain.ErrInvalidApplication):
		return http.StatusBadRequest, "INVALID_APPLICATION", "Invalid application"

	case errors.Is(err, domain.ErrInvalidStatus):
		return http.StatusBadRequest, "INVALID_STATUS", "Invalid application status"

	case errors.Is(err, domain.ErrApplicationNotFound):
		return http.StatusNotFound, "APPLICATION_NOT_FOUND", "Application not found"

	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	}
}
