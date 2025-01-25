package redis

import (
	"context"
	"fmt"
)

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) AddRequestToRedisCounter(ctx context.Context, hostName string) error {
	return r.Client.Incr(ctx, hostName).Err()
}

func (r *RedisClient) RemoveRequestFromRedisCounter(ctx context.Context, hostName string) (int64, error) {
	return r.Client.Decr(ctx, hostName).Result()
}

func (r *RedisClient) FlushDB(ctx context.Context) error {
	return r.Client.FlushDBAsync(context.Background()).Err()
}

func (r *RedisClient) UpsertArrayToRedis(ctx context.Context, key string, values []interface{}) error {
	pipe := r.Client.Pipeline()
	redisCmd := r.Client.Exists(ctx, key)
	exists, err := redisCmd.Val(), redisCmd.Err()
	if err != nil {
		return fmt.Errorf("error checking if key exists: %v", err)
	}
	if exists == 1 {
		pipe.Del(ctx, key)
	}
	err = pipe.RPush(ctx, key, values...).Err()
	if err != nil {
		return err
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetAndRemoveFirstArrayItem(ctx context.Context, key string) (string, error) {
	return r.Client.LPop(ctx, key).Result()
}

func (r *RedisClient) AddItemToArrayTail(ctx context.Context, key string, value interface{}) error {
	return r.Client.RPush(ctx, key, value).Err()
}

func (r *RedisClient) IsKeyExists(ctx context.Context, key string) bool {
	val := r.Client.Exists(ctx, key)
	if val.Err() != nil || val.Val() == 0 {
		return false
	}
	return true
}

func (r *RedisClient) GetArrayKeyLen(ctx context.Context, key string) (int64, error) {
	return r.Client.LLen(ctx, key).Result()
}
