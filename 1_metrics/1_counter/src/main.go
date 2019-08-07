package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var timeUp = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "uptime_seconds_total",
		Help: "Time in seconds the service is up",
	},
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(timeUp)
}

func main() {
	// Increment variable each second inside a goroutine (parallel)
	go func() {
		for {
			time.Sleep(time.Second)
			timeUp.Inc()
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
