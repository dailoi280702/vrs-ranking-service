package video

import (
	"github.com/dailoi280702/vrs-general-service/proto"
	"github.com/dailoi280702/vrs-ranking-service/client/redis"
)

type Usecase struct {
	Rdb                  redis.I
	GeneralSerivceClient proto.ServiceClient
}

func New(rdb *redis.RedisClient, gsc proto.ServiceClient) I {
	return &Usecase{
		Rdb:                  rdb,
		GeneralSerivceClient: gsc,
	}
}
