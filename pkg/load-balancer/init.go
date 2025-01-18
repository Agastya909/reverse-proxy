package loadbalancer

import (
	"net/http"
	"reverse_proxy/config"
	"reverse_proxy/pkg/redis"
	"sync"
)

func NewLoadBalancer(
	env *config.SystemEnv,
	proxyMap *config.ProxyMapping,
	client *http.Client,
	hostMap *sync.Map,
	redis *redis.RedisClient,
) *LoadBalancer {
	return &LoadBalancer{
		Env:            env,
		Proxy:          proxyMap,
		Client:         client,
		HealthyHostMap: hostMap,
		Redis:          redis,
	}
}
