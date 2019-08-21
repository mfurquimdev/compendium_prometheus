package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {
  // Creates new/empty Registry
  upTimeReg := prometheus.NewRegistry()

  // A counter metric for how long (in seconds) the server is up
  var upTime = prometheus.NewCounter(
    prometheus.CounterOpts{
      Name: "uptime_seconds_total",
      Help: "Time in seconds the service is up",
    },
  )

  // Register metric in the Registry
  upTimeReg.MustRegister(upTime)

	// Increment variable each second inside a goroutine (parallel)
	go func() {
		for {
			time.Sleep(time.Second)
			upTime.Inc()
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
  http.Handle("/metrics-counter", promhttp.HandlerFor(upTimeReg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
