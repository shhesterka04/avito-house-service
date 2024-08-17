package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/dto"
)

var ErrFlatExists = errors.New("flat already exists")

type DBFlat interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type FlatRepository struct {
	db DBFlat
}

func NewFlatRepository(db DBFlat) *FlatRepository {
	return &FlatRepository{db: db}
}

func (r *FlatRepository) CreateFlat(ctx context.Context, flat *dto.DtoFlat) (*dto.DtoFlat, error) {
	var existingFlat dto.DtoFlat
	err := r.db.QueryRow(ctx, "SELECT id FROM flats WHERE house_id = $1 AND number = $2", flat.HouseId, flat.Number).Scan(&existingFlat.Id)
	if err == nil {
		return nil, errors.Wrap(ErrFlatExists, "flat already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "query row")
	}

	if _, err = r.db.Exec(ctx, "INSERT INTO flats (house_id, status, number, rooms, price) VALUES ($1, $2, $3, $4, $5)", flat.HouseId, flat.Status, flat.Number, flat.Rooms, flat.Price); err != nil {
		return nil, errors.Wrap(err, "create flat")
	}

	return flat, nil
}

func (r *FlatRepository) UpdateFlat(ctx context.Context, flat *dto.DtoFlat) (*dto.DtoFlat, error) {
	if _, err := r.db.Exec(ctx, "UPDATE flats SET status = $1 WHERE id = $2", flat.Status, flat.Id); err != nil {
		return nil, errors.Wrap(err, "update flat")
	}

	return flat, nil
}

func (r *FlatRepository) GetFlatByID(ctx context.Context, id int) (*dto.DtoFlat, error) {
	row := r.db.QueryRow(ctx, "SELECT id, house_id, status, number, rooms, price FROM flats WHERE id = $1", id)
	flat := &dto.DtoFlat{}
	if err := row.Scan(&flat.Id, &flat.HouseId, &flat.Status, &flat.Number, &flat.Rooms, &flat.Price); err != nil {
		return nil, errors.Wrap(err, "get flat")
	}

	return flat, nil
}

func (r *FlatRepository) GetFlatByHouseID(ctx context.Context, houseId int, userType string) ([]*dto.DtoFlat, error) {
	var rows pgx.Rows
	var err error

	switch userType {
	case string(dto.Client):
		rows, err = r.db.Query(ctx, "SELECT id, house_id, status, number, rooms, price FROM flats WHERE house_id = $1 AND status = 'approved'", houseId)
	case string(dto.Moderator):
		rows, err = r.db.Query(ctx, "SELECT id, house_id, status, number, rooms, price FROM flats WHERE house_id = $1 AND status != 'on moderation'", houseId)
	default:
		return nil, fmt.Errorf("invalid user type")
	}
	if err != nil {
		return nil, errors.Wrap(err, "get flats")
	}
	defer rows.Close()

	var flats []*dto.DtoFlat
	for rows.Next() {
		var flat dto.DtoFlat
		if err = rows.Scan(&flat.Id, &flat.HouseId, &flat.Status, &flat.Number, &flat.Rooms, &flat.Price); err != nil {
			return nil, errors.Wrap(err, "scan flats")
		}
		flats = append(flats, &flat)
	}

	return flats, nil
}
