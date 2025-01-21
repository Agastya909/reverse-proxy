package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"reverse_proxy/config"
)

func (lb *LoadBalancer) GetHostForRequest() (config.Host, error) {
	var (
		host  config.Host
		hosts []string
	)
	hosts = lb.getAvailableHosts()

	if len(hosts) == 0 {
		return host, fmt.Errorf("no available hosts")
	}

	// if len(hosts) == 1 {
	// 	return , nil
	// }

	switch lb.Proxy.Algorithm {
	case "round-robin":
		host, err := lb.RoundRobin(hosts)
		if err != nil {
			return host, err
		}
		return host, nil
	default:
		host, err := lb.RoundRobin(hosts)
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

func (lb *LoadBalancer) RoundRobin(hosts []string) (config.Host, error) {
	var (
		host                 config.Host
		ctx                        = context.Background()
		healthyHostFound           = false
		counter              int64 = 0
		roundRobinListLength int64 = 0
	)

	exists := lb.Redis.IsKeyExists(ctx, "round_robin")
	if !exists {
		var hostList []interface{}
		data, ok := lb.HealthyHostMap.Load(hosts[0])
		if !ok {
			return host, fmt.Errorf("host not found")
		}
		for _, h := range hosts[1:] {
			hostList = append(hostList, h)
		}
		err := lb.Redis.UpsertArrayToRedis(context.Background(), "round_robin", hostList)
		if err != nil {
			log.Println("Failed to add healthy hosts to redis: ", err)
		}
		err = lb.Redis.AddItemToArrayTail(ctx, "round_robin", hosts[0])
		if err != nil {
			log.Println("Failed to add item to round robin list: ", err)
		}
		host = data.(config.Host)
		return host, nil
	}

	roundRobinListLength, err := lb.Redis.GetArrayKeyLen(ctx, "round_robin")
	if err != nil {
		log.Println("Failed to get round robin list length: ", err)
		return host, err
	}

	for !healthyHostFound {
		if counter == roundRobinListLength {
			healthyHostFound = true
		}
		counter++
		value, err := lb.Redis.GetAndRemoveFirstArrayItem(ctx, "round_robin")
		if err != nil {
			log.Println("Failed to get round robin value from redis: ", err)
			continue
		}
		data, ok := lb.HealthyHostMap.Load(value)
		if !ok {
			log.Println("Failed to get host from healthy host map: ", err)
			continue
		}
		host = data.(config.Host)
		err = lb.Redis.AddItemToArrayTail(ctx, "round_robin", value)
		if err != nil {
			log.Println("Failed to add item to round robin list: ", err)
		}
		healthyHostFound = true
	}

	return host, nil
}
