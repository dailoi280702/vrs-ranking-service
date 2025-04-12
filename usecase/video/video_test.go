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
