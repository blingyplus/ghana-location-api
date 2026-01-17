package repositories

import (
	"context"

	"github.com/ghana-location-api/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConstituencyRepository struct {
	pool *pgxpool.Pool
}

func NewConstituencyRepository(pool *pgxpool.Pool) *ConstituencyRepository {
	return &ConstituencyRepository{pool: pool}
}

func (r *ConstituencyRepository) GetBySlug(ctx context.Context, slug string) (*models.Constituency, error) {
	var constituency models.Constituency
	err := r.pool.QueryRow(ctx, "SELECT id, district_id, name, slug FROM constituencies WHERE slug = $1", slug).
		Scan(&constituency.ID, &constituency.DistrictID, &constituency.Name, &constituency.Slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &constituency, nil
}

func (r *ConstituencyRepository) GetByDistrictSlug(ctx context.Context, districtSlug string) ([]models.Constituency, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.district_id, c.name, c.slug
		FROM constituencies c
		JOIN districts d ON c.district_id = d.id
		WHERE d.slug = $1
		ORDER BY c.name
	`, districtSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var constituencies []models.Constituency
	for rows.Next() {
		var constituency models.Constituency
		if err := rows.Scan(&constituency.ID, &constituency.DistrictID, &constituency.Name, &constituency.Slug); err != nil {
			return nil, err
		}
		constituencies = append(constituencies, constituency)
	}

	return constituencies, rows.Err()
}
