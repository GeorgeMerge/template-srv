package router

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"template-srv/internal/transport/http/handlers/ping"
	"template-srv/internal/transport/http/registrar"
)

type Router struct {
	*echo.Echo

	log *slog.Logger
}

func New(log *slog.Logger) *Router {
	r := &Router{
		Echo: echo.New(),
		log:  log,
	}
	r.init()

	return r
}

func (r *Router) init() {
	r.registerHandlers()
	r.registerMiddlewares()
}

func (r *Router) registerHandlers() {
	registrars := []registrar.Echo{
		ping.NewHandler(),
	}

	for _, registrar := range registrars {
		registrar.Register(r.Echo)
	}
}

func (r *Router) registerMiddlewares() {
	r.Use(middleware.Recover())

	skipper := func(c echo.Context) bool {
		paths := map[string]any{
			// list uris to skip logging
		}

		_, ok := paths[c.Request().URL.Path]

		return ok
	}

	r.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		Skipper:   skipper,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			r.log.Info("request", "uri", v.URI, "status", v.Status)

			return nil
		},
	}))
}
