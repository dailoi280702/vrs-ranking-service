package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPHandler() *echo.Echo {
	e := echo.New()
	// cfg = config.GetConfig()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "ok"})
	})

	return e
}
