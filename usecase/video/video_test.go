package video_test

import (
	"context"
	"testing"

	mockgsc "github.com/dailoi280702/vrs-general-service/mock/service"
	mockredis "github.com/dailoi280702/vrs-ranking-service/mock/redis"
	"github.com/dailoi280702/vrs-ranking-service/usecase/video"
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
