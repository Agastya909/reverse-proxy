package trafficproxy

import (
	"reverse_proxy/config"
	"sync"
)

type TrafficProxy struct {
	SystemEnv      config.SystemEnv
	SetupEnv       config.ProxyMapping
	HealthyHostMap *sync.Map
}
