package video

import (
	"context"

	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/type/response"
)

type I interface {
	UpdateInteraction(ctx context.Context, req request.UpdateInteraction) error
	GetTopVideos(ctx context.Context, req request.GetTopVideos) ([]response.Video, error)
}
