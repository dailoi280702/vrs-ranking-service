package echoutil

import (
	"errors"
	"net/http"

	"github.com/dailoi280702/vrs-ranking-service/log"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/labstack/echo/v4"
)

func ReponseErr(c echo.Context, err error) error {
	var (
		logger   = log.Logger()
		ctx      = c.Request().Context()
		appError *apperror.AppError
	)

	if err != nil {
		if errors.As(err, &appError) {
			if appError.Code == http.StatusInternalServerError {
				logger.ErrorContext(ctx, "Internal server error", "error", err)
			}

			return c.JSON(appError.Code, map[string]any{
				"message": appError.Message,
			})
		}
	}

	logger.ErrorContext(ctx, "Internal server error", "error", err)

	return c.JSON(http.StatusInternalServerError, map[string]any{
		"message": err.Error(),
	})
}
