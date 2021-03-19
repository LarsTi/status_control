package main
import(
	"log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"fmt"
	"os"
	"encoding/json"
)
type service struct{
	Status          string `json:"status"`
	Id              string `json:"status_id"`
	Secret          string `json:"secret,omitempty"`
	StatusIP        string `json:"status_ip_changed"`
	Default         string `json:"default,omitempty"`
	Answer          string `json:"answer,omitempty"`
}
var (
	services []service
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "status_http_duration_seconds",
		Help: "Duration of HTTP requests",
	}, []string{"route"})
)

func readEnv() []service{
	i := 1
	var ret []service
	for{
		s := service{
			Status: os.Getenv(fmt.Sprintf("DEFAULT_%d", i)),
			Secret: os.Getenv(fmt.Sprintf("SECRET_%d", i)),
			Id: os.Getenv(fmt.Sprintf("STATUS_%d", i)),
			Default: os.Getenv(fmt.Sprintf("DEFAULT_%d", i)),
		}
		if s.Secret != "" && s.Id != "" {
			log.Printf("Loaded service %s", s.Id)
			ret = append(ret, s)
			i ++
		}else{
			break
		}
	}
	log.Printf("Loaded %d services", len(ret))
	return ret
}
func main() {
	services = readEnv()

	router := mux.NewRouter()
	router.HandleFunc("/get", get).Methods("GET")
	router.HandleFunc("/update", update).Methods("GET", "POST")
	router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	router.Use(loggingMiddleware)
	router.Use(prometheusMiddleware)
	log.Println("works")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080) , router))
}


func update( w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	s := service{}
	s.Answer = "Nicht erlaubt!"
	if len(params["id"]) != 1 ||
	len(params["pw"]) != 1 ||
	len(params["status"]) != 1{

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(s)
		return
	}
	for index, entry := range services {
		if entry.Id == params["id"][0] &&
		entry.Secret == params["pw"][0]{
			entry.Status = params["status"][0]
			if entry.Status == "" {
				entry.Status = entry.Default
			}
			entry.StatusIP = r.Header.Get("X-Forwarded-For")
			services[index] = entry

			s = entry
			s.Default = ""
			s.Secret = ""
			s.Answer = "Ok"
			break
		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)

}
func get( w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	s := service{}
	if len(params["id"]) == 1 {
		for _, s = range services {
			if s.Id == params["id"][0] {
				break
			}
		}
	}
	if s.Status == "" {
		s.Status = s.Default
	}
	s.Default = ""
	s.Secret = ""
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[webapi-request] %s: %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

