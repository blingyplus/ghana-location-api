package repositories

import (
	"context"

	"github.com/ghana-location-api/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DistrictRepository struct {
	pool *pgxpool.Pool
}

func NewDistrictRepository(pool *pgxpool.Pool) *DistrictRepository {
	return &DistrictRepository{pool: pool}
}

func (r *DistrictRepository) GetBySlug(ctx context.Context, slug string) (*models.District, error) {
	var district models.District
	err := r.pool.QueryRow(ctx, "SELECT id, region_id, name, slug, type, capital FROM districts WHERE slug = $1", slug).
		Scan(&district.ID, &district.RegionID, &district.Name, &district.Slug, &district.Type, &district.Capital)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &district, nil
}

func (r *DistrictRepository) GetByRegionSlug(ctx context.Context, regionSlug string) ([]models.District, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT d.id, d.region_id, d.name, d.slug, d.type, d.capital 
		FROM districts d
		JOIN regions r ON d.region_id = r.id
		WHERE r.slug = $1
		ORDER BY d.name
	`, regionSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var districts []models.District
	for rows.Next() {
		var district models.District
		if err := rows.Scan(&district.ID, &district.RegionID, &district.Name, &district.Slug, &district.Type, &district.Capital); err != nil {
			return nil, err
		}
		districts = append(districts, district)
	}

	return districts, rows.Err()
}
