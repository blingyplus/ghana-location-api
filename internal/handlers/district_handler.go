package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ghana-location-api/internal/errors"
	"github.com/ghana-location-api/internal/services"
)

type DistrictHandler struct {
	service *services.LocationService
}

func NewDistrictHandler(service *services.LocationService) *DistrictHandler {
	return &DistrictHandler{service: service}
}

func (h *DistrictHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		errors.WriteError(w, http.StatusBadRequest, "district slug is required")
		return
	}

	district, err := h.service.GetDistrictBySlug(r.Context(), slug)
	if err != nil {
		if err == errors.ErrNotFound {
			errors.WriteError(w, http.StatusNotFound, "district not found")
			return
		}
		if err == errors.ErrInvalidSlug {
			errors.WriteError(w, http.StatusBadRequest, "invalid slug format")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch district")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(district)
}

func (h *DistrictHandler) GetConstituencies(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		errors.WriteError(w, http.StatusBadRequest, "district slug is required")
		return
	}

	constituencies, err := h.service.GetConstituenciesByDistrictSlug(r.Context(), slug)
	if err != nil {
		if err == errors.ErrInvalidSlug {
			errors.WriteError(w, http.StatusBadRequest, "invalid slug format")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch constituencies")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(constituencies)
}
