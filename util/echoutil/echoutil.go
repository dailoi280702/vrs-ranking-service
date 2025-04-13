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
			if appError.Code >= 500 {
				logger.ErrorContext(ctx, "Internal server error", "error", err)
			}

			return c.JSON(appError.Code, appError)
		}
	}

	logger.ErrorContext(ctx, "Internal server error", "error", err)

	return c.JSON(http.StatusInternalServerError, apperror.ErrInternal().WithMessage(err.Error()))
}

type Data[T any] struct {
	Code    int32  `json:"code" example:"200" description:"Response code"`
	Message string `json:"message" example:"success" description:"Response message"`
	Data    T      `json:"data" description:"Response data"`
}

func ReponseData[T any](c echo.Context, data T) error {
	return c.JSON(http.StatusOK, Data[T]{
		Code:    http.StatusOK,
		Data:    data,
		Message: "success",
	})
}
