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

	var count int

	pool.QueryRow(ctx, "SELECT COUNT(*) FROM countries").Scan(&count)
	fmt.Printf("Countries: %d\n", count)

	pool.QueryRow(ctx, "SELECT COUNT(*) FROM regions").Scan(&count)
	fmt.Printf("Regions: %d\n", count)

	pool.QueryRow(ctx, "SELECT COUNT(*) FROM districts").Scan(&count)
	fmt.Printf("Districts: %d\n", count)

	pool.QueryRow(ctx, "SELECT COUNT(*) FROM constituencies").Scan(&count)
	fmt.Printf("Constituencies: %d\n", count)

	pool.QueryRow(ctx, "SELECT COUNT(*) FROM cities").Scan(&count)
	fmt.Printf("Cities: %d\n", count)
}
