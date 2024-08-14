package main

import (
	"context"

	"github.com/shhesterka04/house-service/internal/app"
	"github.com/shhesterka04/house-service/pkg/logger"
)

func main() {
	ctx := context.Background()

	logger.Init()

	if err := app.Run(ctx); err != nil {
		logger.Fatalf(ctx, "app run error: %v", err)
	}
}
