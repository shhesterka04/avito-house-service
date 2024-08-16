package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type DBHouse interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type House struct {
	Id        int
	Address   string
	Year      int
	Developer string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HouseRepository struct {
	db DBHouse
}

func NewHouseRepository(db DBHouse) *HouseRepository {
	return &HouseRepository{db: db}
}

func (r *HouseRepository) CreateHouse(ctx context.Context, house *House) (*House, error) {
	if _, err := r.db.Exec(ctx, "INSERT INTO house (address, year, developer) VALUES ($1, $2, $3)", house.Address, house.Year, house.Developer); err != nil {
		return nil, errors.Wrap(err, "create house")
	}

	return house, nil
}

func (r *HouseRepository) UpdateHouse(ctx context.Context, id int, updAt time.Time) (*House, error) {
	if _, err := r.db.Exec(ctx, "UPDATE house SET updated_at = $1 WHERE id = $2", updAt, id); err != nil {
		return nil, errors.Wrap(err, "update house")
	}

	row := r.db.QueryRow(ctx, "SELECT id, address, year, developer, created_at, updated_at FROM house WHERE id = $1", id)
	house := &House{}
	if err := row.Scan(&house.Id, &house.Address, &house.Year, &house.Developer, &house.CreatedAt, &house.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "get house")
	}

	return house, nil
}
