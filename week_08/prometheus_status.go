package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// create prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", nil)
}
