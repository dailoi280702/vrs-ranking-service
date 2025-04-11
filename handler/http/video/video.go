package video

import (
	"net/http"

	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/usecase"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/echoutil"
	"github.com/labstack/echo/v4"
)

func SetupRoute(gr *echo.Group) {
	g := gr.Group("/videos")
	g.POST("/interactions", updateInteraction)
	g.GET("/top", getTop)
}

func updateInteraction(c echo.Context) error {
	var (
		uc  = usecase.New()
		req request.UpdateInteraction
	)

	if err := c.Bind(&req); err != nil {
		return echoutil.ReponseErr(c, apperror.ErrBadRequest().WithMessage(err.Error()))
	}

	err := uc.Video.UpdateInteraction(c.Request().Context(), req)
	if err != nil {
		return echoutil.ReponseErr(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func getTop(c echo.Context) error {
	var (
		uc  = usecase.New()
		req request.GetTopVideos
	)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	resp, err := uc.Video.GetTopVideos(c.Request().Context(), req)
	if err != nil {
		return echoutil.ReponseErr(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}
