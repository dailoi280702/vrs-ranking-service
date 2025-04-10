package http

import (
	"github.com/dailoi280702/vrs-ranking-service/handler/http/video"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPHandler() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "ok"})
	})

	g := e.Group("/api/v1")

	video.SetupRoute(g)

	return e
}
