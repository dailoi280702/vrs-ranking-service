package video

import (
	"github.com/dailoi280702/vrs-ranking-service/client/redis"
)

type Usecase struct {
	Rdb redis.I
}

func New(rdb *redis.RedisClient) I {
	return &Usecase{
		Rdb: rdb,
	}
}
