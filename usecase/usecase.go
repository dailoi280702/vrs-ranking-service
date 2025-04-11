package usecase

import (
	"github.com/dailoi280702/vrs-ranking-service/client/redis"
	"github.com/dailoi280702/vrs-ranking-service/usecase/video"
)

type Usecase struct {
	Video video.I
}

func New() *Usecase {
	rdb := redis.GetClient()

	return &Usecase{
		Video: video.New(rdb),
	}
}
