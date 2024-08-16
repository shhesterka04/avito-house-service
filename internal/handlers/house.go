package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shhesterka04/house-service/internal/repository"
)

type CreateHouseRequest struct {
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
}

type CreateHouseResponse struct {
	ID        int       `json:"id"`
	Address   string    `json:"address"`
	Year      int       `json:"year"`
	Developer string    `json:"developer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HouseRepo interface {
	CreateHouse(ctx context.Context, house *repository.House) (*repository.House, error)
}

type HouseService struct {
	houseRepo HouseRepo
}

func NewHouseService(houseRepo HouseRepo) *HouseService {
	return &HouseService{houseRepo: houseRepo}
}

func (s *HouseService) CreateHouse(w http.ResponseWriter, r *http.Request) {
	var req CreateHouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	house := &repository.House{
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
	}

	house, err := s.houseRepo.CreateHouse(r.Context(), house)
	if err != nil {
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	res := CreateHouseResponse{
		ID:        house.Id,
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: house.CreatedAt,
		UpdatedAt: house.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *HouseService) SubscribeToHouse(w http.ResponseWriter, r *http.Request) {}
