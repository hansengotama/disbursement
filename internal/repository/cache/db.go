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
	Get(param GetParam) Response
	Set(param SetParam) Response
}

type Response interface {
	Result() (string, error)
}

type CacheRedis struct {
	RedisConn *redis.Client
}

func NewCacheRedis() ICacheRepository {
	return CacheRedis{
		RedisConn: redishelper.GetConnection(),
	}
}

func (r CacheRedis) Get(param GetParam) Response {
	return r.RedisConn.Get(param.Context, param.Key)
}

func (r CacheRedis) Set(param SetParam) Response {
	return r.RedisConn.Set(param.Context, param.Key, param.Value, param.TTL)
}
