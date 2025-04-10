package redis

import (
	"context"
	"fmt"

	"github.com/dailoi280702/vrs-ranking-service/config"
	"github.com/dailoi280702/vrs-ranking-service/log"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

var redisClient *RedisClient

func init() {
	cfg := config.GetConfig()

	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Pass,
		DB:       cfg.Redis.DB,
	})

	ping, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Logger().Error("Failed to connect Redis", "error", err, "config", cfg.Redis)
	}

	redisClient = &RedisClient{rdb: conn}

	log.Logger().Info(fmt.Sprintf("Sent a ping to redis, got: %v", ping))
}

func GetClient() *RedisClient {
	return redisClient
}

func (r RedisClient) Zadd(ctx context.Context, key string, member any, score float64) error {
	return r.rdb.ZAdd(ctx, key, redis.Z{Member: member, Score: score}).Err()
}

func (r RedisClient) ZIncrBy(ctx context.Context, key, member string, incrment float64) error {
	return r.rdb.ZIncrBy(ctx, key, incrment, member).Err()
}

func (r RedisClient) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.rdb.ZRevRange(ctx, key, start, stop).Result()
}
