package video_test

import (
	"context"
	"testing"

	mockgsc "github.com/dailoi280702/vrs-general-service/mock/service"
	"github.com/dailoi280702/vrs-general-service/proto"
	mockredis "github.com/dailoi280702/vrs-ranking-service/mock/redis"
	"github.com/dailoi280702/vrs-ranking-service/type/request"
	"github.com/dailoi280702/vrs-ranking-service/usecase/video"
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite

	ctx       context.Context
	uc        *video.Usecase
	mockRedis *mockredis.MockI
	mockGSC   *mockgsc.MockServiceClient
}

func (s *TestSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockRedis = &mockredis.MockI{}
	s.mockGSC = &mockgsc.MockServiceClient{}

	s.uc = &video.Usecase{
		Rdb:                  s.mockRedis,
		GeneralSerivceClient: s.mockGSC,
	}
}

func TestVideoUsecase(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) TestUpdateInteraction_Success() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionLike,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), float64(0)).Return(nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
}

func (s *TestSuite) TestGetTopVideos_Success() {
	req := request.GetTopVideos{
		UserId: nil,
	}

	video1 := &proto.Video{
		Id:        1,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
		Name:      "Video 1",
	}

	video2 := &proto.Video{
		Id:        2,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
		Name:      "Video 2",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mockRedis.EXPECT().ZRevRange(ctx, "video_rank", int64(0), int64(20)).Return([]string{"1", "2"}, nil)
	s.mockGSC.EXPECT().GetVideosByIds(s.ctx, &proto.GetVideosByIdsRequest{Ids: []int64{1, 2}}).Return(&proto.Videos{Videos: []*proto.Video{video1, video2}}, nil)

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.NoError(err)
	s.Len(videos, 2)
	s.Equal(int64(1), videos[0].ID)
	s.Equal("Video 1", videos[0].Name)
	s.Equal(int64(2), videos[1].ID)
	s.Equal("Video 2", videos[1].Name)
}

func (s *TestSuite) TestGetTopVideos_Success_With_User_ID() {
	userID := int64(123)
	req := request.GetTopVideos{
		UserId: &userID,
	}

	video1 := &proto.Video{
		Id:        1,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
		Name:      "Video 1",
	}

	video2 := &proto.Video{
		Id:        2,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
		Name:      "Video 2",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mockRedis.EXPECT().ZRevRange(ctx, "video_rank", int64(0), int64(20)).Return([]string{"1", "2"}, nil)
	s.mockGSC.EXPECT().GetUserWatchedHistory(ctx, &proto.IdRequest{Id: userID}).Return(&proto.Videos{Videos: []*proto.Video{}}, nil)
	s.mockGSC.EXPECT().GetVideosByIds(s.ctx, &proto.GetVideosByIdsRequest{Ids: []int64{1, 2}}).Return(&proto.Videos{Videos: []*proto.Video{video1, video2}}, nil)

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.NoError(err)
	s.Len(videos, 2)
	s.Equal(int64(1), videos[0].ID)
	s.Equal("Video 1", videos[0].Name)
	s.Equal(int64(2), videos[1].ID)
	s.Equal("Video 2", videos[1].Name)
}

func (s *TestSuite) TestGetTopVideos_Success_With_User_ID_And_Watched_Videos() {
	userID := int64(123)
	req := request.GetTopVideos{
		UserId: &userID,
	}

	video1 := &proto.Video{
		Id:        1,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
		Name:      "Video 1",
	}

	video2 := &proto.Video{
		Id:        2,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
		Name:      "Video 2",
	}

	watchedVideo := &proto.Video{
		Id:        3,
		Likes:     5,
		Comments:  1,
		Shares:    1,
		Views:     1,
		Length:    100,
		WatchTime: 25,
		Name:      "Watched Video",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mockRedis.EXPECT().ZRevRange(ctx, "video_rank", int64(0), int64(20)).Return([]string{"1", "2", "3"}, nil)
	s.mockGSC.EXPECT().GetUserWatchedHistory(ctx, &proto.IdRequest{Id: userID}).Return(&proto.Videos{Videos: []*proto.Video{watchedVideo}}, nil)
	s.mockGSC.EXPECT().GetVideosByIds(s.ctx, &proto.GetVideosByIdsRequest{Ids: []int64{1, 2}}).Return(&proto.Videos{Videos: []*proto.Video{video1, video2}}, nil)

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.NoError(err)
	s.Len(videos, 2)
	s.Equal(int64(1), videos[0].ID)
	s.Equal("Video 1", videos[0].Name)
	s.Equal(int64(2), videos[1].ID)
	s.Equal("Video 2", videos[1].Name)
}

func (s *TestSuite) TestGetTopVideos_InvalidRequest() {
	req := request.GetTopVideos{
		UserId: new(int64), // Invalid user ID (0)
	}

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.Error(err)
	s.Nil(videos)
	s.Equal(400, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestGetTopVideos_FetchWatchedVideoIds_Error() {
	userID := int64(123)
	req := request.GetTopVideos{
		UserId: &userID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mockRedis.EXPECT().ZRevRange(ctx, "video_rank", int64(0), int64(20)).Return([]string{"1", "2"}, nil)
	s.mockGSC.EXPECT().GetUserWatchedHistory(ctx, &proto.IdRequest{Id: userID}).Return(nil, apperror.ErrInternal())

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.Error(err)
	s.Nil(videos)
	s.Equal(500, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestGetTopVideos_FetchRankedVideoIds_Error() {
	req := request.GetTopVideos{
		UserId: nil,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mockRedis.EXPECT().ZRevRange(ctx, "video_rank", int64(0), int64(20)).Return(nil, apperror.ErrInternal())

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.Error(err)
	s.Nil(videos)
	s.Equal(500, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestGetTopVideos_GetVideosByIds_Error() {
	req := request.GetTopVideos{
		UserId: nil,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mockRedis.EXPECT().ZRevRange(ctx, "video_rank", int64(0), int64(20)).Return([]string{"1", "2"}, nil)
	s.mockGSC.EXPECT().GetVideosByIds(s.ctx, &proto.GetVideosByIdsRequest{Ids: []int64{1, 2}}).Return(nil, apperror.ErrInternal())

	videos, err := s.uc.GetTopVideos(s.ctx, req)
	s.Error(err)
	s.Nil(videos)
	s.Equal(500, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestUpdateInteraction_WatchTimeExceedsLength() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionWatch,
		WatchTimeSeconds: 101,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.Error(err)
	s.Equal(400, err.(*apperror.AppError).Code)
	s.Equal("Watch time can not exceed video length", err.(*apperror.AppError).Message)
}

func (s *TestSuite) TestUpdateInteraction_InvalidRequest() {
	req := request.UpdateInteraction{}

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.Error(err)
	s.Equal(400, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestUpdateInteraction_GetVideoByID_Error() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionLike,
		WatchTimeSeconds: 0,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(nil, apperror.ErrInternal())

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.Error(err)
	s.Equal(500, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestUpdateInteraction_UpdateVideo_Error() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionLike,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, apperror.ErrInternal())

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.Error(err)
	s.Equal(500, err.(*apperror.AppError).Code)
}

func (s *TestSuite) TestUpdateInteraction_Redis_Error() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionLike,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), mock.Anything).Return(apperror.ErrInternal())

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
}

func (s *TestSuite) TestUpdateInteraction_Success_Comment() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionComment,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), mock.Anything).Return(nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
	s.Equal(int64(1), video.GetComments())
}

func (s *TestSuite) TestUpdateInteraction_Success_View() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionView,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), mock.Anything).Return(nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
	s.Equal(int64(1), video.GetViews())
}

func (s *TestSuite) TestUpdateInteraction_Success_Watch() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionWatch,
		WatchTimeSeconds: 10,
	}

	video := &proto.Video{
		Id:        req.VideoId,
		Likes:     0,
		Comments:  0,
		Shares:    0,
		Views:     0,
		Length:    100,
		WatchTime: 0,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), mock.Anything).Return(nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
	s.Equal(int64(10), video.GetWatchTime())
}

func (s *TestSuite) TestUpdateInteraction_Success_Share() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionShare,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:       req.VideoId,
		Likes:    0,
		Comments: 0,
		Shares:   0,
		Views:    0,
		Length:   100,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), mock.Anything).Return(nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
	s.Equal(int64(1), video.GetShares())
}

func (s *TestSuite) TestUpdateInteraction_Success_Score_Calculation() {
	req := request.UpdateInteraction{
		VideoId:          1,
		Type:             request.VideoInteractionLike,
		WatchTimeSeconds: 0,
	}

	video := &proto.Video{
		Id:        req.VideoId,
		Likes:     10,
		Comments:  5,
		Shares:    2,
		Views:     2,
		Length:    100,
		WatchTime: 50,
	}

	s.mockGSC.EXPECT().GetVideoByID(s.ctx, &proto.IdRequest{Id: req.VideoId}).Return(video, nil)
	s.mockGSC.EXPECT().UpdateVideo(s.ctx, video).Return(nil, nil)
	s.mockRedis.EXPECT().Zadd(s.ctx, "video_rank", video.GetId(), mock.Anything).Return(nil)

	err := s.uc.UpdateInteraction(s.ctx, req)
	s.NoError(err)
}
