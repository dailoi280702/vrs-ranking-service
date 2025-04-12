package video

import (
	"strconv"

	"github.com/dailoi280702/vrs-ranking-service/type/request"
	_ "github.com/dailoi280702/vrs-ranking-service/type/response"
	"github.com/dailoi280702/vrs-ranking-service/usecase"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/echoutil"
	"github.com/labstack/echo/v4"
)

// updateInteractions godoc
//
//	@Summary		Update video interaction
//	@Description	Update video interaction (like, comment, share, view, watch)
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.UpdateInteraction	true	"Request body"
//	@Success		200		{object}	response.Data[string]
//	@Failure		400		{object}	response.Error
//	@Failure		500		{object}	response.Error
//	@Param			id		path		string	true	"Video ID"	example(5)
//	@Router			/videos/{id}/interactions [post]
func updateInteraction(c echo.Context) error {
	var (
		uc  = usecase.New()
		req request.UpdateInteraction
	)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echoutil.ReponseErr(c, apperror.ErrBadRequest().WithMessage("Invalid video Id"))
	}

	if err := c.Bind(&req); err != nil {
		return echoutil.ReponseErr(c, apperror.ErrBadRequest().WithMessage(err.Error()))
	}

	req.VideoId = id

	if err := uc.Video.UpdateInteraction(c.Request().Context(), req); err != nil {
		return echoutil.ReponseErr(c, err)
	}

	return echoutil.ReponseData(c, "")
}
