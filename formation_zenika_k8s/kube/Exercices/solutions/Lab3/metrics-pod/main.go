package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_elapsed_seconds",
		Help: "The total number of elapsed seconds",
	})
)

func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":80", nil)
}
