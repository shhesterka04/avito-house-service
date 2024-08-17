package service

import (
	"context"

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

	house, err := s.houseRepo.CreateHouse(ctx, house)
	if err != nil {
		return nil, errors.Wrap(err, "create house")
	}

	return house, nil
}
