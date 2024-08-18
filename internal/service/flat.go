//go:generate mockgen -source ./flat.go -destination=./mocks/flat.go -package=mocks
package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shhesterka04/house-service/internal/dto"
)

type FlatRepo interface {
	CreateFlat(ctx context.Context, flat *dto.DtoFlat) (*dto.DtoFlat, error)
	UpdateFlat(ctx context.Context, flat *dto.DtoFlat) (*dto.DtoFlat, error)
	GetFlatByHouseID(ctx context.Context, houseID int, userType string) ([]*dto.DtoFlat, error)
	GetFlatByID(ctx context.Context, id int) (*dto.DtoFlat, error)
}

type HouseFlatRepo interface {
	UpdateHouse(ctx context.Context, id int, updAt time.Time) (*dto.House, error)
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

func (s *FlatService) CreateFlat(ctx context.Context, req dto.CreateFlatRequest) (*dto.DtoFlat, error) {
	flat := &dto.DtoFlat{
		HouseId: req.HouseID,
		Number:  req.Number,
		Rooms:   req.Rooms,
		Price:   req.Price,
		Status:  string(dto.Created),
	}

	createdFlat, err := s.flatRepo.CreateFlat(ctx, flat)
	if err != nil {
		return nil, err
	}

	if _, err = s.houseFlatRepo.UpdateHouse(ctx, req.HouseID, time.Now()); err != nil {
		return nil, err
	}

	return createdFlat, nil
}

func (s *FlatService) UpdateFlat(ctx context.Context, req dto.PostFlatUpdateJSONRequestBody) (*dto.DtoFlat, error) {
	validStatuses := map[dto.Status]struct{}{
		dto.Created:      {},
		dto.Approved:     {},
		dto.Declined:     {},
		dto.OnModeration: {},
	}

	if _, ok := validStatuses[*req.Status]; !ok {
		return nil, errors.New("invalid status")
	}

	flat, err := s.flatRepo.GetFlatByID(ctx, req.Id)
	if err != nil {
		return nil, errors.New("DtoFlat not found")
	}

	flat.Status = string(*req.Status)

	updatedFlat, err := s.flatRepo.UpdateFlat(ctx, flat)
	if err != nil {
		return nil, err
	}

	_, err = s.houseFlatRepo.UpdateHouse(ctx, updatedFlat.HouseId, time.Now())
	if err != nil {
		return nil, err
	}

	return updatedFlat, nil
}

func (s *FlatService) GetFlatsByHouseID(ctx context.Context, houseIDStr, token string) ([]*dto.DtoFlat, error) {
	houseID, err := strconv.Atoi(houseIDStr)
	if err != nil {
		return nil, errors.New("invalid house ID")
	}

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	userType := claims.Subject
	flats, err := s.flatRepo.GetFlatByHouseID(ctx, houseID, userType)
	if err != nil {
		return nil, err
	}

	return flats, nil
}
