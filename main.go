// A simple Prometheus exporter that can be used to temporarily increment the `livepeer_value_redeemed`
// vector counter metric with a given value and publish it on http://localhost:7935/metrics. This can
// be used to correct the `livepeer_value_redeemed` metric due to a bug that was present in the `go-livepeer`
// binary before https://github.com/livepeer/go-livepeer/pull/2916 was merged.
package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const WAIT_TIME = 60

// Create `livepeer_value_redeemed` counter vector metric.
var counterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "livepeer_value_redeemed",
		Help: " Winning ticket value redeemed",
	},
	[]string{"node_id", "node_type"},
)

func main() {
	log.Println("Starting `livepeer_value_redeemed` exporter...")

	// Retrieve increment value and node ID from command line arguments.
	valuePtr := flag.String("value", "", "Value to be published")
	nodeIDPtr := flag.String("node_id", "", "Node ID")

	// If -port is provided, use it to override the default port.
	portPtr := flag.String("port", "7935", "Port to be used")

	// Parse the command line arguments.
	flag.Parse()

	if *valuePtr == "" || *nodeIDPtr == "" {
		log.Fatal("Both value and node_id must be provided.")
	}
	log.Println("Value published to `livepeer_value_redeemed` metric:", *valuePtr)
	log.Println("Node ID label to be used:", *nodeIDPtr)

	// Register counter vector metric and initialize it with a value of 0.
	log.Println("Registering `livepeer_value_redeemed` counter vector metric with a value of 0...")
	prometheus.MustRegister(counterVec)
	counter := counterVec.With(prometheus.Labels{"node_id": *nodeIDPtr, "node_type": "orch"})
	counter.Add(0)

	// Create a goroutine that will increment the counter with the provided value after 60 seconds.
	go func() {
		// Wait for 60 seconds.
		log.Println("Waiting for 60 seconds...")
		time.Sleep(WAIT_TIME * time.Second)

		// Convert value to float64 and increment counter.
		parsedValue, err := strconv.ParseFloat(*valuePtr, 64)
		if err != nil {
			log.Printf("Failed to convert value to float64: %v", err)
			return
		}
		counter.Add(parsedValue)
		log.Println("The `livepeer_value_redeemed` counter vector was incremented with", parsedValue)
	}()

	// Start HTTP server and expose metrics.
	log.Println("Metric published on port:", *portPtr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":"+*portPtr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
