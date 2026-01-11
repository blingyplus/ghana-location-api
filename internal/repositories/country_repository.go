package repositories

import (
	"context"

	"github.com/ghana-location-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CountryRepository struct {
	pool *pgxpool.Pool
}

func NewCountryRepository(pool *pgxpool.Pool) *CountryRepository {
	return &CountryRepository{pool: pool}
}

func (r *CountryRepository) GetAll(ctx context.Context) ([]models.Country, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, code, name FROM countries ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(&country.ID, &country.Code, &country.Name); err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	return countries, rows.Err()
}

func (r *CountryRepository) GetByCode(ctx context.Context, code string) (*models.Country, error) {
	var country models.Country
	err := r.pool.QueryRow(ctx, "SELECT id, code, name FROM countries WHERE code = $1", code).
		Scan(&country.ID, &country.Code, &country.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}
