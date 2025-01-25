package trafficproxy

import (
	"fmt"
	"log"
	"net/http"
	"reverse_proxy/config"
	loadbalancer "reverse_proxy/pkg/load-balancer"
	redisClient "reverse_proxy/pkg/redis"
	"reverse_proxy/utils"

	"github.com/gorilla/mux"
)

func (t *TrafficProxy) Start(rc *redisClient.RedisClient, lb *loadbalancer.LoadBalancer) {
	mux := mux.NewRouter()
	mux.PathPrefix("/").HandlerFunc(t.handleRequest(lb))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", t.SystemEnv.ProxyHttp.Host, t.SystemEnv.ProxyHttp.Port), mux); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}

func (t *TrafficProxy) handleRequest(lb *loadbalancer.LoadBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		allowedHosts, routeBasedMatching, pattern := utils.MatchUrl(urlPath, t.SetupEnv.RouteMatching)
		if routeBasedMatching {
			urlPath = ""
		}
		host, err := lb.GetHostForRequest(routeBasedMatching, allowedHosts, pattern)
		if err != nil || host == (config.Host{}) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		w.Write([]byte("Proxying request to " + host.Name))
	}
}
