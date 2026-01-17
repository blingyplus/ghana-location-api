package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ghana-location-api/pkg/errors"
	"github.com/ghana-location-api/pkg/services"
)

type CityHandler struct {
	service *services.LocationService
}

func NewCityHandler(service *services.LocationService) *CityHandler {
	return &CityHandler{service: service}
}

func (h *CityHandler) GetByDistrict(w http.ResponseWriter, r *http.Request) {
	districtSlug := r.URL.Query().Get("district")
	if districtSlug == "" {
		errors.WriteError(w, http.StatusBadRequest, "district parameter is required")
		return
	}

	cities, err := h.service.GetCitiesByDistrictSlug(r.Context(), districtSlug)
	if err != nil {
		if err == errors.ErrInvalidSlug {
			errors.WriteError(w, http.StatusBadRequest, "invalid slug format")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch cities")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cities)
}
