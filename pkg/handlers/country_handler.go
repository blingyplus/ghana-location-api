package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ghana-location-api/pkg/errors"
	"github.com/ghana-location-api/pkg/services"
)

type CountryHandler struct {
	service *services.LocationService
}

func NewCountryHandler(service *services.LocationService) *CountryHandler {
	return &CountryHandler{service: service}
}

func (h *CountryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	countries, err := h.service.GetAllCountries(r.Context())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch countries")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(countries)
}

func (h *CountryHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		errors.WriteError(w, http.StatusBadRequest, "country code is required")
		return
	}

	country, err := h.service.GetCountryByCode(r.Context(), code)
	if err != nil {
		if err == errors.ErrNotFound {
			errors.WriteError(w, http.StatusNotFound, "country not found")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch country")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(country)
}
