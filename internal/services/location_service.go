package services

import (
	"context"
	"strings"

	"github.com/ghana-location-api/internal/models"
	"github.com/ghana-location-api/internal/repositories"
	"github.com/ghana-location-api/internal/errors"
)

type LocationService struct {
	countryRepo      *repositories.CountryRepository
	regionRepo       *repositories.RegionRepository
	districtRepo     *repositories.DistrictRepository
	constituencyRepo *repositories.ConstituencyRepository
	cityRepo         *repositories.CityRepository
}

func NewLocationService(
	countryRepo *repositories.CountryRepository,
	regionRepo *repositories.RegionRepository,
	districtRepo *repositories.DistrictRepository,
	constituencyRepo *repositories.ConstituencyRepository,
	cityRepo *repositories.CityRepository,
) *LocationService {
	return &LocationService{
		countryRepo:      countryRepo,
		regionRepo:       regionRepo,
		districtRepo:     districtRepo,
		constituencyRepo: constituencyRepo,
		cityRepo:         cityRepo,
	}
}

func (s *LocationService) validateSlug(slug string) error {
	if slug == "" || strings.Contains(slug, " ") || strings.ContainsAny(slug, "!@#$%^&*()") {
		return errors.ErrInvalidSlug
	}
	return nil
}

// Country methods
func (s *LocationService) GetAllCountries(ctx context.Context) ([]models.Country, error) {
	return s.countryRepo.GetAll(ctx)
}

func (s *LocationService) GetCountryByCode(ctx context.Context, code string) (*models.Country, error) {
	country, err := s.countryRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, errors.ErrNotFound
	}
	return country, nil
}

// Region methods
func (s *LocationService) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	return s.regionRepo.GetAll(ctx)
}

func (s *LocationService) GetRegionBySlug(ctx context.Context, slug string) (*models.Region, error) {
	if err := s.validateSlug(slug); err != nil {
		return nil, err
	}
	region, err := s.regionRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, errors.ErrNotFound
	}
	return region, nil
}

func (s *LocationService) GetDistrictsByRegionSlug(ctx context.Context, regionSlug string) ([]models.District, error) {
	if err := s.validateSlug(regionSlug); err != nil {
		return nil, err
	}
	return s.districtRepo.GetByRegionSlug(ctx, regionSlug)
}

// District methods
func (s *LocationService) GetDistrictBySlug(ctx context.Context, slug string) (*models.District, error) {
	if err := s.validateSlug(slug); err != nil {
		return nil, err
	}
	district, err := s.districtRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if district == nil {
		return nil, errors.ErrNotFound
	}
	return district, nil
}

func (s *LocationService) GetConstituenciesByDistrictSlug(ctx context.Context, districtSlug string) ([]models.Constituency, error) {
	if err := s.validateSlug(districtSlug); err != nil {
		return nil, err
	}
	return s.constituencyRepo.GetByDistrictSlug(ctx, districtSlug)
}

// Constituency methods
func (s *LocationService) GetConstituencyBySlug(ctx context.Context, slug string) (*models.Constituency, error) {
	if err := s.validateSlug(slug); err != nil {
		return nil, err
	}
	constituency, err := s.constituencyRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if constituency == nil {
		return nil, errors.ErrNotFound
	}
	return constituency, nil
}

// City methods
func (s *LocationService) GetCitiesByDistrictSlug(ctx context.Context, districtSlug string) ([]models.City, error) {
	if err := s.validateSlug(districtSlug); err != nil {
		return nil, err
	}
	return s.cityRepo.GetByDistrictSlug(ctx, districtSlug)
}
