package http

import (
	"application-service/internal/domain"
	usecase "application-service/internal/usecase/application"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IntegrationHandler struct {
	svc    *usecase.Service
	logger *zap.Logger
}

func NewIntegrationHandler(
	svc *usecase.Service,
	logger *zap.Logger,
) *IntegrationHandler {
	return &IntegrationHandler{
		svc:    svc,
		logger: logger,
	}
}

type createIntegrationApplicationRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Comment string `json:"comment"`
	Source  string `json:"source"`
}

func (h *IntegrationHandler) CreateApplication(c *gin.Context) {
	var req createIntegrationApplicationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(
			"invalid integration create application request",
			zap.Error(err),
		)

		Error(
			c,
			400,
			"INVALID_REQUEST",
			"Invalid request body",
		)
		return
	}

	application, err := h.svc.CreateApplication(
		c,
		req.Name,
		req.Phone,
		req.Comment,
		req.Source,
	)
	if err != nil {
		status, code, msg := MapError(err)
		Error(c, status, code, msg)
		return
	}

	Success(c, gin.H{
		"application": application,
	})
}

func (h *IntegrationHandler) ListApplications(c *gin.Context) {
	filter := domain.ApplicationFilter{}

	statusValue := c.Query("status")
	if statusValue != "" {
		status := domain.ApplicationStatus(statusValue)
		filter.Status = &status
	}

	applications, err := h.svc.ListApplications(c, filter)
	if err != nil {
		status, code, msg := MapError(err)
		Error(c, status, code, msg)
		return
	}

	Success(c, gin.H{
		"items": applications,
		"total": len(applications),
	})
}
