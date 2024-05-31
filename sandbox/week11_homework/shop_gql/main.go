package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	gqlgen_extension "github.com/99designs/gqlgen/graphql/handler/extension"
)

func main() {
	panic("not yet")
}

func GetApp() http.Handler {
	gqlResolver := &Resolver{
		// UsersRepo:  usersRepo,
	}
	cfg := Config{
		Resolvers: gqlResolver,
	}

	srv := graphql_handler.NewDefaultServer(NewExecutableSchema(cfg))
	srv.Use(gqlgen_extension.FixedComplexityLimit(500))

	return srv
}

// --- http handler, empty -------------------------------------------

type ShopHttpApp struct {
	data any
}

var _ http.Handler = NewShopHttpApp() // type check

func NewShopHttpApp() *ShopHttpApp {
	return &ShopHttpApp{
		data: nil,
	}
}

// ServeHTTP implements http.Handler.
func (s *ShopHttpApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("ShopHttpApp.ServeHTTP, request: ", r.Method, r.URL.Path, r.Form)
	bodyBytes, err := io.ReadAll(r.Body)
	panicOnError("io.ReadAll failed", err)
	show(fmt.Sprintf("request body: %s", string(bodyBytes)))
	http.Error(w, "oops", http.StatusNotImplemented)
}

// --- useful little functions ---

var atomicCounter = new(atomic.Uint64)

func nextID_36() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 36)
}

func nextID_10() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 10)
}

func cutPrefix(s, prefix string) string {
	res, _ := strings.CutPrefix(s, prefix)
	return res
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

func strRef(in string) *string {
	return &in
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
