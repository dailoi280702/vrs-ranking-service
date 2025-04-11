package video

import (
	"errors"
	"net/http"

	"github.com/dailoi280702/vrs-ranking-service/log"
	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func SetupRoute(gr *echo.Group) {
	g := gr.Group("/videos")
	g.POST("/interactions", updateInterfaction)
	g.GET("/top", getTop)
}

func updateInterfaction(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		logger = log.Logger()
		uc     = usecase.New()
		req    request.UpdateInteraction
	)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	err := uc.Video.UpdateInteraction(c.Request().Context(), req)
	if err != nil {
		var validationErrrs validator.ValidationErrors

		if errors.As(err, &validationErrrs) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
			})
		}

		logger.ErrorContext(ctx, "Failed to register user", "error", err, "request", req)

		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusInternalServerError, map[string]any{
		"message": "ok",
	})
}

func getTop(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		logger = log.Logger()
		uc     = usecase.New()
		req    request.GetTopVideos
	)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	resp, err := uc.Video.GetTopVideos(c.Request().Context(), req)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to register user", "error", err)

		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusInternalServerError, resp)
}
