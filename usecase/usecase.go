package usecase

import (
	"github.com/dailoi280702/vrs-ranking-service/client/generalservice"
	"github.com/dailoi280702/vrs-ranking-service/client/redis"
	videouc "github.com/dailoi280702/vrs-ranking-service/usecase/video"
)

type Usecase struct {
	Video videouc.I
}

func New() *Usecase {
	rdb := redis.GetClient()
	gsc := generalservice.GetClient()

	return &Usecase{
		Video: videouc.New(rdb, gsc),
	}
}
