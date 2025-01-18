package trafficproxy

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	redisClient "reverse_proxy/pkg/redis"

	"github.com/gorilla/mux"
)

func (t *TrafficProxy) Start(rc *redisClient.RedisClient) {
	mux := mux.NewRouter()
	mux.PathPrefix("/").HandlerFunc(t.handleRequest(rc))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", t.SystemEnv.ProxyHttp.Host, t.SystemEnv.ProxyHttp.Port), mux); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}

func (t *TrafficProxy) handleRequest(rc *redisClient.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		t.HealthyHostMap.Range(func(key, value interface{}) bool {
			log.Println("Healthy Hosts in proxy: ", key)
			return true
		})
		id := vars["id"]
		// lowestCount := math.MaxInt64
		// hostToUse := ""
		// for _, value := range t.SetupEnv.ProxyServers {
		// 	reqCountStr, err := rc.Get(value.Host.Name)
		// 	log.Println("Request count for host ", value.Host.Name, " is ", reqCountStr)
		// 	if err == redis.Nil {
		// 		reqCountStr = "0"
		// 	} else if err != nil {
		// 		log.Printf("Failed to get request count for host %s: %v", value.Host.Name, err)
		// 		continue
		// 	}
		// 	reqCount, err := strconv.Atoi(reqCountStr)
		// 	if err != nil {
		// 		log.Printf("Failed to convert request count for host %s: %v", value.Host.Name, err)
		// 		continue
		// 	}
		// 	if reqCount < lowestCount {
		// 		lowestCount = reqCount
		// 		hostToUse = value.Host.Name
		// 	}
		// }
		// rc.AddRequestToRedisCounter(hostToUse)
		// time.Sleep(time.Second * time.Duration(rand.Intn(10)))
		// log.Println("Request completed with id : and data : ", id, rand.Intn(100))
		// rc.RemoveRequestFromRedisCounter(hostToUse)
		w.Write([]byte(fmt.Sprintf("Request completed with id : %s and data : %d", id, rand.Intn(100))))
	}
}
