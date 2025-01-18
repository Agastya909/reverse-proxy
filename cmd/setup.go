package main

import (
	"fmt"
	"net/http"
	"reverse_proxy/config"
	loadbalancer "reverse_proxy/pkg/load-balancer"
	"reverse_proxy/pkg/redis"
	trafficproxy "reverse_proxy/pkg/traffic-proxy"
	"sync"
)

func Setup() {
	var (
		httpClient    = &http.Client{}
		healtyHostMap = &sync.Map{}
	)
	env, proxyConfig := getConfigs()
	redis := setupRedis(&env)
	lb := loadbalancer.NewLoadBalancer(&env, &proxyConfig, httpClient, healtyHostMap, redis)
	proxy := trafficproxy.NewTrafficProxy(&env, &proxyConfig, healtyHostMap)
	lb.Start()
	proxy.Start(redis)
}

func getConfigs() (config.SystemEnv, config.ProxyMapping) {
	env := config.LoadSystemConfig()
	proxyConfig := config.LoadProxyConfig()
	return env, proxyConfig
}

func setupRedis(env *config.SystemEnv) *redis.RedisClient {
	redis := redis.NewRedisClient(env.Redis.DB, fmt.Sprintf("%s:%d", env.Redis.Host, env.Redis.Port), env.Redis.Password)
	redis.FlushDB()
	return redis
}
