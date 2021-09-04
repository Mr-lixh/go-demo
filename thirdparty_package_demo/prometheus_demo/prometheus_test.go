package prometheus_demo

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"testing"
	"time"
)

var (
	requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_processed_total",
		Help: "The total number of processed events",
	})
)

func TestPromCounter(t *testing.T) {
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	requestsProcessed.Inc()

	done := make(chan bool)

	log.Println("Server is ready to handle requests at: 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not listen on 8080: %v\n", err)
	}

	<-done
}
