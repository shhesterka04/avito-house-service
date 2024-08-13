package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/logger"
)

type Database struct {
	cluster *pgxpool.Pool
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
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool connect")
	}

	c.pgpool = pool
	return &Database{cluster: pool}, nil
}

func (c *Client) Close() {
	if c.pgpool != nil {
		c.pgpool.Close()
	}
}

func (c *Client) generateDsn() string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", c.dbHost, c.dbPort, c.dbUser, c.dbUserPass, c.dbName)
}
