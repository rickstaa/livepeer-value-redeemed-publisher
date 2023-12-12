// A simple Prometheus exporter that can be used to temporarily increment the `livepeer_value_redeemed`
// vector counter metric with a given value and publish it on http://localhost:7936/metrics. This can
// be used to correct the `livepeer_value_redeemed` metric due to a bug that was present in the `go-livepeer`
// binary before https://github.com/livepeer/go-livepeer/pull/2916 was merged.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create `livepeer_value_redeemed` counter vector metric.
var counterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "livepeer_value_redeemed",
		Help: "This counter shows the value redeemed",
	},
	[]string{"node_id", "node_type"},
)

func main() {
	// Retrieve increment value and node ID from command line arguments.
	valuePtr := flag.String("value", "", "Value to be published")
	nodeIDPtr := flag.String("node_id", "", "Node ID")
	flag.Parse()
	if *valuePtr == "" || *nodeIDPtr == "" {
		log.Fatal("Both value and node_id must be provided.")
	}
	fmt.Println("Value to be published to livepeer_value_redeemed:", *valuePtr)
	fmt.Println("Node ID label to be used:", *nodeIDPtr)

	// Register counter vector metric and initialize it with a value of 0.
	prometheus.MustRegister(counterVec)
	counter := counterVec.With(prometheus.Labels{"node_id": *nodeIDPtr, "node_type": "orch"})
	counter.Add(0)

	// Create a goroutine that will increment the counter with the provided value after 5 seconds.
	go func() {
		time.Sleep(5 * time.Second)

		// Convert value to float64 and increment counter.
		parsedValue, err := strconv.ParseFloat(*valuePtr, 64)
		if err != nil {
			log.Fatal("Failed to convert value to float64.")
		}
		counter.Add(parsedValue)
		log.Println("Value published to livepeer_value_redeemed:", *valuePtr)
	}()

	// Start HTTP server and expose metrics.
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7936", nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
