package postgres

import (
	"application-service/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepo struct {
	pool *pgxpool.Pool
}

func NewApplicationRepo(pool *pgxpool.Pool) *ApplicationRepo {
	return &ApplicationRepo{pool: pool}
}

func (r *ApplicationRepo) exec(ctx context.Context) DBExecutor {
	return getExecutor(ctx, r.pool)
}

func (r *ApplicationRepo) CreateApplication(
	ctx context.Context,
	a *domain.Application,
) error {
	const q = `
		INSERT INTO applications (
			name,
			phone,
			comment,
			source,
			status
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	exec := r.exec(ctx)

	if err := exec.QueryRow(
		ctx,
		q,
		a.Name,
		a.Phone,
		a.Comment,
		a.Source,
		a.Status,
	).Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create application: %w", err)
	}

	return nil
}

func (r *ApplicationRepo) GetApplicationByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Application, error) {
	const q = `
		SELECT
			id,
			name,
			phone,
			comment,
			source,
			status,
			created_at,
			updated_at
		FROM applications
		WHERE id = $1
	`

	exec := r.exec(ctx)

	var a domain.Application

	if err := exec.QueryRow(ctx, q, id).Scan(
		&a.ID,
		&a.Name,
		&a.Phone,
		&a.Comment,
		&a.Source,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrApplicationNotFound
		}

		return nil, fmt.Errorf("get application by id: %w", err)
	}

	return &a, nil
}

func (r *ApplicationRepo) ListApplications(
	ctx context.Context,
	filter domain.ApplicationFilter,
) ([]domain.Application, error) {
	query := `
		SELECT
			id,
			name,
			phone,
			comment,
			source,
			status,
			created_at,
			updated_at
		FROM applications
	`

	args := make([]any, 0, 1)

	if filter.Status != nil {
		query += ` WHERE status = $1`
		args = append(args, *filter.Status)
	}

	query += ` ORDER BY created_at DESC`

	exec := r.exec(ctx)

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}
	defer rows.Close()

	applications := make([]domain.Application, 0)

	for rows.Next() {
		var a domain.Application

		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Phone,
			&a.Comment,
			&a.Source,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan application: %w", err)
		}

		applications = append(applications, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applications: %w", err)
	}

	return applications, nil
}

func (r *ApplicationRepo) UpdateApplicationStatus(
	ctx context.Context,
	id uuid.UUID,
	status domain.ApplicationStatus,
) error {
	const q = `
		UPDATE applications
		SET
			status = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	exec := r.exec(ctx)

	tag, err := exec.Exec(ctx, q, id, status)
	if err != nil {
		return fmt.Errorf("update application status: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrApplicationNotFound
	}

	return nil
}
