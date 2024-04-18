package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	html_tmpl "html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	gorilla_mux "github.com/gorilla/mux"
	http_router "github.com/julienschmidt/httprouter"
	pkg_err "github.com/pkg/errors"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	go_valid "github.com/asaskevich/govalidator"
	gorilla_schema "github.com/gorilla/schema"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	gorilla_ws "github.com/gorilla/websocket"
	"github.com/icrowley/fake"

	tmpl_item "week05/item"     // separate lib, used in template
	hero_tmpl "week05/template" // generated template
)

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func main() {
	// middlwareDemo()
	// contextValueDemo()

	// basicErrorDemo()
	// namedErrorDemo()
	// ownErrorDemo()
	// pkgErrorDemo() // preferred error handling

	// gorillaRouterDemo() // powerful but slow
	// httprouterDemo() // middle-class
	// mixedRoutersDemo() // multiple
	// fasthttpDemo() // fast but ugly

	// validationDemo()
	// loggingDemo()
	// websocketDemo()
	heroTemplatesDemo()
}

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

//go:generate hero -source=./template/
func heroTemplatesDemo() {
	show("heroTemplatesDemo: program started ...")

	var ExampleItems = []*tmpl_item.Item{
		{Id: 1, Title: "foo", Description: "bar"},
		{Id: 2, Title: "username", Description: "freelancer"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		buffer := new(bytes.Buffer)
		hero_tmpl.Index(ExampleItems, buffer) // call to generated method
		w.Write(buffer.Bytes())
	})

	show("Starting server (/) at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func websocketDemo() {
	show("websocketDemo: program started ...")

	newMessage := func() []byte {
		data, _ := json.Marshal(map[string]string{
			"email":   fake.EmailAddress(),
			"name":    fake.FirstName() + " " + fake.LastName(),
			"subject": fake.Product() + " " + fake.Model(),
		})
		return data
	}

	pushWSNotifications := func(client *gorilla_ws.Conn) {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			w, err := client.NextWriter(gorilla_ws.TextMessage)
			if err != nil {
				break
			}
			w.Write(newMessage()) // json text
			w.Close()

			<-ticker.C // wait for next event
		}
	}

	var connUpgrader = gorilla_ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// setup web server

	tmpl := html_tmpl.Must(html_tmpl.ParseFiles("week05/ws.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// open from JS WS
	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		conn, err := connUpgrader.Upgrade(w, r, nil) // add WS proto
		if err != nil {
			log.Fatal(err)
		} else {
			go pushWSNotifications(conn) // async endless loop
		}
	})

	show("Starting server (/, /notifications) at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func loggingDemo() {
	show("loggingDemo: program started ...")

	addr := host
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: false})

	// у zap-а нет логгера по-умолчанию
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flush

	// start with printing some lines

	// stdout print
	fmt.Printf("STDOUT starting server at %s:%d", addr, port)
	// std logger (default)
	log.Printf("STDLOG starting server at %s:%d", addr, port)

	// zap
	zapLogger.Info("starting server",
		zap.String("logger", "ZAP"),
		zap.String("host", addr),
		zap.Int("port", port),
	) // msg and fields

	// logrus
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   addr,
		"port":   port,
	}).Info("Starting server")

	// build web server logger

	accessLog := new(ThreeLoggers)
	accessLog.StdLogger = log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile)
	accessLog.ZapLogger = zapLogger.Sugar().With(
		zap.String("mode", "[access_log]"),
		zap.String("logger", "ZAP"),
	)
	accessLog.LogrusLogger = logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// setup web server

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	})
	siteWithLogging := accessLog.accessLogMiddleware(siteMux) // log 5 lines, all different methods

	show("Starting server (/) at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, siteWithLogging)
	show("end of program. ", err)
	/*
		1 2024-04-18T07:56:50.590Z: loggingDemo: program started ...

		2 STDOUT starting server at 127.0.0.1:8080
		  2024/04/18 10:56:50 STDLOG starting server at 127.0.0.1:8080
		3 {"level":"info","ts":1713427010.590113,"caller":"week05/main.go:116","msg":"starting server","logger":"ZAP","host":"127.0.0.1","port":8080}
		4 INFO[0000] Starting server                               host=127.0.0.1 logger=LOGRUS port=8080

		5 2024-04-18T07:56:50.590Z: Starting server (/) at: string(127.0.0.1:8080);
		6 2024-04-18T07:56:50.590Z: Open url http://localhost:8080/

		7  FMT [GET] [::1]:57042, / 7.911µs
		8  2024/04/18 10:56:57 LOG [GET] [::1]:57042, / 45.876µs
		9  STD main.go:81: [GET] [::1]:57042, / 51.211µs
		10 {"level":"info","ts":1713427017.5010386,"caller":"week05/main.go:84","msg":"/{method 15 0 GET <nil>} {remote_addr 15 0 [::1]:57042 <nil>} {url 15 0 / <nil>} {work_time 8 75764  <nil>}","mode":"[access_log]","logger":"ZAP"}
		11 {"level":"info","logger":"LOGRUS","method":"GET","mode":"[access_log]","msg":"/","remote_addr":"[::1]:57042","time":"2024-04-18T10:56:57+03:00","work_time":125042}

		12 FMT [GET] [::1]:57042, /favicon.ico 22.876µs
		13 2024/04/18 10:56:57 LOG [GET] [::1]:57042, /favicon.ico 88.669µs
		14 STD main.go:81: [GET] [::1]:57042, /favicon.ico 107.216µs
		15 {"level":"info","ts":1713427017.5220897,"caller":"week05/main.go:84","msg":"/favicon.ico{method 15 0 GET <nil>} {remote_addr 15 0 [::1]:57042 <nil>} {url 15 0 /favicon.ico <nil>} {work_time 8 171958  <nil>}","mode":"[access_log]","logger":"ZAP"}
		16 {"level":"info","logger":"LOGRUS","method":"GET","mode":"[access_log]","msg":"/favicon.ico","remote_addr":"[::1]:57042","time":"2024-04-18T10:56:57+03:00","work_time":282405}
	*/
}

type ThreeLoggers struct {
	StdLogger    *log.Logger
	ZapLogger    *zap.SugaredLogger
	LogrusLogger *logrus.Entry
}

func (ac *ThreeLoggers) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // call internal handler

		// log line 1
		fmt.Printf("FMT [%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
		// log line 2
		log.Printf("LOG [%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
		// line 3
		ac.StdLogger.Printf("[%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
		// line 4
		ac.ZapLogger.Info(r.URL.Path,
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Duration("work_time", time.Since(start)),
		)
		// line 5
		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	}) // handler wrapper
}

func validationDemo() {
	show("validationDemo: program started ...")

	init := func() {
		go_valid.CustomTypeTagMap.Set("msgSubject", go_valid.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			subject, ok := i.(string) // 'subject' field, string
			if !ok {
				return false
			}
			if len(subject) < 1 || len(subject) > 10 { // 'subject' size limit
				return false
			}
			return true
		}))
	}
	init() // module init, custom validator registration

	type SendMessage struct {
		Id        int    `valid:",optional"` // valid: describe validation rules
		Priority  string `valid:"in(low|normal|high)"`
		Recipient string `schema:"to" valid:"email"` // schema: describe parsing rules
		Subject   string `valid:"msgSubject"`        // call custom validator
		Inner     string `schema:"-" valid:"-"`      // skip
		flag      int    // private, no schema, no validation
	}
	// http://localhost:8080/?to=fooATbar.com&priority=top&subject=Helloqqqqqqqqqqqqqq!&inner=ignored&id=12&flag=23

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("request as is: " + r.URL.String() + "\n\n"))

		msg := &SendMessage{}
		decoder := gorilla_schema.NewDecoder()
		decoder.IgnoreUnknownKeys(true)
		err := decoder.Decode(msg, r.URL.Query()) // fill msg struct from url, parse
		if err != nil {
			fmt.Println(err)
			http.Error(w, "failed to decode request parameters", 500)
			return
		}

		w.Write([]byte(fmt.Sprintf("Decoded message: %#v\n\n", msg)))

		_, err = go_valid.ValidateStruct(msg) // validate values
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Validator: malformed message: %s\n\n", err))) // cumulative error
			if allErrs, ok := err.(go_valid.Errors); ok {                             // get detailed list of errors
				for _, e := range allErrs.Errors() {
					data := []byte(fmt.Sprintf("malformed value, error: %#v\n\n", e))
					w.Write(data)
				}
			}
		} else {
			w.Write([]byte(fmt.Sprintf("Validator: message OK\n\n")))
		}
	}

	http.HandleFunc("/", handler)

	show("Starting server (/) at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func fasthttpDemo() {
	show("fasthttpDemo: program started ...")

	rootHandler := func(ctx *fasthttp.RequestCtx) {
		users := []string{"foo_bar"}
		body, _ := json.Marshal(users)
		// any order will do:
		ctx.SetBody(body)
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	getUserHandler := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "view user: `%s`\n", ctx.UserValue("id"))
	}

	router := fasthttprouter.New()
	// curl -v -X GET http://localhost:8080/
	router.GET("/", rootHandler)
	// curl -v -X GET http://localhost:8080/users/foo
	router.GET("/users/:id", getUserHandler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	log.Fatal(fasthttp.ListenAndServe(portStr, router.Handler))
	// err := http.ListenAndServe(addrStr, nil)
	// show("end of program. ", err)
}

func mixedRoutersDemo() {
	show("mixedRoutersDemo: program started ...")

	httpRouterHandler := func(w http.ResponseWriter, r *http.Request, ps http_router.Params) {
		fmt.Fprintf(w, "Request with high hitrate, id: `%s`\n", ps.ByName("id"))
	}
	gorillaHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Request with complex routing logic\n")
	}
	stdHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Request with averange complexity\n")
	}

	// curl -v -X GET http://localhost:8080/fast/foo
	fastApiHandler := http_router.New()
	fastApiHandler.GET("/fast/:id", httpRouterHandler)

	// curl -v -X GET -H "X-Requested-With: XMLHttpRequest" -d '{"name":"Foo Bar"}' http://localhost:8080/complex/
	complexApiHandler := gorilla_mux.NewRouter()
	complexApiHandler.HandleFunc("/complex/", gorillaHandler).Headers("X-Requested-With", "XMLHttpRequest") // ajax

	// curl -v -X GET http://localhost:8080/std/
	stdApiHandler := http.NewServeMux()
	stdApiHandler.HandleFunc("/std/", stdHandler)

	// combine 3 routers alltogether
	siteMux := http.NewServeMux()
	siteMux.Handle("/fast/", fastApiHandler)
	siteMux.Handle("/complex/", complexApiHandler)
	siteMux.Handle("/std/", stdApiHandler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, siteMux)
	show("end of program. ", err)
}

func httprouterDemo() {
	show("httprouterDemo: program started ...")

	List := func(w http.ResponseWriter, r *http.Request, _ http_router.Params) {
		fmt.Fprint(w, "users list ...\n")
	}

	Get := func(w http.ResponseWriter, r *http.Request, ps http_router.Params) {
		fmt.Fprintf(w, "view user by id: `%s`\n", ps.ByName("id"))
	}

	Create := func(w http.ResponseWriter, r *http.Request, ps http_router.Params) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		fmt.Fprintf(w, "creating a new user: `%s`\n", body)
	}

	Update := func(w http.ResponseWriter, r *http.Request, ps http_router.Params) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		fmt.Fprintf(w, "updating user: `%s`, new data: `%s`\n", ps.ByName("login"), body)
	}

	router := http_router.New()
	router.GET("/", List)
	// curl -v -X GET http://localhost:8080/users
	router.GET("/users", List)

	// curl -v -X GET http://localhost:8080/users/foo_bar
	router.GET("/users/:id", Get)

	// curl -v -X PUT -H "Content-Type: application/json" -d '{"login":"foo_bar"}' http://localhost:8080/users
	router.PUT("/users", Create)

	// curl -v -X POST -H "Content-Type: application/json" -d '{"name":"Foo Bar"}' http://localhost:8080/users/foo_bar
	router.POST("/users/:login", Update)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, router)
	show("end of program. ", err)
}

func gorillaRouterDemo() {
	show("gorillaRouterDemo: program started ...")

	listUsers := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "users list: ...\n")
	}

	userByID := func(w http.ResponseWriter, r *http.Request) {
		routeVars := gorilla_mux.Vars(r)
		fmt.Fprintf(w, "view user id: `%s`\n", routeVars["id"])
	}

	createUser := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		fmt.Fprintf(w, "creating a new user: `%s`\n", body)
	}

	updateUser := func(w http.ResponseWriter, r *http.Request) {
		routeVars := gorilla_mux.Vars(r)

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		fmt.Fprintf(w, "updating user: `%s`, new data: `%s`\n", routeVars["login"], body)
	}

	r := gorilla_mux.NewRouter()
	r.HandleFunc("/", listUsers)

	// curl -v -X GET http://localhost:8080/users/123
	r.HandleFunc("/users/{id:[0-9]+}", userByID)

	// curl -v -X PUT -H "Content-Type: application/json" -d '{"login":"foo_bar"}' http://localhost:8080/users
	r.HandleFunc("/users", createUser).Methods("PUT")

	// curl -v -X POST -H "Content-Type: application/json"  -H "X-Auth: test" -d '{"name":"Foo Bar"}' http://localhost:8080/users/foo_bar
	r.HandleFunc("/users/{login}", updateUser).Methods("POST").Headers("X-Auth", "test")

	// curl -v -X GET http://localhost:8080/users
	r.HandleFunc("/users", listUsers).Host("localhost").Methods("GET")

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, r)
	show("end of program. ", err)
}

func pkgErrorDemo() {
	show("pkgErrorDemo: program started ...")

	client := http.Client{Timeout: time.Duration(time.Millisecond)}
	getRemoteResource := func() error {
		url := "http://127.0.0.1:9999/pages?id=123"
		_, err := client.Get(url)
		if err != nil {
			return pkg_err.Wrap(err, "resource error") // stack trace captured here
		}
		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		err := getRemoteResource()
		if err != nil {
			fmt.Printf("full err: %+v\n", err)        // stack trace printed here
			switch err := pkg_err.Cause(err).(type) { // check wrapped error
			case *url.Error:
				fmt.Printf("resource %s err: %+v\n", err.URL, err.Err)
				http.Error(w, "remote resource error", 500)
			default:
				fmt.Printf("%+v\n", err)
				http.Error(w, "parsing error", 500)
			}
			return
		}
		w.Write([]byte("all is OK"))
	}

	http.HandleFunc("/", handler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
	/*
		2024-04-17T12:52:06.671Z: pkgErrorDemo: program started ...
		2024-04-17T12:52:06.671Z: Starting server at: string(127.0.0.1:8080);
		2024-04-17T12:52:06.671Z: Open url http://localhost:8080/
		full err: Get "http://127.0.0.1:9999/pages?id=123": dial tcp 127.0.0.1:9999: connect: connection refused
		resource error
		main.pkgErrorDemo.func1
		        prj/sandbox/week05/main.go:48
		main.pkgErrorDemo.func2
		        prj/sandbox/week05/main.go:54
		net/http.HandlerFunc.ServeHTTP
		        ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.2.linux-amd64/src/net/http/server.go:2166
		net/http.(*ServeMux).ServeHTTP
		        ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.2.linux-amd64/src/net/http/server.go:2683
		net/http.serverHandler.ServeHTTP
		        ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.2.linux-amd64/src/net/http/server.go:3137
		net/http.(*conn).serve
		        ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.2.linux-amd64/src/net/http/server.go:2039
		runtime.goexit
		        ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.2.linux-amd64/src/runtime/asm_amd64.s:1695
		resource http://127.0.0.1:9999/pages?id=123 err: dial tcp 127.0.0.1:9999: connect: connection refused
	*/
}

type ResourceError struct {
	URL string
	Err error
}

func (re *ResourceError) Error() string { // custom error type
	return fmt.Sprintf("Resource error: URL: %s, err: %v", re.URL, re.Err)
}
func ownErrorDemo() {
	show("ownErrorDemo: program started ...")

	client := http.Client{Timeout: time.Duration(time.Millisecond)}

	getRemoteResource := func() error {
		url := "http://127.0.0.1:9999/pages?id=123"
		_, err := client.Get(url)
		if err != nil {
			return &ResourceError{URL: url, Err: err} // custom error type
		}
		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		err := getRemoteResource()
		if err != nil {
			switch err.(type) { // check for known custom errors
			case *ResourceError:
				err := err.(*ResourceError)
				fmt.Printf("resource %s err: %s\n", err.URL, err.Err)
				http.Error(w, "remote resource error", 500)
			default:
				fmt.Printf("internal error: %+v\n", err)
				http.Error(w, "internal error", 500)
			}
			return
		}
		w.Write([]byte("all is OK"))
	}

	http.HandleFunc("/", handler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
	/*
		2024-04-17T12:39:05.813Z: ownErrorDemo: program started ...
		2024-04-17T12:39:05.813Z: Starting server at: string(127.0.0.1:8080);
		2024-04-17T12:39:05.813Z: Open url http://localhost:8080/
		resource http://127.0.0.1:9999/pages?id=123 err: Get "http://127.0.0.1:9999/pages?id=123":
			dial tcp 127.0.0.1:9999: connect: connection refused
	*/
}

func namedErrorDemo() {
	show("namedErrorDemo: program started ...")

	var (
		client      = http.Client{Timeout: time.Duration(time.Millisecond)}
		ErrResource = errors.New("resource error")
	)
	getRemoteResource := func() error {
		url := "http://127.0.0.1:9999/pages?id=123"
		_, err := client.Get(url)
		if err != nil {
			return ErrResource // error as a constant
		}
		return nil
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := getRemoteResource()
		if err != nil {
			fmt.Printf("error happend: %+v\n", err)
			switch err { // if you have limited number of error cases, you may use constants
			case ErrResource:
				http.Error(w, "remote resource error", 500)
			default:
				http.Error(w, "internal error", 500)
			}
			return
		}
		w.Write([]byte("all is OK"))
	}

	http.HandleFunc("/", handler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
	/*
		2024-04-17T12:32:17.965Z: namedErrorDemo: program started ...
		2024-04-17T12:32:17.965Z: Starting server at: string(127.0.0.1:8080);
		2024-04-17T12:32:17.965Z: Open url http://localhost:8080/
		error happend: resource error
	*/
}

func basicErrorDemo() {
	show("basicErrorDemo: program started ...")

	client := http.Client{Timeout: time.Duration(time.Millisecond)}

	getRemoteResource := func() error {
		url := "http://127.0.0.1:9999/pages?id=123"
		_, err := client.Get(url)
		if err != nil {
			// можно так:
			// return err // вернётся `timed out`. и что?
			// будет `res error: time out`. а где?

			// можно так:
			// return fmt.Errof("getRemoteResource: %+v", err)

			// но лучше так:
			// return fmt.Errorf("getRemoteResource: %s at %s", err, url)
			return fmt.Errorf("getRemoteResource: %+v at %s", err, url) // same shit
		}
		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		err := getRemoteResource()
		if err != nil {
			fmt.Printf("error happend: %+v\n", err) // logging
			http.Error(w, "internal error", 500)
			return
		}

		w.Write([]byte("no problem"))
	}

	http.HandleFunc("/", handler)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
	/*
		2024-04-17T12:16:04.604Z: basicErrorDemo: program started ...
		2024-04-17T12:16:04.604Z: Starting server at: string(127.0.0.1:8080);
		2024-04-17T12:16:04.604Z: Open url http://localhost:8080/
		error happend: getRemoteResource: Get "http://127.0.0.1:9999/pages?id=123":
			dial tcp 127.0.0.1:9999: connect: connection refused (Client.Timeout exceeded while awaiting headers)
			at http://127.0.0.1:9999/pages?id=123
	*/
}

func contextValueDemo() {
	show("contextValueDemo: program started ...")

	// сколько в среднем спим при эмуляции работы
	const AvgSleep = 50
	// линтер ругается если используем базовые типы в Value контекста
	// типа так безопаснее разграничивать
	type key int
	const timingsKey key = 1
	type Timing struct {
		Count    int
		Duration time.Duration
	}
	type ctxTimings struct {
		sync.Mutex
		Data map[string]*Timing
	}

	logContextTimings := func(ctx context.Context, path string, start time.Time) {
		// получаем тайминги из контекста
		// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
		timings, ok := ctx.Value(timingsKey).(*ctxTimings)
		if !ok {
			return
		}
		totalReal := time.Since(start)
		buf := bytes.NewBufferString(path)
		var totalRegistered time.Duration
		for timing, value := range timings.Data {
			totalRegistered += value.Duration
			buf.WriteString(fmt.Sprintf("\n\t%s(%d): %s", timing, value.Count, value.Duration))
		}
		buf.WriteString(fmt.Sprintf("\n\ttotal: %s", totalReal))
		buf.WriteString(fmt.Sprintf("\n\ttracked: %s", totalRegistered))
		buf.WriteString(fmt.Sprintf("\n\tunkn: %s", totalReal-totalRegistered))

		fmt.Println(buf.String())
	}

	// track time on exit, call in worker: `defer trackContextTimings(ctx, workName, time.Now())`
	trackContextTimings := func(ctx context.Context, recordName string, start time.Time) {
		// получаем тайминги из контекста
		// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
		timings, ok := ctx.Value(timingsKey).(*ctxTimings)
		if !ok { // wrong context
			return
		}

		elapsed := time.Since(start)
		// лочимся на случай конкурентной записи в мапку
		timings.Lock()
		defer timings.Unlock()

		// update context
		// если меткри ещё нет - мы её создадим, если есть - допишем в существующую
		if metric, metricExist := timings.Data[recordName]; !metricExist {
			// create new record
			timings.Data[recordName] = &Timing{
				Count:    1,
				Duration: elapsed,
			}
		} else {
			metric.Count++
			metric.Duration += elapsed
		}
	}

	emulateWork := func(ctx context.Context, workName string) {
		defer trackContextTimings(ctx, workName, time.Now()) // track on exit

		rnd := time.Duration(rand.Intn(AvgSleep))
		time.Sleep(time.Millisecond * rnd)
	}

	// pages handler for server
	loadPostsHandle := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		emulateWork(ctx, "checkCache")
		emulateWork(ctx, "loadPosts")
		emulateWork(ctx, "loadPosts")
		emulateWork(ctx, "loadPosts")
		time.Sleep(10 * time.Millisecond) // untracked
		emulateWork(ctx, "loadSidebar")
		emulateWork(ctx, "loadComments")

		fmt.Fprintln(w, "Request done")
	}

	// middleware
	timingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context() // create structure
			ctx = context.WithValue(ctx,
				timingsKey,
				&ctxTimings{
					Data: make(map[string]*Timing),
				})
			defer logContextTimings(ctx, r.URL.Path, time.Now()) // do calculations on exit

			next.ServeHTTP(w, r.WithContext(ctx)) // call internal processing
		})
	}

	// setup server

	// rand.Seed(time.Now().UTC().UnixNano()) // deprecated
	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)
	siteWithTiming := timingMiddleware(siteMux)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, siteWithTiming)
	show("end of program. ", err)
	/*
		2024-04-17T11:22:31.524Z: contextValueDemo: program started ...
		2024-04-17T11:22:31.524Z: Starting server at: string(127.0.0.1:8080);
		2024-04-17T11:22:31.524Z: Open url http://localhost:8080/
		/admin/panic
		        checkCache(1): 36.727597ms
		        loadPosts(3): 76.790603ms
		        loadSidebar(1): 1.437164ms
		        loadComments(1): 31.489476ms
		        total: 157.023282ms
		        tracked: 146.44484ms
		        unkn: 10.578442ms
		/favicon.ico
		        checkCache(1): 44.547705ms
		        loadPosts(3): 62.60222ms
		        loadSidebar(1): 1.22976ms
		        loadComments(1): 19.259073ms
		        total: 138.079069ms
		        tracked: 127.638758ms
		        unkn: 10.440311ms
	*/
}

func middlwareDemo() {
	show("middlwareDemo: program started ...")

	// naive web-page handler
	var pageWithAllChecks = func(w http.ResponseWriter, r *http.Request) {
		// recover from panic
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()

		// logging
		defer func(start time.Time) {
			fmt.Printf("[%s] %s, %s %s\n",
				r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
		}(time.Now())

		// auth
		_, err := r.Cookie("session_id")
		// учебный пример! это не проверка авторизации!
		if err != nil {
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// your logic
		var businessLogic = func() {
			show("foo")
			panic("oops")
		}
		businessLogic()
	}
	show(`демонстрация обработчика запроса, где шаги пайплайна реализуются самостоятельно, без мидлвари`, pageWithAllChecks)

	// server with middleware (MW) pattern in use

	// pages

	// /admin/
	adminIndex := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<a href="/">site index</a>`)
		fmt.Fprintln(w, "Admin main page")
	}
	// /admin/panic
	panicPage := func(w http.ResponseWriter, r *http.Request) {
		panic("this must be recovered")
	}
	// /login
	loginPage := func(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().Add(10 * time.Hour)
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   "foo",
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	// /logout
	logoutPage := func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	// /
	mainPage := func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		// учебный пример! это не проверка авторизации!
		loggedIn := (err != http.ErrNoCookie)
		if loggedIn {
			fmt.Fprintln(w, `<a href="/logout">logout</a>`)
			fmt.Fprintln(w, "Welcome, "+session.Value)
		} else {
			fmt.Fprintln(w, `<a href="/login">login</a>`)
			fmt.Fprintln(w, "You need to login")
		}
	}

	// middleware

	// auth MW, one layer
	adminAuthMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("adminAuthMiddleware", r.URL.Path)
			_, err := r.Cookie("session_id")
			// учебный пример! это не проверка авторизации!
			if err != nil {
				fmt.Println("no auth at", r.URL.Path)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r) // call internal handler
		})
	}
	// log MW, another layer
	accessLogMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("accessLogMiddleware", r.URL.Path)
			start := time.Now()

			next.ServeHTTP(w, r) // call internal handler

			fmt.Printf("[%s] %s, %s %s\n",
				r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
		})
	}

	panicMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("panicMiddleware", r.URL.Path)
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("recovered", err)
					http.Error(w, "Internal server error", 500)
				}
			}()

			next.ServeHTTP(w, r) // call internal handler
		})
	}

	// server setup

	adminPages := http.NewServeMux()
	adminPages.HandleFunc("/admin/", adminIndex)
	adminPages.HandleFunc("/admin/panic", panicPage)

	// set middleware: auth only for admin
	adminPagesWithAuth := adminAuthMiddleware(adminPages)

	sitePages := http.NewServeMux()
	sitePages.Handle("/admin/", adminPagesWithAuth)
	sitePages.HandleFunc("/login", loginPage)
	sitePages.HandleFunc("/logout", logoutPage)
	sitePages.HandleFunc("/", mainPage)

	// set middleware
	sitePagesWithLog := accessLogMiddleware(sitePages)
	sitePagesWithLogAndRecover := panicMiddleware(sitePagesWithLog)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, sitePagesWithLogAndRecover) // 1) recover, 2) log, 3) auth-if-admin
	show("end of program. ", err)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
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
