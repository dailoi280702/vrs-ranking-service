package video

import (
	"github.com/labstack/echo/v4"
)

func SetupRoute(gr *echo.Group) {
	g := gr.Group("/videos")
	g.GET("/top", getTop)
	g.POST("/:id/interactions", updateInteraction)
}
