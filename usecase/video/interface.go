package video

import (
	"context"

	"github.com/dailoi280702/vrs-ranking-service/type/model"
	"github.com/dailoi280702/vrs-ranking-service/type/request"
)

type I interface {
	UpdateInteraction(ctx context.Context, req request.UpdateInteraction) error
	GetTopVideos(ctx context.Context, req request.GetTopVideos) ([]model.Video, error)
}
