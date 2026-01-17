package repositories

import (
	"context"

	"github.com/ghana-location-api/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegionRepository struct {
	pool *pgxpool.Pool
}

func NewRegionRepository(pool *pgxpool.Pool) *RegionRepository {
	return &RegionRepository{pool: pool}
}

func (r *RegionRepository) GetAll(ctx context.Context) ([]models.Region, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, country_id, name, slug, capital FROM regions ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []models.Region
	for rows.Next() {
		var region models.Region
		if err := rows.Scan(&region.ID, &region.CountryID, &region.Name, &region.Slug, &region.Capital); err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}

	return regions, rows.Err()
}

func (r *RegionRepository) GetBySlug(ctx context.Context, slug string) (*models.Region, error) {
	var region models.Region
	err := r.pool.QueryRow(ctx, "SELECT id, country_id, name, slug, capital FROM regions WHERE slug = $1", slug).
		Scan(&region.ID, &region.CountryID, &region.Name, &region.Slug, &region.Capital)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &region, nil
}
