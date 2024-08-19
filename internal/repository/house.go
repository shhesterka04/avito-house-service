//go:generate mockgen -source ./house.go -destination=./mocks/house_db.go -package=mocks
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/pkg/logger"
)

var ErrHouseExists = errors.New("house already exists")

type DBHouse interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type RowDBHouse interface {
	Scan(dest ...any) error
}

type HouseRepository struct {
	db DBHouse
}

func NewHouseRepository(db DBHouse) *HouseRepository {
	return &HouseRepository{db: db}
}

func (r *HouseRepository) CreateHouse(ctx context.Context, house *dto.House) (*dto.House, error) {
	var existingHouse dto.House
	err := r.db.QueryRow(ctx, "SELECT id FROM house WHERE address = $1", house.Address).Scan(&existingHouse.Id)
	if err == nil {
		return nil, errors.Wrap(ErrHouseExists, "house already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "query row")
	}

	if _, err = r.db.Exec(ctx, "INSERT INTO house (address, year, developer) VALUES ($1, $2, $3)", house.Address, house.Year, house.Developer); err != nil {
		return nil, errors.Wrap(err, "create house")
	}

	_ = r.db.QueryRow(ctx, "SELECT * FROM house WHERE address = $1", house.Address).Scan(&house.Id, &house.Address, &house.Year, &house.Developer, &house.CreatedAt, &house.UpdateAt)

	return house, nil
}

func (r *HouseRepository) UpdateHouse(ctx context.Context, id int, updAt time.Time) (*dto.House, error) {
	if _, err := r.db.Exec(ctx, "UPDATE house SET updated_at = $1 WHERE id = $2", updAt, id); err != nil {
		return nil, errors.Wrap(err, "update house")
	}

	row := r.db.QueryRow(ctx, "SELECT id, address, year, developer, created_at, updated_at FROM house WHERE id = $1", id)
	house := &dto.House{}
	if err := row.Scan(&house.Id, &house.Address, &house.Year, &house.Developer, &house.CreatedAt, &house.UpdateAt); err != nil {
		return nil, errors.Wrap(err, "get house")
	}

	logger.Infof(ctx, "house updated: %v", house)

	return house, nil
}
