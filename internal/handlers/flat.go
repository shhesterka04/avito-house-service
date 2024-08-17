package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/service"
)

type FlatHandler struct {
	flatService *service.FlatService
}

func NewFlatHandler(flatService *service.FlatService) *FlatHandler {
	return &FlatHandler{flatService: flatService}
}

func (h *FlatHandler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFlatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdFlat, err := h.flatService.CreateFlat(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create flat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdFlat)
}

func (h *FlatHandler) UpdateFlat(w http.ResponseWriter, r *http.Request) {
	var req dto.PostFlatUpdateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedFlat, err := h.flatService.UpdateFlat(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedFlat)
}

func (h *FlatHandler) GetFlatsByHouseID(w http.ResponseWriter, r *http.Request) {
	houseIDStr := strings.TrimPrefix(r.URL.Path, "/house/")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	flats, err := h.flatService.GetFlatsByHouseID(r.Context(), houseIDStr, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flats)
}
