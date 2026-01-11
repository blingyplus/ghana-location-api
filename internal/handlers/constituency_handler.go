package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ghana-location-api/internal/errors"
	"github.com/ghana-location-api/internal/services"
)

type ConstituencyHandler struct {
	service *services.LocationService
}

func NewConstituencyHandler(service *services.LocationService) *ConstituencyHandler {
	return &ConstituencyHandler{service: service}
}

func (h *ConstituencyHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		errors.WriteError(w, http.StatusBadRequest, "constituency slug is required")
		return
	}

	constituency, err := h.service.GetConstituencyBySlug(r.Context(), slug)
	if err != nil {
		if err == errors.ErrNotFound {
			errors.WriteError(w, http.StatusNotFound, "constituency not found")
			return
		}
		if err == errors.ErrInvalidSlug {
			errors.WriteError(w, http.StatusBadRequest, "invalid slug format")
			return
		}
		errors.WriteError(w, http.StatusInternalServerError, "failed to fetch constituency")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(constituency)
}
