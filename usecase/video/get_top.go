package video

import (
	"context"

	"github.com/dailoi280702/vrs-general-service/proto"
	"github.com/dailoi280702/vrs-ranking-service/type/model"
	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/constant"
	"github.com/dailoi280702/vrs-ranking-service/util/converter"
)

func (u *Usecase) GetTopVideos(ctx context.Context, req request.GetTopVideos) ([]model.Video, error) {
	strIds, err := u.Rdb.ZRevRange(ctx, constant.RedisVideoRankKey, 0, constant.RedisTopVideoMaxItems)
	if err != nil {
		return nil, apperror.ErrInternal().WithError(err).WithMessage(err.Error())
	}

	videoIds := converter.StringToInt64SliceIgnoreError(strIds)

	resp, err := u.GeneralSerivceClient.GetVideosByIds(ctx, &proto.GetVideosByIdsRequest{
		Ids: videoIds,
	})
	if err != nil {
		return nil, apperror.ErrGRPC(err)
	}

	mVideo := map[int64]model.Video{}

	for _, v := range resp.Videos {
		mVideo[v.GetId()] = model.Video{
			ID:        v.GetId(),
			Comments:  v.GetComments(),
			Length:    v.GetLength(),
			Likes:     v.GetLikes(),
			Name:      v.GetName(),
			Shares:    v.GetShares(),
			Views:     v.GetViews(),
			WatchTime: v.GetWatchTime(),
		}
	}

	// Preserve Redis order
	res := []model.Video{}
	for _, id := range videoIds {
		if v, ok := mVideo[id]; ok {
			res = append(res, v)
		}
	}

	return res, nil
}
