package trafficproxy

import (
	"reverse_proxy/config"
	"sync"
)

func NewTrafficProxy(cfg *config.SystemEnv, proxyCfg *config.ProxyMapping, healtyHosts *sync.Map) *TrafficProxy {
	return &TrafficProxy{
		SystemEnv:      *cfg,
		SetupEnv:       *proxyCfg,
		HealthyHostMap: healtyHosts,
	}
}
