package repositories

import (
	"context"

	"github.com/ghana-location-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CityRepository struct {
	pool *pgxpool.Pool
}

func NewCityRepository(pool *pgxpool.Pool) *CityRepository {
	return &CityRepository{pool: pool}
}

func (r *CityRepository) GetByDistrictSlug(ctx context.Context, districtSlug string) ([]models.City, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.district_id, c.name, c.lat, c.lng
		FROM cities c
		JOIN districts d ON c.district_id = d.id
		WHERE d.slug = $1
		ORDER BY c.name
	`, districtSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		var lat, lng *float64
		err := rows.Scan(&city.ID, &city.DistrictID, &city.Name, &lat, &lng)
		if err != nil {
			return nil, err
		}
		city.Lat = lat
		city.Lng = lng
		cities = append(cities, city)
	}

	return cities, rows.Err()
}
