//go:generate mockgen -source ./house.go -destination=./mocks/house.go -package=mocks
package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/dto"
)

type HouseRepo interface {
	CreateHouse(ctx context.Context, house *dto.House) (*dto.House, error)
}

type HouseService struct {
	houseRepo HouseRepo
}

func NewHouseService(houseRepo HouseRepo) *HouseService {
	return &HouseService{houseRepo: houseRepo}
}

func (s *HouseService) CreateHouse(ctx context.Context, req dto.PostHouseCreateJSONRequestBody) (*dto.House, error) {
	house := &dto.House{
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
	}

	if !validateHouseRequest(*house) {
		return nil, errors.New("invalid request")
	}

	house, err := s.houseRepo.CreateHouse(ctx, house)
	if err != nil {
		return nil, errors.Wrap(err, "create house")
	}

	return house, nil
}

func validateHouseRequest(h dto.House) bool {
	if h.Address == "" {
		return false
	}

	if h.Year <= 0 && h.Year > time.Now().Year() {
		return false
	}

	return true
}
