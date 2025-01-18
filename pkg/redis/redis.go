package redis

import "context"

func (r *RedisClient) Set(key string, value interface{}) error {
	return r.Client.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(context.Background(), key).Result()
}

func (r *RedisClient) AddRequestToRedisCounter(hostName string) error {
	return r.Client.Incr(context.Background(), hostName).Err()
}

func (r *RedisClient) RemoveRequestFromRedisCounter(hostName string) (int64, error) {
	return r.Client.Decr(context.Background(), hostName).Result()
}

func (r *RedisClient) FlushDB() error {
	return r.Client.FlushDBAsync(context.Background()).Err()
}
