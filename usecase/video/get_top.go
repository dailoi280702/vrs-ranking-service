package video

import (
	"context"
	"slices"
	"sync"

	"github.com/dailoi280702/vrs-general-service/proto"
	"github.com/dailoi280702/vrs-ranking-service/type/model"
	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/constant"
	"github.com/dailoi280702/vrs-ranking-service/util/converter"
)

func (u *Usecase) GetTopVideos(ctx context.Context, req request.GetTopVideos) ([]model.Video, error) {
	if err := req.Validate(); err != nil {
		return nil, apperror.ErrBadRequest().WithMessage(err.Error())
	}

	errCh := make(chan error)
	watchedId := []int64{}
	cachedIds := []int64{}
	var wg sync.WaitGroup
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if req.UserId != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()
			ids, err := u.fetchWatchedVideoIds(cancelCtx, *req.UserId)
			watchedId = ids
			errCh <- err
		}()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		ids, err := u.fetchRankedVideoIds(cancelCtx)
		cachedIds = ids
		errCh <- err
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			cancel()

			return nil, apperror.ErrInternal().WithError(err).WithMessage(err.Error())
		}
	}

	// Filter by user watch history
	ids := []int64{}
	for i := range cachedIds {
		if !slices.Contains(watchedId, cachedIds[i]) {
			ids = append(ids, cachedIds[i])
		}
	}

	resp, err := u.GeneralSerivceClient.GetVideosByIds(ctx, &proto.GetVideosByIdsRequest{
		Ids: ids,
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

	// Preserve cache order
	res := []model.Video{}
	for _, id := range ids {
		if v, ok := mVideo[id]; ok {
			res = append(res, v)
		}
	}

	return res, nil
}

func (u *Usecase) fetchRankedVideoIds(ctx context.Context) ([]int64, error) {
	strIds, err := u.Rdb.ZRevRange(ctx, constant.RedisVideoRankKey, 0, constant.RedisTopVideoMaxItems)
	if err != nil {
		return nil, err
	}

	return converter.StringToInt64SliceIgnoreError(strIds), nil
}

func (u *Usecase) fetchWatchedVideoIds(ctx context.Context, userId int64) ([]int64, error) {
	resp, err := u.GeneralSerivceClient.GetUserWatchedHistory(ctx, &proto.IdRequest{Id: userId})
	if err != nil {
		return nil, err
	}

	res := make([]int64, len(resp.Videos))
	for i := range res {
		res[i] = resp.Videos[i].GetId()
	}

	return res, nil
}
