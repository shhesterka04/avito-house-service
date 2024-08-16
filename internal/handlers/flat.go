package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shhesterka04/house-service/internal/middleware"
	"github.com/shhesterka04/house-service/internal/repository"
)

type CreateFlatRequest struct {
	HouseID int `json:"house_id"`
	Number  int `json:"number"`
	Rooms   int `json:"rooms"`
	Price   int `json:"price"`
}

type CreateFlatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Status  string `json:"status"`
	Number  int    `json:"number"`
	Rooms   int    `json:"rooms"`
	Price   int    `json:"price"`
}

type UpdateFlatRequest struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type UpdateFlatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Status  string `json:"status"`
	Number  int    `json:"number"`
	Rooms   int    `json:"rooms"`
	Price   int    `json:"price"`
}

type GetFlatsByHouseIDRequest struct {
	HouseID int `json:"house_id"`
}

type GetFlatsByHouseIDResponse struct {
	Flats []*repository.Flat `json:"flats"`
}

type FlatRepo interface {
	CreateFlat(ctx context.Context, flat *repository.Flat) (*repository.Flat, error)
	UpdateFlat(ctx context.Context, flat *repository.Flat) (*repository.Flat, error)
	GetFlatByHouseID(ctx context.Context, houseID int, userType string) ([]*repository.Flat, error)
	GetFlatByID(ctx context.Context, id int) (*repository.Flat, error)
}

type HouseFlatRepo interface {
	UpdateHouse(ctx context.Context, id int, updAt time.Time) (*repository.House, error)
}

type FlatService struct {
	flatRepo      FlatRepo
	houseFlatRepo HouseFlatRepo
}

func NewFlatService(flatRepo FlatRepo, houseFlatRepo HouseFlatRepo) *FlatService {
	return &FlatService{
		flatRepo:      flatRepo,
		houseFlatRepo: houseFlatRepo,
	}
}

func (s *FlatService) CreateFlat(w http.ResponseWriter, r *http.Request) {
	var req CreateFlatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	flat := &repository.Flat{
		HouseId: req.HouseID,
		Number:  req.Number,
		Rooms:   req.Rooms,
		Price:   req.Price,
	}

	createdFlat, err := s.flatRepo.CreateFlat(r.Context(), flat)
	if err != nil {
		http.Error(w, "Failed to create flat", http.StatusInternalServerError)
		return
	}

	if _, err = s.houseFlatRepo.UpdateHouse(r.Context(), req.HouseID, time.Now()); err != nil {
		http.Error(w, "Failed to update house last added date", http.StatusInternalServerError)
		return
	}

	res := CreateFlatResponse{
		ID:      createdFlat.ID,
		HouseID: createdFlat.HouseId,
		Status:  createdFlat.Status,
		Number:  createdFlat.Number,
		Rooms:   createdFlat.Rooms,
		Price:   createdFlat.Price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *FlatService) UpdateFlat(w http.ResponseWriter, r *http.Request) {
	var req UpdateFlatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	validStatuses := map[string]struct{}{
		"created":       {},
		"approved":      {},
		"declined":      {},
		"on moderation": {},
	}

	if _, ok := validStatuses[req.Status]; !ok {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	flat, err := s.flatRepo.GetFlatByID(r.Context(), req.ID)
	if err != nil {
		http.Error(w, "Flat not found", http.StatusNotFound)
		return
	}

	flat.Status = req.Status

	updatedFlat, err := s.flatRepo.UpdateFlat(r.Context(), flat)
	if err != nil {
		http.Error(w, "Failed to update flat", http.StatusInternalServerError)
		return
	}

	res := UpdateFlatResponse{
		ID:      updatedFlat.ID,
		HouseID: updatedFlat.HouseId,
		Status:  updatedFlat.Status,
		Number:  updatedFlat.Number,
		Rooms:   updatedFlat.Rooms,
		Price:   updatedFlat.Price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *FlatService) GetFlatsByHouseID(w http.ResponseWriter, r *http.Request) {
	houseIDStr := strings.TrimPrefix(r.URL.Path, "/house/")
	houseID, err := strconv.Atoi(houseIDStr)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userType, valid := middleware.ValidTokens[token]
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	flats, err := s.flatRepo.GetFlatByHouseID(r.Context(), houseID, userType)
	if err != nil {
		http.Error(w, "Failed to get flats", http.StatusInternalServerError)
		return
	}

	res := GetFlatsByHouseIDResponse{
		Flats: flats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
