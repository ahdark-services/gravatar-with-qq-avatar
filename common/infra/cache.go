package infra

import (
	"github.com/AH-dark/gravatar-with-qq-support/pkg/cache"
	"github.com/redis/go-redis/v9"
)

func NewCacheDriver(redisClient redis.UniversalClient) cache.Driver {
	return cache.NewRedisDriver(redisClient)
}
