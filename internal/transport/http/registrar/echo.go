package registrar

import "github.com/labstack/echo/v4"

type Echo interface {
	Register(router *echo.Echo)
}
