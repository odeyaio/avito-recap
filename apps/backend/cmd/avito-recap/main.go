package main

import (
	"avito-recap/internal/app"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	const op = "main.run"

	config, err := app.LoadConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, config, logger)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := application.Run(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
