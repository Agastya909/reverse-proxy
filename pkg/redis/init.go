package redis

import "github.com/redis/go-redis/v9"

func NewRedisClient(db int, addr string, password *string) *RedisClient {
	opts := redis.Options{
		Addr: addr,
		DB:   db,
	}
	if password != nil {
		opts.Password = *password
	}
	redisClient := redis.NewClient(&opts)
	return &RedisClient{
		Client: redisClient,
		Opts:   &opts,
	}
}
