package video

import (
	"context"

	"github.com/dailoi280702/vrs-general-service/proto"
	"github.com/dailoi280702/vrs-ranking-service/log"
	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/constant"
)

func (u *Usecase) UpdateInteraction(ctx context.Context, req request.UpdateInteraction) error {
	if err := req.Validate(); err != nil {
		return apperror.ErrBadRequest().WithMessage(err.Error())
	}

	video, err := u.GeneralSerivceClient.GetVideoByID(ctx, &proto.IdRequest{Id: req.VideoId})
	if err != nil {
		return apperror.ErrGRPC(err)
	}

	switch req.Type {
	case request.VideoInteractionComment:
		video.Comments++
	case request.VideoInteractionLike:
		video.Likes++
	case request.VideoInteractionShare:
		video.Shares++
	case request.VideoInteractionView:
		video.Views++
	case request.VideoInteractionWatch:
		video.WatchTime += req.WatchTimeSeconds
	}

	_, err = u.GeneralSerivceClient.UpdateVideo(ctx, video)
	if err != nil {
		return apperror.ErrGRPC(err)
	}

	var score float64
	if video.GetViews() != 0 {
		score += float64(video.GetComments()+video.GetLikes()+video.GetShares()) / float64(video.GetViews())
	}

	if video.GetLength() != 0 && video.GetViews() != 0 {
		score += float64(video.GetWatchTime()) / float64(video.GetLength()*video.GetViews())
	}

	go func(videoId int64, score float64) {
		logger := log.Logger()
		if err := u.Rdb.Zadd(ctx, constant.RedisVideoRankKey, video.GetId(), score); err != nil {
			logger.Error("Redis error", "error", err, "key", constant.RedisVideoRankKey, "member", videoId)
		}
	}(req.VideoId, score)

	return nil
}
