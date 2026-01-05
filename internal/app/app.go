package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"template-srv/internal/config"
	"template-srv/internal/transport/http/router"
)

type App struct {
	log        *slog.Logger
	cfg        config.App
	httpserver *http.Server
}

func New(log *slog.Logger, appCfg config.App) (*App, error) {
	a := &App{
		log: log,
		cfg: appCfg,
	}

	err := a.init(context.Background())
	if err != nil {
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

	err := a.httpserver.Shutdown(shutdownCtx)
	if err != nil {
		a.log.ErrorContext(ctx, "failed to shutdown httpserver, closing forcefully", "error", err)

		return a.httpserver.Close()
	}

	return nil
}

func (a *App) init(ctx context.Context) error {
	router := router.New(a.log)

	a.httpserver = &http.Server{
		Addr:    fmt.Sprintf(":%v", a.cfg.HTTP.Port),
		Handler: router,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadTimeout: a.cfg.HTTP.ReadTimeout,
	}

	return nil
}
