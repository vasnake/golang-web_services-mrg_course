package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// one metric data
type Timing struct {
	Count    int
	Duration time.Duration
}

// all metrics for given request
type ctxTimings struct {
	sync.Mutex
	Data map[string]*Timing
}

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type key int

const timingsKey key = 1 // key for access to context values

// сколько в среднем спим при эмуляции работы
const AvgSleep = 50

func trackContextTimings(ctx context.Context, metricName string, start time.Time) {
	// update some metric by metric name, called using defer in request processing handler

	// получаем тайминги из контекста
	// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}

	elapsed := time.Since(start)

	// лочимся на случай конкурентной записи в мапку
	timings.Lock()
	defer timings.Unlock()

	// если меткри ещё нет - мы её создадим, если есть - допишем в существующую
	if metric, metricExist := timings.Data[metricName]; !metricExist {
		timings.Data[metricName] = &Timing{
			Count:    1,
			Duration: elapsed,
		}
	} else {
		metric.Count++
		metric.Duration += elapsed
	}
}

func logContextTimings(ctx context.Context, path string, start time.Time) {
	// calc and log metrics summary after all work is done,
	// called using defer in request processing handler

	// получаем тайминги из контекста
	// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}

	totalReal := time.Since(start)

	buf := bytes.NewBufferString(path)
	var total time.Duration // zero by default!
	for timing, value := range timings.Data {
		total += value.Duration
		buf.WriteString(fmt.Sprintf("\n\t%s(%d): %s", timing, value.Count, value.Duration))
	}

	buf.WriteString(fmt.Sprintf("\n\ttotal: %s", totalReal))
	buf.WriteString(fmt.Sprintf("\n\ttracked: %s", total))
	buf.WriteString(fmt.Sprintf("\n\tunkn: %s", totalReal-total))

	fmt.Println(buf.String())
}

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new metrics storage
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			timingsKey,
			&ctxTimings{
				Data: make(map[string]*Timing),
			})

		// log metrics after
		defer logContextTimings(ctx, r.URL.Path, time.Now())
		// process request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func emulateWork(ctx context.Context, workName string) {
	defer trackContextTimings(ctx, workName, time.Now())

	rnd := time.Duration(rand.Intn(AvgSleep))
	time.Sleep(time.Millisecond * rnd)
}

func loadPostsHandle(w http.ResponseWriter, req *http.Request) {
	// page handler, process with metrics collection
	ctx := req.Context()

	emulateWork(ctx, "checkCache")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")

	// emulate some unregistered work
	time.Sleep(10 * time.Millisecond)

	emulateWork(ctx, "loadSidebar")
	emulateWork(ctx, "loadComments")

	fmt.Fprintln(w, "Request done")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)

	// wrap with timings. It's somehow stupid: mtrics logic distributed across middleware and business logic
	siteHandler := timingMiddleware(siteMux)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", siteHandler)
}
