package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"reverse_proxy/config"
)

func (lb *LoadBalancer) GetHostForRequest(routeBasedMatching bool, allowedHosts []string, key string) (config.Host, error) {
	var (
		host         config.Host
		healthyHosts []string
		redisKey     string = lb.Proxy.Algorithm
	)
	healthyHosts = lb.getAvailableHosts()

	if len(healthyHosts) == 0 {
		return host, fmt.Errorf("no available hosts")
	}

	if routeBasedMatching {
		redisKey = key
	}
	// if len(hosts) == 1 {
	// 	return , nil
	// }

	switch lb.Proxy.Algorithm {
	case "round_robin":
		host, err := lb.RoundRobin(healthyHosts, allowedHosts, routeBasedMatching, redisKey)
		if err != nil {
			return host, err
		}
		return host, nil
	default:
		host, err := lb.RoundRobin(healthyHosts, allowedHosts, routeBasedMatching, redisKey)
		if err != nil {
			return host, err
		}
		return host, nil
	}
}

func (lb *LoadBalancer) getAvailableHosts() []string {
	var hosts []string
	lb.HealthyHostMap.Range(func(key, value interface{}) bool {
		hosts = append(hosts, key.(string))
		return true
	})
	return hosts
}

func (lb *LoadBalancer) RoundRobin(healthyHosts, allowedHosts []string, routeBasesdMatching bool, redisKey string) (config.Host, error) {
	var (
		host config.Host
		ctx  = context.Background()
	)

	exists := lb.Redis.IsKeyExists(ctx, redisKey)
	if !exists {
		if routeBasesdMatching {
			host, err := lb.saveNewHostsToRedis(ctx, allowedHosts, redisKey)
			if err != nil {
				return config.Host{}, err
			}
			return host, nil
		}

		host, err := lb.saveNewHostsToRedis(ctx, healthyHosts, redisKey)
		if err != nil {
			return config.Host{}, err
		}
		return host, nil
	}

	if routeBasesdMatching {
		host, err := lb.getHostAndUpdateRedisRoundRobin(ctx, allowedHosts, redisKey)
		if err != nil {
			return config.Host{}, err
		}
		return host, nil
	}

	host, err := lb.getHostAndUpdateRedisRoundRobin(ctx, healthyHosts, redisKey)
	if err != nil {
		return config.Host{}, err
	}

	return host, nil
}

func (lb *LoadBalancer) saveNewHostsToRedis(ctx context.Context, hosts []string, key string) (config.Host, error) {
	var (
		host     interface{}
		hostList []interface{}
	)

	for _, h := range hosts {
		hostData, ok := lb.HealthyHostMap.Load(h)
		if ok {
			hostList = append(hostList, h)
			host = hostData
		}
	}

	if len(hostList) == 0 {
		return config.Host{}, fmt.Errorf("no healthy hosts found")
	}

	err := lb.Redis.UpsertArrayToRedis(ctx, key, hostList)
	if err != nil {
		return config.Host{}, fmt.Errorf("failed to add healthy hosts to redis: %v", err)
	}

	return host.(config.Host), nil
}

func (lb *LoadBalancer) getHostAndUpdateRedisRoundRobin(ctx context.Context, hosts []string, key string) (config.Host, error) {
	var (
		roundRobinListLength int64 = 0
		counter              int64 = 0
		host                 config.Host
		healthyHostFound     bool = false
	)
	roundRobinListLength, err := lb.Redis.GetArrayKeyLen(ctx, key)
	if err != nil {
		log.Println("Failed to get round robin list length: ", err)
		return config.Host{}, err
	}

	for !healthyHostFound {
		if counter == roundRobinListLength {
			healthyHostFound = true
		}
		counter++

		value, err := lb.Redis.GetAndRemoveFirstArrayItem(ctx, key)
		if err != nil {
			log.Printf("Failed to get %s value from redis: %s ", key, err)
			continue
		}

		data, ok := lb.HealthyHostMap.Load(value)
		if !ok {
			log.Println("Failed to get host from healthy host map: ", err)
			continue
		}

		host = data.(config.Host)
		err = lb.Redis.AddItemToArrayTail(ctx, key, value)
		if err != nil {
			log.Println("Failed to add item to round robin list: ", err)
		}
		healthyHostFound = true
	}

	return host, nil
}
