package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/metadata"
)

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func main() {
	// config_Flag_Json_Ldflags() // go run main.go --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"
	// go run -ldflags="-X 'main.Version=$(git rev-parse HEAD)' -X 'main.Branch=$(git rev-parse --abbrev-ref HEAD)'" main.go

	consulConfig()
}

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func consulConfig() {
	show("consulConfig: program started ...")
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

	loadConfig()
	go runConfigUpdater()

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)

	siteHandler := configMiddleware(siteMux)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, siteHandler)
	show("end of program. ", err)
}
func runConfigUpdater() {
	ticker := time.Tick(3 * time.Second)
	for _ = range ticker {
		loadConfig()
	}
}

func loadConfig() {
	qo := &consulapi.QueryOptions{
		WaitIndex: consulLastIndex,
	}
	kvPairs, qm, err := consul.KV().List(cfgPrefix, qo)
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
			continue
		}
		fmt.Printf("item[%d] %#v\n", idx, item)
		key := prefixStripper.Replace(item.Key)
		newConfig[key] = string(item.Value)
	}

	globalCfgMu.Lock()
	globalCfg = newConfig
	consulLastIndex = qm.LastIndex
	globalCfgMu.Unlock()

	fmt.Printf("config updated to version %v\n\t%#v\n\n", consulLastIndex, newConfig)
}
func loadPostsHandle(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	localCfg, err := ConfigFromContext(ctx)

	if err != nil {
		http.Error(w, "no config!", 500)
		return
	}

	fmt.Fprintf(w, "localCfg version %v\n%#v\n", consulLastIndex, localCfg)
	fmt.Fprintln(w, "Request done")
}
func ConfigFromContext(ctx context.Context) (map[string]string, error) {
	cfg, ok := ctx.Value(configKey).(map[string]string)
	if !ok {
		return nil, fmt.Errorf("config not found")
	}
	return cfg, nil
}

type key string

const configKey key = "configKey"

var (
	// consulAddr   = flag.String("addr", "192.168.99.100:32769", "consul addr (8500 in original consul)")
	consulAddr      = flag.String("addr", "127.0.0.1:8500", "consul addr (8500 in original consul)")
	consul          *consulapi.Client
	consulLastIndex uint64 = 0
	globalCfg              = make(map[string]string)
	globalCfgMu            = &sync.RWMutex{}
	cfgPrefix              = "myapi/"
	prefixStripper         = strings.NewReplacer(cfgPrefix, "")
)

func configMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		globalCfgMu.RLock()
		localCfg := globalCfg
		globalCfgMu.RUnlock()

		ctx = context.WithValue(ctx,
			configKey,
			localCfg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// go run flag.go --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"
var (
	commentsEnabled  = flag.Bool("comments", false, "Enable comments after post")
	commentsLimit    = flag.Int("limit", 10, "Comments number per page")
	commentsServices = &AddrList{}

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

func config_Flag_Json_Ldflags() {
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
