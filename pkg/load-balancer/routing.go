package loadbalancer

import (
	"fmt"
	"reverse_proxy/config"
)

func (lb *LoadBalancer) GetHostForRequest() (string, error) {
	hosts := lb.getAvailableHosts()
	if len(hosts) == 0 {
		return "", fmt.Errorf("no available hosts")
	}

	switch lb.Proxy.Algorithm {
	case "round-robin":
		lb.RoundRobin(hosts)
	default:
		lb.RoundRobin(hosts)
	}

	return hosts[0], nil
}

func (lb *LoadBalancer) getAvailableHosts() []string {
	var hosts []string
	lb.HealthyHostMap.Range(func(key, value interface{}) bool {
		fmt.Println("Key: value: ", key, value)
		hosts = append(hosts, key.(string))
		return true
	})
	return hosts
}

func (lb *LoadBalancer) RoundRobin(hosts []string) (config.Host, error) {
	var host config.Host
	return host, nil
}
