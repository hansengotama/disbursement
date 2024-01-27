package cacherepo

import (
	"context"
	redishelper "github.com/hansengotama/disbursement/internal/lib/redis"
	"github.com/redis/go-redis/v9"
	"time"
)

type GetParam struct {
	Context context.Context
	Key     string
}

type SetParam struct {
	Context context.Context
	Key     string
	Value   any
	TTL     time.Duration
}

type ICacheRepository interface {
	Get(param GetParam) (string, error)
	Set(param SetParam) (string, error)
}

type CacheRedis struct {
	RedisConn *redis.Client
}

func NewCacheRedis() ICacheRepository {
	return CacheRedis{
		RedisConn: redishelper.GetConnection(),
	}
}

func (r CacheRedis) Get(param GetParam) (string, error) {
	return r.RedisConn.Get(param.Context, param.Key).Result()
}

func (r CacheRedis) Set(param SetParam) (string, error) {
	return r.RedisConn.Set(param.Context, param.Key, param.Value, param.TTL).Result()
}
