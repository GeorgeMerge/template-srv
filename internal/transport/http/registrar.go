package http

import "github.com/labstack/echo/v4"

type EchoRegistrar interface {
	Register(router *echo.Echo)
}
