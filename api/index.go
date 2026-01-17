package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ghana-location-api/internal/config"
	"github.com/ghana-location-api/internal/handlers"
	"github.com/ghana-location-api/internal/repositories"
	"github.com/ghana-location-api/internal/services"
)

var router http.Handler

func init() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic(err)
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

	router = r
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
