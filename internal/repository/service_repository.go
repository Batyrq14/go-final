package repository

import (
	"context"
	"qasynda/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Service, error)
	List(ctx context.Context) ([]*models.Service, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type serviceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, service *models.Service) error {
	query := `
		INSERT INTO services (id, name, description, icon_url, created_at, updated_at)
		VALUES (:id, :name, :description, :icon_url, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, service)
	return err
}

func (r *serviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	var service models.Service
	query := `SELECT id, name, description, icon_url, created_at, updated_at FROM services WHERE id = $1`
	err := r.db.GetContext(ctx, &service, query, id)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) List(ctx context.Context) ([]*models.Service, error) {
	var services []*models.Service
	query := `SELECT id, name, description, icon_url, created_at, updated_at FROM services ORDER BY name`
	err := r.db.SelectContext(ctx, &services, query)
	return services, err
}

func (r *serviceRepository) Update(ctx context.Context, service *models.Service) error {
	query := `
		UPDATE services 
		SET name = :name, description = :description, icon_url = :icon_url, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, service)
	return err
}

func (r *serviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

