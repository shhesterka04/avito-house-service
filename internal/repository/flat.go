package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type DBFlat interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Flat struct {
	ID      int
	HouseId int
	Status  string
	Number  int
	Rooms   int
	Price   int
}

type FlatRepository struct {
	db DBFlat
}

func NewFlatRepository(db DBFlat) *FlatRepository {
	return &FlatRepository{db: db}
}

func (r *FlatRepository) CreateFlat(ctx context.Context, flat *Flat) (*Flat, error) {
	if _, err := r.db.Exec(ctx, "INSERT INTO flats (house_id, status, number, rooms, price) VALUES ($1, $2, $3, $4, $5)", flat.HouseId, "created", flat.Number, flat.Rooms, flat.Price); err != nil {
		return nil, errors.Wrap(err, "create flat")
	}

	return flat, nil
}

func (r *FlatRepository) UpdateFlat(ctx context.Context, flat *Flat) (*Flat, error) {
	if _, err := r.db.Exec(ctx, "UPDATE flats SET status = $1 WHERE id = $2", flat.Status, flat.ID); err != nil {
		return nil, errors.Wrap(err, "update flat")
	}

	return flat, nil
}

func (r *FlatRepository) GetFlatByID(ctx context.Context, id int) (*Flat, error) {
	row := r.db.QueryRow(ctx, "SELECT id, house_id, status, number, rooms, price FROM flats WHERE id = $1", id)
	flat := &Flat{}
	if err := row.Scan(&flat.ID, &flat.HouseId, &flat.Status, &flat.Number, &flat.Rooms, &flat.Price); err != nil {
		return nil, errors.Wrap(err, "get flat")
	}

	return flat, nil
}

func (r *FlatRepository) GetFlatByHouseID(ctx context.Context, houseID int, userType string) ([]*Flat, error) {
	var rows pgx.Rows
	var err error

	switch userType {
	case "Client":
		rows, err = r.db.Query(ctx, "SELECT id, house_id, status, number, rooms, price FROM flats WHERE house_id = $1 AND status = 'approved'", houseID)
	case "Moderator":
		rows, err = r.db.Query(ctx, "SELECT id, house_id, status, number, rooms, price FROM flats WHERE house_id = $1", houseID)
	default:
		return nil, fmt.Errorf("invalid user type")
	}
	if err != nil {
		return nil, errors.Wrap(err, "get flats")
	}
	defer rows.Close()

	var flats []*Flat
	for rows.Next() {
		var flat Flat
		if err = rows.Scan(&flat.ID, &flat.HouseId, &flat.Status, &flat.Number, &flat.Rooms, &flat.Price); err != nil {
			return nil, errors.Wrap(err, "scan flats")
		}
		flats = append(flats, &flat)
	}

	return flats, nil
}
