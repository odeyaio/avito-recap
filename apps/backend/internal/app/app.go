package app

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type App struct {
	server echo.StartConfig
	router *echo.Echo
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (a *App) Run(ctx context.Context) error {
	defer a.pool.Close()

	a.logger.InfoContext(ctx, "starting application", "address", a.server.Address)
	return a.server.Start(ctx, a.router)
}
