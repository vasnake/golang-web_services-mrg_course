package main

import (
	"fmt"
	"net/http"
	"runtime"

	"expvar"
)

var (
	hits = expvar.NewMap("hits") // register my metrics
)

// page handler, demo
func handler(w http.ResponseWriter, r *http.Request) {

	hits.Add(r.URL.Path, 1) // update my metrics

	w.Write([]byte("expvar increased"))
}

func init() {

	// hook-up my function to metrics output
	expvar.Publish("mystat", expvar.Func(func() interface{} {

		hits.Init() // reset my metrics

		// add another set of my metrics
		return map[string]int{
			"test":          100500,
			"value":         42,
			"goroutine_num": runtime.NumGoroutine(),
		}
	}))
}

// demo
func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8081")
	http.ListenAndServe(":8081", nil)
}
