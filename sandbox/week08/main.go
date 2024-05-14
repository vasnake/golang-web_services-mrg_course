package main

import (
	"bytes"
	"context"
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"unsafe"

	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/metadata"
)

/*
// Всё что в комментари над `import "C"`` является кодом на C code и будет скомпилирован при помощи GCC.
// У вас должен быть установлен GCC
// see files *.c

void Multiply(int a, int b);
int MMultiply(int a, int b);
void foo();
void free(void*);
unsigned int sleep (unsigned int __seconds);
*/
import "C" //это псевдо-пакет, он реализуется компилятором
//export printResultGolang
func printResultGolang(result C.int) {
	fmt.Printf("result-var internals %T = %+v\n", result, result)
}

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func main() {
	// config_Flag_Json_Ldflags_Demo()
	// go run main.go --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"
	// go run -ldflags="-X 'main.Version=$(git rev-parse HEAD)' -X 'main.Branch=$(git rev-parse --abbrev-ref HEAD)'" main.go

	// consulConfigDemo()

	// expvarMetricsPullDemo()
	// graphiteMetricsPushDemo()
	// prometheusSimpleDemo()
	// prometheusDemo()

	// unsafeDemo()

	// cgoBasicDemo()
	// cgoAndBackDemo()
	// cgoMemLeakDemo()
	cgoThreadsHog()
}

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func cgoThreadsHog() {
	show("cgoThreadsHog: program started ...")
	/*
		в golang sleep не блокирует системный тред
		этот пример надо смотреть перед cgo_sleep
		после запуска надо посомтреть сколько тредов запущено процессом
	*/

	var go_main = func() {
		for i := 0; i < 100; i++ {
			go func() {
				// запускаем ГОшный sleep
				time.Sleep(time.Minute * 3) // num threads = num cores
			}()
		}
		// time.Sleep(time.Minute * 11)
		userInput("hit ENTER") // nolint
	}

	/*
		cgo блокирует весь системный тред
		golang-рантайм не может больше в нём запускать никакие другие горутины
		если в СИ операция была блокирующей, например, sleep, то она заблокирует весь тред
		после запуска надо посомтреть сколько тредов запущено процессом
		будет больше чем в cgo_go_sleep
	*/
	var cgo_main = func() {
		for i := 0; i < 100; i++ {
			go func() {
				// запускаем СИшный sleep
				C.sleep(60 * 3)
			}()
		}
		// time.Sleep(11 * time.Minute)
		userInput("hit ENTER") // nolint
	}

	go_main()
	cgo_main()

	show("end of program. ")
}

func cgoMemLeakDemo() {
	show("cgoMemLeakDemo: program started ...")

	var print = func(s string) {
		cs := C.CString(s) // переход в другую вселенную

		// СИ-шные не собираются через ГО-шный сборщик мусора, их надо освобождать руками
		// закомментируйте эту строку и запустите программу - начнётся утечка памяти
		defer C.free(unsafe.Pointer(cs))

		println(cs)
	}

	for i := 0; i < 10; i++ {
		print("Hello World")
	}

	show("end of program. ")
}

// cgo overhead bench
func CallCgo(n int) {
	for i := 0; i < n; i++ {
		C.foo()
	}
}
func CallGo(n int) {
	for i := 0; i < n; i++ {
		foo()
	}
}
func foo() {}

// BenchmarkCGO-8          24312007                43.25 ns/op
// BenchmarkGo-8           1000000000               0.7206 ns/op

func cgoAndBackDemo() {
	show("cgoAndBackDemo: program started ...")

	a := 2
	b := 3
	// для того чтобы вызвать СИшный крод надо добавить префикс "С."
	// там же туда надо передать СИшные переменные
	C.Multiply(C.int(a), C.int(b)) // void

	show("end of program. ")
	/*
	   2024-05-13T10:46:47.336Z: cgoAndBackDemo: program started ...
	   result-var internals main._Ctype_int = 6
	   2024-05-13T10:46:47.336Z: end of program.
	*/
}

func cgoBasicDemo() {
	show("cgoBasicDemo: program started ...")

	a := 2
	b := 3
	// для того чтобы вызвать СИшный крод надо добавить префикс "С."
	// там же туда надо передать СИшные переменные
	res := C.MMultiply(C.int(a), C.int(b))
	fmt.Printf("Multiply in C: %d * %d = %d\n", a, b, int(res))
	fmt.Printf("с-var internals %T = %+v\n", res, res)

	show("end of program. ")
	/*
	   2024-05-13T10:27:47.273Z: cgoBasicDemo: program started ...
	   Multiply in C: 2 * 3 = 6
	   с-var internals main._Ctype_int = 6
	   2024-05-13T10:27:47.273Z: end of program.
	*/
}

func unsafeDemo() {
	show("unsafeDemo: program started ...")

	var Float64bits = func(f float64) uint64 {
		return *(*uint64)(unsafe.Pointer(&f)) // take ref, make a pointer, set a new interpretation (type) for memory, un-ref value
	}

	a := int64(1)
	fmt.Println("memory pointer for var `a`", unsafe.Pointer(&a))
	fmt.Println("memory size for var `a`", unsafe.Sizeof(a))

	println("-------")

	f := 10.11
	fmt.Printf("10.11 float64 in dec: %d\n", Float64bits(f))
	fmt.Printf("in hex: %#016x\n", Float64bits(f))
	fmt.Printf("in binary: %b\n", Float64bits(f))

	println("-------")

	type Message struct {
		flag1 bool
		flag2 bool
		name  string
	}

	msg := Message{ // nolint
		flag1: false,
		flag2: true,
		name:  "Neque porro quisquam est qui dolorem",
	}

	fmt.Println("memory size for Message struct", unsafe.Sizeof(msg))

	fmt.Println(
		"flag1 Sizeof", unsafe.Sizeof(msg.flag1),
		"Alignof", unsafe.Alignof(msg.flag1),
		"Offsetof", unsafe.Offsetof(msg.flag1),
	)

	fmt.Println(
		"flag2 Sizeof", unsafe.Sizeof(msg.flag2),
		"Alignof", unsafe.Alignof(msg.flag2),
		"Offsetof", unsafe.Offsetof(msg.flag2),
	)

	fmt.Println(
		"name Sizeof", unsafe.Sizeof(msg.name),
		"Alignof", unsafe.Alignof(msg.name),
		"Offsetof", unsafe.Offsetof(msg.name),
	)

	// bytesToStr сорздаёт строку, указывающую на слайс байт, чтобы избежать копирования.
	// Warning: the string returned by the function should be used with care, as the whole input data
	// chunk may be either blocked from being freed by GC because of a single string or the buffer.Data
	// may be garbage-collected even when the string exists.
	var s = "0"
	var bytesToStr = func(data []byte) string {
		sliceHeaderRef := (*reflect.SliceHeader)(unsafe.Pointer(&data)) // nolint
		fmt.Printf("type: %T, value: %+v\n", sliceHeaderRef, sliceHeaderRef)
		fmt.Printf("type: %T, value: %+v\n", sliceHeaderRef.Data, sliceHeaderRef.Data)
		// stringHeaderRef := reflect.StringHeader{Data: sliceHeaderRef.Data, Len: sliceHeaderRef.Len}
		stringHeaderRef := (*reflect.StringHeader)(unsafe.Pointer(&s)) // nolint
		stringHeaderRef.Data = sliceHeaderRef.Data
		stringHeaderRef.Len = sliceHeaderRef.Len
		fmt.Printf("type: %T, value: %+v\n", stringHeaderRef, stringHeaderRef)
		return *(*string)(unsafe.Pointer(stringHeaderRef))
	}

	data := []byte(`some test`)
	show("change header from byte_slice to string", data)
	str := bytesToStr(data)
	// str := string(data)
	fmt.Printf("type: %T, value: %+v\n", str, str)

	show("end of program. ")
	/*
	   2024-05-13T09:55:28.105Z: unsafeDemo: program started ...
	   memory pointer for var `a` 0xc00020e010
	   memory size for var `a` 8
	   -------
	   10.11 float64 in dec: 4621881042083847864
	   in hex: 0x40243851eb851eb8
	   in binary: 100000000100100001110000101000111101011100001010001111010111000
	   -------
	   memory size for Message struct 24
	   flag1 Sizeof 1 Alignof 1 Offsetof 0
	   flag2 Sizeof 1 Alignof 1 Offsetof 1
	   name Sizeof 16 Alignof 8 Offsetof 8
	   2024-05-13T09:55:28.105Z: change header from byte_slice to string[]byte{0x73, 0x6f, 0x6d, 0x65, 0x20, 0x74, 0x65, 0x73, 0x74};
	   type: *reflect.SliceHeader, value: &{Data:824635875376 Len:9 Cap:9}
	   type: uintptr, value: 824635875376
	   type: *reflect.StringHeader, value: &{Data:824635875376 Len:9}
	   type: string, value: some test
	   2024-05-13T09:55:28.105Z: end of program.
	*/
}

func prometheusDemo() {
	show("prometheusDemo: program started ...")

	var mainPage = func(w http.ResponseWriter, r *http.Request) {
		rnd := time.Duration(rand.Intn(50)) // work imitation
		time.Sleep(time.Millisecond * rnd)
		w.Write([]byte("hello world")) // nolint
	}

	var timeTrackingMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r) // serve

			// r.URL.Path приходит от юзера! не делайте так в проде! (unsafe)
			// update registered custom metrics
			pmts_req_timings.
				WithLabelValues(r.URL.Path).
				Observe(float64(time.Since(start).Seconds()))
			pmts_req_counter.
				WithLabelValues(r.URL.Path).
				Inc()
		})
	}

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", mainPage)
	siteMux.Handle("/metrics", promhttp.Handler()) // http://localhost:8080/metrics

	siteHandler := timeTrackingMiddleware(siteMux) // collect metrics for prometheus

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, siteHandler)
	show("end of program. ", err)
}

var (
	pmts_req_timings = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "method_timing",
			Help: "Per method timing",
		},
		[]string{"method"},
	)
	pmts_req_counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "method_counter",
			Help: "Per method counter",
		},
		[]string{"method"},
	)
)

func init() {
	// add custom metrics to other numbers
	prometheus.MustRegister(pmts_req_timings)
	prometheus.MustRegister(pmts_req_counter)
}

func prometheusSimpleDemo() {
	show("prometheusSimpleDemo: program started ...")

	http.Handle("/metrics", promhttp.Handler()) // server process metrics
	// http://localhost:8080/metrics

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func graphiteMetricsPushDemo() {
	/*
		docker copmpose up
		or

		docker run -d\
		--name graphite\
		--restart=always\
		-p 80:80\
		-p 2003-2004:2003-2004\
		-p 2023-2024:2023-2024\
		-p 8125:8125/udp\
		-p 8126:8126\
		graphiteapp/graphite-statsd

		http://localhost/?width=799&height=395&target=coursera.goroutines_num&target=coursera.mem_heap&target=coursera.mem_stack
	*/

	show("graphiteMetricsPushDemo: program started ...")
	flag.Parse()

	var sendStat = func() { // async
		conn, err := net.Dial("tcp", *carbonAddr)
		if err != nil {
			panic(err)
		}

		memstats := &runtime.MemStats{}
		c := time.Tick(3 * time.Second)
		for tickTime := range c {
			runtime.ReadMemStats(memstats)

			buf := bytes.NewBuffer([]byte{})
			fmt.Fprintf(buf, "coursera.mem_heap %d %d\n",
				memstats.HeapInuse, tickTime.Unix())
			fmt.Fprintf(buf, "coursera.mem_stack %d %d\n",
				memstats.StackInuse, tickTime.Unix())
			fmt.Fprintf(buf, "coursera.goroutines_num %d %d\n",
				runtime.NumGoroutine(), tickTime.Unix())

			conn.Write(buf.Bytes()) // nolint
			fmt.Println(buf.String())
		}
	}

	go sendStat()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world")) // nolint
	})

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

var carbonAddr = flag.String("graphite", "localhost:2003", "The address of carbon receiver")

func expvarMetricsPullDemo() {
	show("expvarMetricsPullDemo: program started ...")
	// по урл `/debug/vars` нам доступны некоторые показания
	// http://localhost:8080/debug/vars

	var handler = func(w http.ResponseWriter, r *http.Request) {
		hits_expvarCustomMetrics.Add(r.URL.Path, 1) // update metrics
		w.Write([]byte("expvar increased"))         // nolint
	}

	http.HandleFunc("/", handler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

var (
	hits_expvarCustomMetrics = expvar.NewMap("hits") // just metrics data
)

func init() {
	// expvar custom metrics publisher // some logic to run on each metrics call
	expvar.Publish("mystat", expvar.Func(func() interface{} {
		hits_expvarCustomMetrics.Init() // clear previous data from custom metrics collection
		return map[string]int{
			"test":          100500, // imitation
			"value":         42,
			"goroutine_num": runtime.NumGoroutine(),
		}
	}))
}

func consulConfigDemo() {
	show("consulConfigDemo: program started ...")
	// sandbox\week08\docker-compose.yml => pushd week08 && docker compose up&
	// http://localhost:8500/ui/dc1/kv/myapi/
	// config updated to version 46 map[string]string{"baz":"toodaloo", "foo":"bar"}

	flag.Parse()
	var err error
	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err = consulapi.NewClient(config)

	if err != nil {
		fmt.Println("consul error", err)
		return
	}

	go runConfigUpdater() // update cfg from consul, every x sec

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", showCfgWebPageHandler)

	siteHandler := addConfigMiddleware(siteMux) // add copy of current cfg to request

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, siteHandler)
	show("end of program. ", err)
}
func runConfigUpdater() {
	updateConfigFromConsul()
	ticker := time.Tick(3 * time.Second)
	for _ = range ticker { // nolint
		updateConfigFromConsul()
	}
}

func updateConfigFromConsul() {
	opts := &consulapi.QueryOptions{
		WaitIndex: consulLastIndex,
	}
	kvPairs, qm, err := consul.KV().List(cfgPrefix, opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("remote consulLastIndex", qm.LastIndex)
	if consulLastIndex == qm.LastIndex {
		fmt.Println("consulLastIndex not changed")
		return
	}

	newConfig := make(map[string]string)

	for idx, item := range kvPairs {
		if item.Key == cfgPrefix {
			continue // skip cfg folder item
		}
		fmt.Printf("item[%d] %#v\n", idx, item)
		key := prefixStripper.Replace(item.Key)
		newConfig[key] = string(item.Value)
	}

	globalCfgMu.Lock()
	globalCfg = newConfig // don't mutate the map, make a ref to a new map instead
	consulLastIndex = qm.LastIndex
	globalCfgMu.Unlock()

	fmt.Printf("config updated to version %v\n\t%#v\n\n", consulLastIndex, newConfig)
}
func showCfgWebPageHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	localCfg, err := getConfigFromContext(ctx)

	if err != nil {
		http.Error(w, "no config!", 500)
		return
	}

	fmt.Fprintf(w, "localCfg version %v\n%#v\n", consulLastIndex, localCfg) // no, global cfg version here
	fmt.Fprintln(w, "Request done")
}
func getConfigFromContext(ctx context.Context) (map[string]string, error) {
	cfg, ok := ctx.Value(configKey).(map[string]string)
	if !ok {
		return nil, fmt.Errorf("config not found")
	}
	return cfg, nil
}

type key string

const configKey key = "configKey"

var ( // consul config infrastructure
	// consulAddr   = flag.String("addr", "192.168.99.100:32769", "consul addr (8500 in original consul)")
	consulAddr      = flag.String("addr", "127.0.0.1:8500", "consul addr (8500 in original consul)")
	consul          *consulapi.Client
	consulLastIndex uint64 = 0
	globalCfg              = make(map[string]string) // global mutable config
	globalCfgMu            = &sync.RWMutex{}
	cfgPrefix              = "myapi/"
	prefixStripper         = strings.NewReplacer(cfgPrefix, "")
)

func addConfigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		globalCfgMu.RLock()
		localCfg := globalCfg // make a local copy from global mutable config
		globalCfgMu.RUnlock()

		ctx = context.WithValue(ctx, configKey, localCfg)

		next.ServeHTTP(w, r.WithContext(ctx)) // serve
	})
}

func config_Flag_Json_Ldflags_Demo() {
	// go run flag.go --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"
	show("config_Flag_Json_Ldflags: program started ...")

	// flag
	flag.Parse()
	if *commentsEnabled {
		fmt.Println("Comments per page", *commentsLimit)
		fmt.Println("Comments services", *commentsServices)
	} else {
		fmt.Println("Comments disabled")
	}

	// json
	data, err := os.ReadFile("./week08/config.json")
	if err != nil {
		log.Fatalln("can't read config file:", err)
	}
	type Config struct {
		Comments bool `json:"comments"`
		Limit    int
		Servers  []string
	}
	var config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("can't parse config:", err)
	}
	if config.Comments {
		fmt.Println("Comments per page", config.Limit) // default 0 as int type
		fmt.Println("Comments services", config.Servers)
	} else {
		fmt.Println("Comments disabled")
	}

	// ldflags
	// go run -ldflags="-X 'main.Version=$(git rev-parse HEAD)' -X 'main.Branch=$(git rev-parse --abbrev-ref HEAD)'" main.go
	fmt.Printf("[ldflags] starting version `%s`, branch `%s` \n", Version, Branch)

	show("end of program. ")
	/*
		2024-05-10T12:18:33.858Z: config_Flag_Json_Ldflags: program started ...
		Comments per page 10
		Comments services [127.0.0.1:8081 127.0.0.1:8082]
		Comments per page 0
		Comments services [127.0.0.1:8081 127.0.0.1:8082]
		[ldflags] starting version `54640b14239c9b5cf894f108a3fb207c56e782bd`, branch `main`
		2024-05-10T12:18:33.861Z: end of program.
	*/
}

// go run flag.go --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"
var (
	commentsEnabled  = flag.Bool("comments", false, "Enable comments after post")
	commentsLimit    = flag.Int("limit", 10, "Comments number per page")
	commentsServices = &AddrList{} // custom flag

	// go run -ldflags="-X 'main.Version=$(git rev-parse HEAD)' -X 'main.Branch=$(git rev-parse --abbrev-ref HEAD)'" main.go
	Version = "" // ldflags
	Branch  = "" // ldflags
)

func init() {
	flag.Var(commentsServices, "servers", "Comments number per page")
}

type AddrList []string

func (v *AddrList) String() string {
	return fmt.Sprint(*v)
}
func (v *AddrList) Set(in string) error {
	for _, addr := range strings.Split(in, ",") {
		ipRaw, _, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("bad addr %v", addr)
		}
		ip := net.ParseIP(ipRaw)
		if ip.To4() == nil {
			return fmt.Errorf("invalid ipv4 addr %v", addr)
		}
		*v = append(*v, addr)
	}
	return nil
}

// GetContextMetadataValue returns value under given key from context metadata.
// If no data in context or value is empty, return `dflt` value
func GetContextMetadataValue(ctx context.Context, key, dflt string) string {
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		show("GetContextMetadataValue failed, context w/o metadata")
		return dflt
	}
	x := strings.Join(md.Get(key), "")
	if x == "" {
		return dflt
	}
	return x
}

func userInput(msg string) (res string, err error) {
	show(msg)
	if n, e := fmt.Scanln(&res); n != 1 || e != nil {
		return "", e
	}
	return res, nil
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
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
