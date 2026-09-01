package application

import (
	"application-service/internal/domain"
	"context"

	"github.com/google/uuid"
)

type ApplicationRepository interface {
	CreateApplication(ctx context.Context, application *domain.Application) error
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*domain.Application, error)
	ListApplications(ctx context.Context, filter domain.ApplicationFilter) ([]domain.Application, error)
	UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status domain.ApplicationStatus) error
}
