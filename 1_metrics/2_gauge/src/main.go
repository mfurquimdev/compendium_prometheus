package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Creates new/empty Registry
	funcDurationReg := prometheus.NewRegistry()

	// A counter metric for how long (in seconds) the server is up
	var funcDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "function_duration_seconds",
			Help: "Time in seconds the most recently run of a function has taken to complete",
		},
	)

	// Register metric in the Registry
	funcDurationReg.MustRegister(funcDuration)

  // Execute a function in a goroutine to alter the variable funcDuration
  // based on a random time waiting
	go func() {
		for {
			go func() {
				timer := prometheus.NewTimer(prometheus.ObserverFunc(funcDuration.Set))
				defer timer.ObserveDuration()
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			}()
			time.Sleep(5 * time.Second)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics-gauge", promhttp.HandlerFor(funcDurationReg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
