package loadbalancer

import (
	"net/http"
	"reverse_proxy/config"
	"reverse_proxy/pkg/redis"
	"sync"
)

type LoadBalancer struct {
	Env            *config.SystemEnv
	Proxy          *config.ProxyMapping
	Client         *http.Client
	HealthyHostMap *sync.Map
	Redis          *redis.RedisClient
}
