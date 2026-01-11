package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ghana-location-api/internal/config"
	"github.com/ghana-location-api/internal/handlers"
	"github.com/ghana-location-api/internal/repositories"
	"github.com/ghana-location-api/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Initialize repositories
	countryRepo := repositories.NewCountryRepository(pool)
	regionRepo := repositories.NewRegionRepository(pool)
	districtRepo := repositories.NewDistrictRepository(pool)
	constituencyRepo := repositories.NewConstituencyRepository(pool)
	cityRepo := repositories.NewCityRepository(pool)

	// Initialize services
	locationService := services.NewLocationService(
		countryRepo,
		regionRepo,
		districtRepo,
		constituencyRepo,
		cityRepo,
	)

	// Initialize handlers
	countryHandler := handlers.NewCountryHandler(locationService)
	regionHandler := handlers.NewRegionHandler(locationService)
	districtHandler := handlers.NewDistrictHandler(locationService)
	constituencyHandler := handlers.NewConstituencyHandler(locationService)
	cityHandler := handlers.NewCityHandler(locationService)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Countries
		r.Get("/countries", countryHandler.GetAll)
		r.Get("/countries/{code}", countryHandler.GetByCode)

		// Regions
		r.Get("/regions", regionHandler.GetAll)
		r.Get("/regions/{slug}", regionHandler.GetBySlug)
		r.Get("/regions/{slug}/districts", regionHandler.GetDistricts)

		// Districts
		r.Get("/districts/{slug}", districtHandler.GetBySlug)
		r.Get("/districts/{slug}/constituencies", districtHandler.GetConstituencies)

		// Constituencies
		r.Get("/constituencies/{slug}", constituencyHandler.GetBySlug)

		// Cities
		r.Get("/cities", cityHandler.GetByDistrict)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
