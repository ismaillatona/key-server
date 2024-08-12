package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const baseURL string = "/key/"

var (
	maxSize            int
	srvPort            int
	keyLengthHistogram prometheus.Histogram
	httpStatusCounter  = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_status_codes",
			Help: "Counter of HTTP status codes",
		},
		[]string{"code"},
	)
)

func main() {
	// Define command-line flags
	flag.IntVar(&maxSize, "max-size", 1024, "maximum key size")
	flag.IntVar(&srvPort, "srv-port", 1123, "server listening port")
	flag.Parse()

	// Initialize Histogram bucket using initialized var max-size
	keyLengthHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "key_length_distribution",
		Help:    "Length distribution of generated keys",
		Buckets: prometheus.LinearBuckets(float64(maxSize)/20, float64(maxSize)/20, 20),
	})
	// Register Prometheus metrics
	prometheus.MustRegister(keyLengthHistogram)
	prometheus.MustRegister(httpStatusCounter)

	// Define the HTTP handler for /key/
	http.HandleFunc(baseURL, keyHandler)
	// Expose the registered metrics at /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start the HTTP server
	log.Printf("Server listening on port %d\n", srvPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", srvPort), nil))
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the requested length from the URL path
	lengthStr := r.URL.Path[len(baseURL):]
	if lengthStr == "" {
		http.Error(w, "Length must be specified", http.StatusBadRequest)
		httpStatusCounter.WithLabelValues("400").Inc() // Increment the counter for status code 400
		return
	}

	// Convert length to integer
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		http.Error(w, "Invalid length", http.StatusBadRequest)
		httpStatusCounter.WithLabelValues("400").Inc() // Increment the counter for status code 400
		return
	}

	// Check if the requested length exceeds the maximum size
	if length > maxSize {
		http.Error(w, "Requested length exceeds maximum size", http.StatusBadRequest)
		httpStatusCounter.WithLabelValues("400").Inc() // Increment the counter for status code 400
		return
	}

	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err = rand.Read(randomBytes)
	if err != nil {
		http.Error(w, "Failed to generate random bytes", http.StatusInternalServerError)
		httpStatusCounter.WithLabelValues("500").Inc() // Increment the counter for status code 500
		return
	}

	// Convert random bytes to a hexadecimal string
	response := hex.EncodeToString(randomBytes)

	// Record the length in the histogram
	keyLengthHistogram.Observe(float64(length))

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

	// Increment the counter for status code 200
	httpStatusCounter.WithLabelValues("200").Inc()
}
