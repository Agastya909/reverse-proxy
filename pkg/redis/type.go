package redis

import "github.com/redis/go-redis/v9"

type RedisClient struct {
	Opts   *redis.Options
	Client *redis.Client
}
