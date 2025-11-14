package repository

import (
	"context"
	"fmt"
	"qasynda/internal/models"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProviderRepository interface {
	Create(ctx context.Context, provider *models.ServiceProvider) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceProvider, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.ServiceProvider, error)
	Update(ctx context.Context, provider *models.ServiceProvider) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *models.ProviderFilter) ([]*models.ServiceProvider, int, error)
	AddService(ctx context.Context, providerID, serviceID uuid.UUID) error
	RemoveService(ctx context.Context, providerID, serviceID uuid.UUID) error
	GetProviderServices(ctx context.Context, providerID uuid.UUID) ([]*models.Service, error)
}

type providerRepository struct {
	db *sqlx.DB
}

func NewProviderRepository(db *sqlx.DB) ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) Create(ctx context.Context, provider *models.ServiceProvider) error {
	query := `
		INSERT INTO service_providers (id, user_id, bio, hourly_rate, experience_years, location, is_available, rating, total_reviews, created_at, updated_at)
		VALUES (:id, :user_id, :bio, :hourly_rate, :experience_years, :location, :is_available, :rating, :total_reviews, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, provider)
	return err
}

func (r *providerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceProvider, error) {
	var provider models.ServiceProvider
	query := `
		SELECT id, user_id, bio, hourly_rate, experience_years, location, is_available, rating, total_reviews, created_at, updated_at
		FROM service_providers WHERE id = $1
	`
	err := r.db.GetContext(ctx, &provider, query, id)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *providerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.ServiceProvider, error) {
	var provider models.ServiceProvider
	query := `
		SELECT id, user_id, bio, hourly_rate, experience_years, location, is_available, rating, total_reviews, created_at, updated_at
		FROM service_providers WHERE user_id = $1
	`
	err := r.db.GetContext(ctx, &provider, query, userID)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *providerRepository) Update(ctx context.Context, provider *models.ServiceProvider) error {
	query := `
		UPDATE service_providers 
		SET bio = :bio, hourly_rate = :hourly_rate, experience_years = :experience_years, 
		    location = :location, is_available = :is_available, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, provider)
	return err
}

func (r *providerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM service_providers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *providerRepository) List(ctx context.Context, filter *models.ProviderFilter) ([]*models.ServiceProvider, int, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	baseQuery := `
		SELECT DISTINCT sp.id, sp.user_id, sp.bio, sp.hourly_rate, sp.experience_years, 
		       sp.location, sp.is_available, sp.rating, sp.total_reviews, 
		       sp.created_at, sp.updated_at
		FROM service_providers sp
		LEFT JOIN provider_services ps ON sp.id = ps.provider_id
		WHERE 1=1
	`

	if filter.ServiceID != nil {
		conditions = append(conditions, fmt.Sprintf("ps.service_id = $%d", argPos))
		args = append(args, *filter.ServiceID)
		argPos++
	}

	if filter.City != nil && *filter.City != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(sp.location) LIKE LOWER($%d)", argPos))
		args = append(args, "%"+*filter.City+"%")
		argPos++
	}

	if filter.MinRating != nil {
		conditions = append(conditions, fmt.Sprintf("sp.rating >= $%d", argPos))
		args = append(args, *filter.MinRating)
		argPos++
	}

	if filter.IsAvailable != nil {
		conditions = append(conditions, fmt.Sprintf("sp.is_available = $%d", argPos))
		args = append(args, *filter.IsAvailable)
		argPos++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := strings.Replace(baseQuery, "SELECT DISTINCT sp.id, sp.user_id", "SELECT COUNT(DISTINCT sp.id)", 1)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Sorting
	sortBy := "rating"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	baseQuery += fmt.Sprintf(" ORDER BY sp.%s %s", sortBy, sortOrder)

	// Pagination
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.Limit, offset)

	var providers []*models.ServiceProvider
	err = r.db.SelectContext(ctx, &providers, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return providers, total, nil
}

func (r *providerRepository) AddService(ctx context.Context, providerID, serviceID uuid.UUID) error {
	query := `INSERT INTO provider_services (provider_id, service_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, providerID, serviceID)
	return err
}

func (r *providerRepository) RemoveService(ctx context.Context, providerID, serviceID uuid.UUID) error {
	query := `DELETE FROM provider_services WHERE provider_id = $1 AND service_id = $2`
	_, err := r.db.ExecContext(ctx, query, providerID, serviceID)
	return err
}

func (r *providerRepository) GetProviderServices(ctx context.Context, providerID uuid.UUID) ([]*models.Service, error) {
	var services []*models.Service
	query := `
		SELECT s.id, s.name, s.description, s.icon_url, s.created_at, s.updated_at
		FROM services s
		INNER JOIN provider_services ps ON s.id = ps.service_id
		WHERE ps.provider_id = $1
	`
	err := r.db.SelectContext(ctx, &services, query, providerID)
	return services, err
}

