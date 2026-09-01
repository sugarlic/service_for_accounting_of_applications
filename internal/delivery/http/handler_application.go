package http

import (
	"application-service/internal/domain"
	usecase "application-service/internal/usecase/application"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApplicationHandler struct {
	svc    *usecase.Service
	logger *zap.Logger
}

func NewApplicationHandler(
	svc *usecase.Service,
	logger *zap.Logger,
) *ApplicationHandler {
	return &ApplicationHandler{
		svc:    svc,
		logger: logger,
	}
}

type createApplicationRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Comment string `json:"comment"`
	Source  string `json:"source"`
}

type updateApplicationStatusRequest struct {
	Status domain.ApplicationStatus `json:"status"`
}

func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req createApplicationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid create application request", zap.Error(err))
		Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
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

func (h *ApplicationHandler) ListApplications(c *gin.Context) {
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

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("invalid application id", zap.Error(err))
		Error(c, http.StatusBadRequest, "INVALID_APPLICATION_ID", "Invalid application id")
		return
	}

	application, err := h.svc.GetApplication(c, applicationID)
	if err != nil {
		status, code, msg := MapError(err)
		Error(c, status, code, msg)
		return
	}

	Success(c, gin.H{
		"application": application,
	})
}

func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("invalid application id", zap.Error(err))
		Error(c, http.StatusBadRequest, "INVALID_APPLICATION_ID", "Invalid application id")
		return
	}

	var req updateApplicationStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update application status request", zap.Error(err))
		Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	application, err := h.svc.UpdateApplicationStatus(
		c,
		applicationID,
		req.Status,
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
