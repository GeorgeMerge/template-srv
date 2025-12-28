package app

import (
	"context"
	"log/slog"

	"eventsv1/internal/config"
	"eventsv1/pkg/httpserver"

	"github.com/labstack/echo/v4"
)

type App struct {
	log        *slog.Logger
	cfg        config.App
	httpserver *httpserver.Server
}

func New(log *slog.Logger, appCfg config.App) *App {
	return &App{
		log: log,
		cfg: appCfg,
	}
}

func (a *App) Init(ctx context.Context) error {
	router := echo.New()
	a.httpserver = httpserver.New(ctx, router, a.cfg.HTTP.Port)

	return nil
}

func (a *App) GracefulShutdown(ctx context.Context) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, a.cfg.HTTP.ShutdownTimeout)
	defer shutdownCancel()

	if err := a.httpserver.Shutdown(shutdownCtx); err != nil {
		a.log.Error("failed to shutdown httpserver, closing forcefully", "error", err)
		return a.httpserver.Close()
	}

	return nil
}
