package http

import (
	"github.com/dailoi280702/vrs-ranking-service/handler/http/video"
	"github.com/dailoi280702/vrs-ranking-service/util/echoutil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewHTTPHandler() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return echoutil.ReponseData(c, "ok")
	})

	e.GET("/healthz", func(c echo.Context) error {
		return echoutil.ReponseData(c, "ok")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	g := e.Group("/api/v1")

	video.SetupRoute(g)

	return e
}
