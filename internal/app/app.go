package app

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/config"
	"github.com/shhesterka04/house-service/internal/handlers"
	"github.com/shhesterka04/house-service/internal/middleware"
	"github.com/shhesterka04/house-service/internal/repository"
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

	dbConn, err := pgClient.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}
	logger.Infof(ctx, "connected to database")

	if err = pgClient.Migrate("/migrations"); err != nil {
		return errors.Wrap(err, "migrate")
	}
	defer pgClient.Close()

	userRepo := repository.NewUserRepository(dbConn.Cluster)
	userService := handlers.NewAuthService(userRepo)

	http.HandleFunc("POST /dummyLogin", userService.DummyLogin)
	http.HandleFunc("POST /register", userService.Register)

	protectedRoutes := http.NewServeMux()
	protectedRoutes.Handle("POST /house/create", middleware.AuthMiddleware("moderator")(http.HandlerFunc(handlers.CreateHouse)))
	protectedRoutes.Handle("GET /house/{id}", middleware.AuthMiddleware("client")(http.HandlerFunc(handlers.GetFlatsByHouseID)))
	protectedRoutes.Handle("POST /house/{id}/subscribe", middleware.AuthMiddleware("moderator")(http.HandlerFunc(handlers.SubscribeToHouse)))
	protectedRoutes.Handle("POST /flat/create", middleware.AuthMiddleware("client")(http.HandlerFunc(handlers.CreateFlat)))
	protectedRoutes.Handle("POST /flat/update", middleware.AuthMiddleware("client")(http.HandlerFunc(handlers.UpdateFlat)))

	http.ListenAndServe(cfg.HostAddr, nil)

	return nil
}
