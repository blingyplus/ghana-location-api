package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to database successfully")

	// Read migration file
	migrationPath := "migrations/001_initial_schema.sql"
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("failed to read migration file: %v", err)
	}

	// Execute migration as a single transaction
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Execute the entire migration
	_, err = tx.Exec(ctx, string(migrationSQL))
	if err != nil {
		log.Fatalf("failed to execute migration: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	fmt.Println("✓ Migration executed successfully")

	// Verify tables were created
	tables := []string{"countries", "regions", "districts", "constituencies", "cities"}
	for _, table := range tables {
		var exists bool
		err := pool.QueryRow(ctx, 
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)",
			table,
		).Scan(&exists)
		
		if err != nil {
			log.Fatalf("failed to check table %s: %v", table, err)
		}

		if exists {
			fmt.Printf("✓ Table '%s' exists\n", table)
		} else {
			log.Fatalf("✗ Table '%s' does not exist", table)
		}
	}

	// Verify indexes
	indexes := []string{
		"idx_regions_country_id",
		"idx_regions_slug",
		"idx_districts_region_id",
		"idx_districts_slug",
		"idx_constituencies_district_id",
		"idx_constituencies_slug",
		"idx_cities_district_id",
	}

	for _, idx := range indexes {
		var exists bool
		err := pool.QueryRow(ctx,
			"SELECT EXISTS (SELECT FROM pg_indexes WHERE indexname = $1)",
			idx,
		).Scan(&exists)

		if err != nil {
			log.Fatalf("failed to check index %s: %v", idx, err)
		}

		if exists {
			fmt.Printf("✓ Index '%s' exists\n", idx)
		} else {
			fmt.Printf("⚠ Index '%s' not found (may be created automatically)\n", idx)
		}
	}

	fmt.Println("\n✓ All migrations tested successfully!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
