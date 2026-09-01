package application

import (
	"application-service/internal/domain"
	"context"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	applicationRepo ApplicationRepository
	logger          *zap.Logger
}

func NewService(
	applicationRepo ApplicationRepository,
	logger *zap.Logger,
) *Service {
	return &Service{
		applicationRepo: applicationRepo,
		logger:          logger,
	}
}

func (s *Service) CreateApplication(
	ctx context.Context,
	name string,
	phone string,
	comment string,
	source string,
) (*domain.Application, error) {
	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	comment = strings.TrimSpace(comment)
	source = strings.TrimSpace(source)

	if name == "" {
		return nil, domain.ErrInvalidApplication
	}

	if phone == "" {
		return nil, domain.ErrInvalidApplication
	}

	application := &domain.Application{
		Name:    name,
		Phone:   phone,
		Comment: comment,
		Source:  source,
		Status:  domain.ApplicationStatusNew,
	}

	if err := s.applicationRepo.CreateApplication(ctx, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (s *Service) ListApplications(
	ctx context.Context,
	filter domain.ApplicationFilter,
) ([]domain.Application, error) {
	if filter.Status != nil && !filter.Status.IsValid() {
		return nil, domain.ErrInvalidStatus
	}

	applications, err := s.applicationRepo.ListApplications(ctx, filter)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (s *Service) GetApplication(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Application, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidApplication
	}

	application, err := s.applicationRepo.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (s *Service) UpdateApplicationStatus(
	ctx context.Context,
	id uuid.UUID,
	status domain.ApplicationStatus,
) (*domain.Application, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidApplication
	}

	if !status.IsValid() {
		return nil, domain.ErrInvalidStatus
	}

	if err := s.applicationRepo.UpdateApplicationStatus(ctx, id, status); err != nil {
		return nil, err
	}

	application, err := s.applicationRepo.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return application, nil
}
