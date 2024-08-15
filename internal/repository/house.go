package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBHouse interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type House struct {
}

type HouseRepository struct {
	db DBHouse
}

func NewHouseRepository(db DBHouse) *HouseRepository {
	return &HouseRepository{db: db}
}

//create

//get

//updateHouse

//getFlatByHouseID
