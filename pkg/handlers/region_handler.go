package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ghana-location-api/pkg/errors"
	"github.com/ghana-location-api/pkg/services"
)

type RegionHandler struct {
	service *services.LocationService
}

func NewRegionHandler(service *services.LocationService) *RegionHandler {
	return &RegionHandler{service: service}
}

func (h *RegionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	regions, err := h.service.GetAllRegions(r.Context())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch regions")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regions)
}

func (h *RegionHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		errors.WriteError(w, http.StatusBadRequest, "region slug is required")
		return
	}

	region, err := h.service.GetRegionBySlug(r.Context(), slug)
	if err != nil {
		if err == errors.ErrNotFound {
			errors.WriteError(w, http.StatusNotFound, "region not found")
			return
		}
		if err == errors.ErrInvalidSlug {
			errors.WriteError(w, http.StatusBadRequest, "invalid slug format")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch region")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(region)
}

func (h *RegionHandler) GetDistricts(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		errors.WriteError(w, http.StatusBadRequest, "region slug is required")
		return
	}

	districts, err := h.service.GetDistrictsByRegionSlug(r.Context(), slug)
	if err != nil {
		if err == errors.ErrInvalidSlug {
			errors.WriteError(w, http.StatusBadRequest, "invalid slug format")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch districts")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(districts)
}
