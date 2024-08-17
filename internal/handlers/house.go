package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/service"
	"github.com/shhesterka04/house-service/pkg/logger"
)

type HouseHandler struct {
	houseService *service.HouseService
}

func NewHouseHandler(houseService *service.HouseService) *HouseHandler {
	return &HouseHandler{houseService: houseService}
}

func (h *HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	var req dto.PostHouseCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf(r.Context(), "Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	house, err := h.houseService.CreateHouse(r.Context(), req)
	if err != nil {
		logger.Errorf(r.Context(), "Error creating house: %v", err)
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(house)
}

func (h *HouseHandler) SubscribeToHouse(w http.ResponseWriter, r *http.Request) {}
