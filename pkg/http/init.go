package http

import "github.com/redis/go-redis/v9"

type HttpServer struct {
	Redis redis.Client
}
