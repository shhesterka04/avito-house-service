package db

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/shhesterka04/house-service/pkg/logger"
)

type Database struct {
	Cluster *pgxpool.Pool
}

type Client struct {
	pgpool     *pgxpool.Pool
	dbName     string
	dbUser     string
	dbUserPass string
	dbHost     string
	dbPort     int
}

func NewClient(dbName, dbUser, dbUserPass, dbHost string, dbPort int) *Client {
	return &Client{
		dbName:     dbName,
		dbUser:     dbUser,
		dbUserPass: dbUserPass,
		dbHost:     dbHost,
		dbPort:     dbPort,
	}
}

func (c *Client) Connect(ctx context.Context) (*Database, error) {
	dsn := c.generateDsn()
	logger.Infof(ctx, "connecting to database: %v", dsn)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool connect")
	}

	c.pgpool = pool
	return &Database{Cluster: pool}, nil
}

func (c *Client) MigrateUp(migrationDir string) error {
	db := stdlib.OpenDBFromPool(c.pgpool)
	if err := goose.Up(db, migrationDir); err != nil {
		return errors.Wrap(err, "goose up")
	}

	return nil
}

func (c *Client) MigrateDown(migrationDir string) error {
	db := stdlib.OpenDBFromPool(c.pgpool)
	if err := goose.Down(db, migrationDir); err != nil {
		return errors.Wrap(err, "goose up")
	}

	return nil
}

func (c *Client) ResetMigrations(migrationDir string) error {
	db := stdlib.OpenDBFromPool(c.pgpool)
	if err := goose.Reset(db, migrationDir); err != nil {
		return errors.Wrap(err, "goose reset")
	}
	return nil
}

func (c *Client) Close() {
	if c.pgpool != nil {
		c.pgpool.Close()
	}
}

func (c *Client) generateDsn() string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", c.dbHost, c.dbPort, c.dbUser, c.dbUserPass, c.dbName)
}
