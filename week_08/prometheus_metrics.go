package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (

	// my timings metrics: method: /some/url, timing: 100500 ms
	timings = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "method_timing",
			Help: "Per method timing", // time to process each endpoint
		},
		[]string{"method"}, // method will be = url
	)

	// my counters metrics: method: /some/url, count: 100500 times
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "method_counter",
			Help: "Per method counter", // call count for each endpoint
		},
		[]string{"method"}, // url aka endpoint
	)
)

func init() {
	// hook-up my accumulators
	prometheus.MustRegister(timings)
	prometheus.MustRegister(counter)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	// work imitation
	rnd := time.Duration(rand.Intn(50))
	time.Sleep(time.Millisecond * rnd)
	w.Write([]byte("hello world"))
}

func timeTrackingMiddleware(next http.Handler) http.Handler {
	// decorator, update metrics accumulators
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r) // do work

		// r.URL.Path приходит от юзера! не делайте так в проде!
		timings.
			WithLabelValues(r.URL.Path).
			Observe(float64(time.Since(start).Seconds()))

		counter.
			WithLabelValues(r.URL.Path).
			Inc()
	})
}

func main() {
	siteMux := http.NewServeMux()

	// business
	siteMux.HandleFunc("/", mainPage)

	// metrics
	siteMux.Handle("/metrics", promhttp.Handler())

	// decorate (metrics endpoint will update my accum, do I want it?)
	siteHandler := timeTrackingMiddleware(siteMux)

	fmt.Println("starting server at :8083")
	http.ListenAndServe(":8083", siteHandler)
}
