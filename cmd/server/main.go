package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"template-srv/pkg/log_utils"

	"template-srv/internal/app"
	"template-srv/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: log_utils.MapSlogLevel(cfg.Logger.Level),
	}))

	app, err := app.New(log, cfg.App)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := app.GracefulShutdown(context.WithoutCancel(ctx)); err != nil {
			log.Error("failed to shutdown app", "error", err)
		}
	}()

	errCh := app.Run(ctx)

	log.Info("service started")
	defer log.Info("service stopped")

	select {
	case <-ctx.Done():
		log.Info("signal received, shutting down", "reason", ctx.Err())
	case err := <-errCh:
		log.Error("http server error", "error", err)
	}
}
