package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type CountryData struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type RegionData struct {
	Name    string  `json:"name"`
	Slug    string  `json:"slug"`
	Capital *string `json:"capital,omitempty"`
}

type DistrictData struct {
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	Type       string  `json:"type"`
	Capital    *string `json:"capital,omitempty"`
	RegionSlug string  `json:"region_slug"`
}

type ConstituencyData struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	RegionSlug  string  `json:"region_slug"`
	DistrictSlug *string `json:"district_slug,omitempty"`
}

type CityData struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Lat         *float64 `json:"lat,omitempty"`
	Lng         *float64 `json:"lng,omitempty"`
	DistrictSlug string  `json:"district_slug"`
}

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatalf("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to database successfully")

	// Seed countries
	fmt.Println("\nSeeding countries...")
	if err := seedCountries(ctx, pool); err != nil {
		log.Fatalf("failed to seed countries: %v", err)
	}
	fmt.Println("✓ Countries seeded")

	// Seed regions
	fmt.Println("\nSeeding regions...")
	regionMap, err := seedRegions(ctx, pool)
	if err != nil {
		log.Fatalf("failed to seed regions: %v", err)
	}
	fmt.Printf("✓ Regions seeded (%d regions)\n", len(regionMap))

	// Seed districts
	fmt.Println("\nSeeding districts...")
	districtMap, err := seedDistricts(ctx, pool, regionMap)
	if err != nil {
		log.Fatalf("failed to seed districts: %v", err)
	}
	fmt.Printf("✓ Districts seeded (%d districts)\n", len(districtMap))

	// Seed constituencies
	fmt.Println("\nSeeding constituencies...")
	if err := seedConstituencies(ctx, pool, districtMap); err != nil {
		log.Fatalf("failed to seed constituencies: %v", err)
	}
	fmt.Println("✓ Constituencies seeded")

	// Seed cities
	fmt.Println("\nSeeding cities...")
	if err := seedCities(ctx, pool, districtMap); err != nil {
		log.Fatalf("failed to seed cities: %v", err)
	}
	fmt.Println("✓ Cities seeded")

	fmt.Println("\n✓ Database seeding completed successfully!")
}

func seedCountries(ctx context.Context, pool *pgxpool.Pool) error {
	data, err := os.ReadFile("data/countries.json")
	if err != nil {
		return fmt.Errorf("failed to read countries.json: %w", err)
	}

	var countries []CountryData
	if err := json.Unmarshal(data, &countries); err != nil {
		return fmt.Errorf("failed to parse countries.json: %w", err)
	}

	for _, country := range countries {
		_, err := pool.Exec(ctx,
			"INSERT INTO countries (code, name) VALUES ($1, $2) ON CONFLICT (code) DO NOTHING",
			country.Code, country.Name,
		)
		if err != nil {
			return fmt.Errorf("failed to insert country %s: %w", country.Code, err)
		}
	}

	return nil
}

func seedRegions(ctx context.Context, pool *pgxpool.Pool) (map[string]string, error) {
	data, err := os.ReadFile("data/regions.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read regions.json: %w", err)
	}

	var regions []RegionData
	if err := json.Unmarshal(data, &regions); err != nil {
		return nil, fmt.Errorf("failed to parse regions.json: %w", err)
	}

	// Get country ID for Ghana
	var countryID string
	err = pool.QueryRow(ctx, "SELECT id FROM countries WHERE code = 'GH'").Scan(&countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get country ID for GH: %w", err)
	}

	regionMap := make(map[string]string)

	for _, region := range regions {
		var regionID string
		err := pool.QueryRow(ctx,
			`INSERT INTO regions (country_id, name, slug, capital) 
			 VALUES ($1, $2, $3, $4) 
			 ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, capital = EXCLUDED.capital
			 RETURNING id`,
			countryID, region.Name, region.Slug, region.Capital,
		).Scan(&regionID)

		if err != nil {
			return nil, fmt.Errorf("failed to insert region %s: %w", region.Slug, err)
		}

		regionMap[region.Slug] = regionID
	}

	return regionMap, nil
}

func seedDistricts(ctx context.Context, pool *pgxpool.Pool, regionMap map[string]string) (map[string]string, error) {
	data, err := os.ReadFile("data/districts.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read districts.json: %w", err)
	}

	var districts []DistrictData
	if err := json.Unmarshal(data, &districts); err != nil {
		return nil, fmt.Errorf("failed to parse districts.json: %w", err)
	}

	districtMap := make(map[string]string)

	for _, district := range districts {
		regionID, exists := regionMap[district.RegionSlug]
		if !exists {
			return nil, fmt.Errorf("region not found: %s", district.RegionSlug)
		}

		var districtID string
		err := pool.QueryRow(ctx,
			`INSERT INTO districts (region_id, name, slug, type, capital) 
			 VALUES ($1, $2, $3, $4, $5) 
			 ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, type = EXCLUDED.type, capital = EXCLUDED.capital
			 RETURNING id`,
			regionID, district.Name, district.Slug, district.Type, district.Capital,
		).Scan(&districtID)

		if err != nil {
			return nil, fmt.Errorf("failed to insert district %s: %w", district.Slug, err)
		}

		districtMap[district.Slug] = districtID
	}

	return districtMap, nil
}

func seedConstituencies(ctx context.Context, pool *pgxpool.Pool, districtMap map[string]string) error {
	data, err := os.ReadFile("data/constituencies.json")
	if err != nil {
		return fmt.Errorf("failed to read constituencies.json: %w", err)
	}

	var constituencies []ConstituencyData
	if err := json.Unmarshal(data, &constituencies); err != nil {
		return fmt.Errorf("failed to parse constituencies.json: %w", err)
	}

	for _, constituency := range constituencies {
		var districtID *string
		if constituency.DistrictSlug != nil {
			id, exists := districtMap[*constituency.DistrictSlug]
			if !exists {
				// Log warning but continue
				fmt.Printf("  ⚠ District not found for constituency %s: %s\n", constituency.Slug, *constituency.DistrictSlug)
			} else {
				districtID = &id
			}
		}

		_, err := pool.Exec(ctx,
			`INSERT INTO constituencies (district_id, name, slug) 
			 VALUES ($1, $2, $3) 
			 ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, district_id = EXCLUDED.district_id`,
			districtID, constituency.Name, constituency.Slug,
		)

		if err != nil {
			return fmt.Errorf("failed to insert constituency %s: %w", constituency.Slug, err)
		}
	}

	return nil
}

func seedCities(ctx context.Context, pool *pgxpool.Pool, districtMap map[string]string) error {
	data, err := os.ReadFile("data/cities.json")
	if err != nil {
		return fmt.Errorf("failed to read cities.json: %w", err)
	}

	var cities []CityData
	if err := json.Unmarshal(data, &cities); err != nil {
		return fmt.Errorf("failed to parse cities.json: %w", err)
	}

	for _, city := range cities {
		districtID, exists := districtMap[city.DistrictSlug]
		if !exists {
			// Log warning but continue
			fmt.Printf("  ⚠ District not found for city %s: %s\n", city.Name, city.DistrictSlug)
			continue
		}

		_, err := pool.Exec(ctx,
			`INSERT INTO cities (district_id, name, lat, lng) 
			 VALUES ($1, $2, $3, $4) 
			 ON CONFLICT (district_id, name) DO UPDATE SET lat = EXCLUDED.lat, lng = EXCLUDED.lng`,
			districtID, city.Name, city.Lat, city.Lng,
		)

		if err != nil {
			return fmt.Errorf("failed to insert city %s: %w", city.Name, err)
		}
	}

	return nil
}
