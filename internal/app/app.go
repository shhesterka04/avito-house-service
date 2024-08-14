package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/config"
	"github.com/shhesterka04/house-service/pkg/db"
	"github.com/shhesterka04/house-service/pkg/logger"
)

func Run(ctx context.Context) error {
	logger.Infof(ctx, "starting app")

	cfg, err := config.LoadConfig("")
	if err != nil {
		logger.Errorf(ctx, "config load error: %v", err)
		return errors.Wrap(err, "load config")
	}

	pgClient := db.NewClient(
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	_, err = pgClient.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}

	logger.Infof(ctx, "connected to database")

	defer pgClient.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	server := &http.Server{
		Addr: cfg.HostAddr,
	}

	if err = server.ListenAndServe(); err != nil {
		logger.Errorf(ctx, "Server start error: %v\n", err)
	}

	return nil
}
