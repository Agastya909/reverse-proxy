package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func (lb *LoadBalancer) Start() {
	go lb.HealthCheck()
}

func (lb *LoadBalancer) HealthCheck() {
	var (
		ticker = time.NewTicker(time.Duration(lb.Env.HealthCheckPeriod) * time.Second)
	)
	defer ticker.Stop()

	for range ticker.C {
		var wg sync.WaitGroup
		for _, server := range lb.Proxy.ProxyServers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				url := fmt.Sprintf("%s://%s:%v%s", server.Host.Protocol, server.Host.Address, server.Host.Port, server.Host.Health)
				isHealthy := lb.isHealthy(url)
				if isHealthy {
					lb.HealthyHostMap.Store(server.Host.Name, server.Host)
				} else {
					lb.HealthyHostMap.Delete(server.Host.Name)
				}
			}()
		}
		lb.HealthyHostMap.Range(func(key, value interface{}) bool {
			log.Println("Healthy Hosts after check: ", key)
			return true
		})
		wg.Wait()
	}
}

func (lb *LoadBalancer) isHealthy(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}
	res, err := lb.Client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}
