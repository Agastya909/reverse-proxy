package trafficproxy

import (
	"fmt"
	"log"
	"net/http"
	"reverse_proxy/config"
	loadbalancer "reverse_proxy/pkg/load-balancer"
	redisClient "reverse_proxy/pkg/redis"

	"github.com/gorilla/mux"
)

func (t *TrafficProxy) Start(rc *redisClient.RedisClient, lb *loadbalancer.LoadBalancer) {
	mux := mux.NewRouter()
	mux.PathPrefix("/").HandlerFunc(t.handleRequest(rc, lb))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", t.SystemEnv.ProxyHttp.Host, t.SystemEnv.ProxyHttp.Port), mux); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}

func (t *TrafficProxy) handleRequest(rc *redisClient.RedisClient, lb *loadbalancer.LoadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, err := lb.GetHostForRequest()
		if err != nil || host == (config.Host{}) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		w.Write([]byte("Proxying request to " + host.Name))
	}
}
