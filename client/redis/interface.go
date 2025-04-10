package redis

import "context"

type Interface interface {
	Zadd(ctx context.Context, key string, member any, score float64) error
	ZIncrBy(ctx context.Context, key, member string, incrment float64) error
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}
