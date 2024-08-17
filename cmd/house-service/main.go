package main

import (
	"context"

	"github.com/shhesterka04/house-service/internal/app"
	"github.com/shhesterka04/house-service/pkg/logger"
)

const (
	levelConfig = "debug"
	development = true
)

func main() {
	ctx := context.Background()

	logger.Init(
		logger.Config{
			Level:       levelConfig,
			Development: development,
		})

	if err := app.Run(ctx); err != nil {
		logger.Fatalf(ctx, "app run error: %v", err)
	}
}
