package app

import (
	"context"
	"log/slog"

	"template-srv/internal/config"
	"template-srv/internal/transport/http"
	"template-srv/internal/transport/http/ping"
	"template-srv/pkg/httpserver"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	log        *slog.Logger
	cfg        config.App
	httpserver *httpserver.Server
}

func New(log *slog.Logger, appCfg config.App) (*App, error) {
	a := &App{
		log: log,
		cfg: appCfg,
	}

	if err := a.init(context.Background()); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) RunAsync() chan error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		errCh <- a.httpserver.ListenAndServe()
	}()

	return errCh
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

func (a *App) init(ctx context.Context) error {
	router := echo.New()

	a.registerMiddlewares(router)
	a.registerHandlers(router)
	a.httpserver = httpserver.New(ctx, router, a.cfg.HTTP.Port)

	return nil
}

func (a *App) registerHandlers(router *echo.Echo) {
	registrars := []http.EchoRegistrar{
		ping.NewHandler(),
	}

	for _, registrar := range registrars {
		registrar.Register(router)
	}
}

func (a *App) registerMiddlewares(router *echo.Echo) {
	router.Use(middleware.Recover())

	skipper := func(c echo.Context) bool {
		paths := map[string]any{
			// list uris to skip logging
		}

		_, ok := paths[c.Request().URL.Path]
		return ok
	}

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		Skipper:   skipper,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			a.log.Info("request", "uri", v.URI, "status", v.Status)
			return nil
		},
	}))
}
