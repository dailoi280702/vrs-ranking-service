package video

import (
	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/type/response"
	"github.com/dailoi280702/vrs-ranking-service/usecase"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/echoutil"
	"github.com/labstack/echo/v4"
)

// getTop godoc
//
//	@Summary		Get top videos
//	@Description	Get top videos, optionally filtered by user watch history
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		integer	false	"User ID"	example(3)
//	@Success		200		{object}	response.Data[[]response.Video]
//	@Failure		400		{object}	response.Error
//	@Failure		500		{object}	response.Error
//	@Router			/videos/top [get]
func getTop(c echo.Context) error {
	var (
		uc  = usecase.New()
		req request.GetTopVideos
	)

	if err := c.Bind(&req); err != nil {
		return echoutil.ReponseErr(c, apperror.ErrBadRequest().WithMessage(err.Error()))
	}

	resp, err := uc.Video.GetTopVideos(c.Request().Context(), req)
	if err != nil {
		return echoutil.ReponseErr(c, err)
	}

	return echoutil.ReponseData(c, response.FormVideosResponse(resp...))
}
