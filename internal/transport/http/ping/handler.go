package ping

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Register(router *echo.Echo) {
	router.POST("/ping", h.ping)
}

func (h *Handler) ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
