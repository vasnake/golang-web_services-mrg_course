package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

func net_listen() {
	show("program started ...")

	var handleConnection = func(conn net.Conn) {
		defer conn.Close() // on exit

		name := conn.RemoteAddr().String()
		show("Connected: ", name)
		conn.Write([]byte("Hello, " + name + ". Type `Exit` to quit from your session.\n\r"))

		currLine := bufio.NewScanner(conn)
		for currLine.Scan() {
			text := currLine.Text()
			// ignore empty lines
			if text == "Exit" {
				conn.Write([]byte("Bye\n\r"))
				show("Disconnected: ", name)
				break
			} else if text != "" {
				show("Got y from x, (x, y): ", name, text)
				conn.Write([]byte("You typed: `" + text + "`\n\r"))
			}
		} // end loop
	}

	listner, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	for {
		show("Waiting for connection ... ", listner)
		conn, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		show("Got new connection: ", conn)

		// async session
		go handleConnection(conn)
	}
	// show("end of program.")
	/*
	   2023-12-07T14:31:09.205Z: program started ...
	   2023-12-07T14:31:09.206Z: Waiting for connection ... *net.TCPListener(&{0xc000126000 {<nil> 0 0}});
	   ^Z
	   sandbox$ bg
	   sandbox$ netcat 127.0.0.1 8080
	   2023-12-07T14:34:39.541Z: Got new connection: *net.TCPConn(&{{0xc000126080}});
	   2023-12-07T14:34:39.541Z: Waiting for connection ... *net.TCPListener(&{0xc000126000 {<nil> 0 0}});
	   2023-12-07T14:34:39.541Z: Connected: string(127.0.0.1:59492);
	   Hello, 127.0.0.1:59492. Type `Exit` to quit from your session.
	   ff
	   2023-12-07T14:35:02.131Z: Got y from x, (x, y): string(127.0.0.1:59492); string(ff);
	   You typed: `ff`
	   Exit
	   Bye
	   2023-12-07T14:35:20.403Z: Disconnected: string(127.0.0.1:59492);
	   sandbox$ netcat 127.0.0.1 8080
	   2023-12-07T14:35:50.339Z: Got new connection: *net.TCPConn(&{{0xc000126100}});
	   2023-12-07T14:35:50.339Z: Waiting for connection ... *net.TCPListener(&{0xc000126000 {<nil> 0 0}});
	   2023-12-07T14:35:50.339Z: Connected: string(127.0.0.1:33696);
	   Hello, 127.0.0.1:33696. Type `Exit` to quit from your session.
	   qq
	   2023-12-07T14:36:07.322Z: Got y from x, (x, y): string(127.0.0.1:33696); string(qq);
	   You typed: `qq`
	   Exit
	   2023-12-07T14:36:11.794Z: Disconnected: string(127.0.0.1:33696);
	   Bye
	*/
}

func httpDemo() {
	show("program started ...")

	const addrStr = ":8080"
	const servingUrlPattern = "/"

	var handler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Привет, мир!")
		w.Write([]byte("!!!\n"))
	}

	// register handler in server mux
	http.HandleFunc(servingUrlPattern, handler)

	show("Starting server at: ", addrStr)
	err := http.ListenAndServe(addrStr, nil) // handler=nil => using mux

	show("end of program. ", err)
	/*
	   2023-12-07T14:56:19.644Z: program started ...
	   2023-12-07T14:56:19.644Z: Starting server at: string(:8080);
	   ^Z
	   sandbox$  bg
	   sandbox$ curl 127.0.0.1:8080
	   Привет, мир!
	   !!!
	*/
}

func pagesDemo() {
	show("program started ...")

	var mainPageHandler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Main page")
	}

	// register handlers in mux

	http.HandleFunc("/", mainPageHandler)

	// N.B. route w.o. ending slash
	http.HandleFunc(
		"/page",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Single page:", r.URL.String())
		},
	)

	// N.B. route with ending slash, it works as a prefix for a class of supported routes
	http.HandleFunc(
		"/pages/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Multiple pages:", r.URL.String())
		},
	)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	err := http.ListenAndServe(addrStr, nil)

	show("end of program. ", err)
	/*
	   2023-12-07T15:08:06.651Z: program started ...
	   2023-12-07T15:08:06.651Z: Starting server at: string(:8080);
	   ^Z
	   sandbox$  bg
	   sandbox$ curl 127.0.0.1:8080
	   Main page
	   sandbox$ curl 127.0.0.1:8080/page
	   Single page: /page
	   sandbox$ curl 127.0.0.1:8080/pages/
	   Multiple pages: /pages/
	   sandbox$ curl 127.0.0.1:8080/pages/foo
	   Multiple pages: /pages/foo
	   sandbox$ curl 127.0.0.1:8080/pages/bar/foo
	   Multiple pages: /pages/bar/foo
	*/
}

func servehttpDemo() {
	show("program started ...")

	// All requests served by one instance of a struct.
	// But, different routes could use different instances of a struct.
	testHandlerRef := &Handler_servehttp{Name: "test"}
	http.Handle("/test/", testHandlerRef)

	rootHandlerRef := &Handler_servehttp{Name: "root"}
	http.Handle("/", rootHandlerRef)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
	/*
	   2023-12-07T15:28:49.715Z: program started ...
	   2023-12-07T15:28:49.715Z: Starting server at: string(:8080);
	   ^Z
	   sandbox$  bg
	   sandbox$ curl 127.0.0.1:8080/
	   Name: root URL: /
	   sandbox$ curl 127.0.0.1:8080/test/
	   Name: test URL: /test/
	*/
}

func muxDemo() {
	show("program started ...")

	var handler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Processed URL:", r.URL.String())
	}

	// multiplexor w/o server
	muxRef := http.NewServeMux()
	// register handler for path
	muxRef.HandleFunc("/", handler) // vs http.Handle("/", rootHandlerRef)

	// server+mux defined with optional parameters
	const addrStr = ":8080"
	server := http.Server{
		Addr:         addrStr,
		Handler:      muxRef,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	show("Starting server at: ", addrStr)
	err := server.ListenAndServe() // vs err := http.ListenAndServe(addrStr, nil)
	show("end of program.", err)
	/*
	   2023-12-07T15:41:31.726Z: program started ...
	   2023-12-07T15:41:31.726Z: Starting server at: string(:8080);
	   ...
	   sandbox$ curl 127.0.0.1:8080/test/
	   Processed URL: /test/
	*/
}

func serversDemo() {
	show("program started ...")

	// different addr, same behaviour (way to scale server?)
	var runServer = func(addr string) {
		muxRef := http.NewServeMux()
		muxRef.HandleFunc(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Root handler, processed request: Addr:", addr, "URL:", r.URL.String())
			},
		)

		server := http.Server{
			Addr:    addr,
			Handler: muxRef,
		}

		show("Starting server at: ", addr)
		err := server.ListenAndServe()
		show("Server stopped: ", err)
	}

	// two async servers on two different ports
	var addrStr1, addrStr2 = ":8081", ":8080"
	// async
	go runServer(addrStr1)
	// main goroutine
	runServer(addrStr2)

	show("end of program.")
	/*
	   2023-12-07T15:51:08.757Z: program started ...
	   2023-12-07T15:51:08.757Z: Starting server at: string(:8080);
	   2023-12-07T15:51:08.757Z: Starting server at: string(:8081);
	   ...
	   sandbox$ curl 127.0.0.1:8080/test/
	   Root handler, processed request: Addr: :8080 URL: /test/
	   sandbox$ curl 127.0.0.1:8081/test/
	   Root handler, processed request: Addr: :8081 URL: /test/
	*/
}

func get_paramsDemo() {
	show("program started ...")

	var handler = func(w http.ResponseWriter, r *http.Request) {
		// from url string
		getParamValue := r.URL.Query().Get("param")
		if getParamValue != "" {
			fmt.Fprintln(w, "`param`:", getParamValue)
		}

		// from all params storages
		queryKeyValue := r.FormValue("key")
		if queryKeyValue != "" {
			fmt.Fprintln(w, "`key`:", queryKeyValue)
		}
	}

	http.HandleFunc("/", handler)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
	/*
	   2023-12-07T16:13:33.721Z: program started ...
	   2023-12-07T16:13:33.721Z: Starting server at: string(:8080);
	   ...
	   sandbox$ curl -X GET "http://127.0.0.1:8080/test/foo?key=bar&param=baz"
	   `param`: baz
	   `key`: bar
	   sandbox$ curl -X PUT -d param=putParam -d key=putKey "http://127.0.0.1:8080/test/foo?key=getKey&param=getParam"
	   `param`: getParam
	   `key`: putKey
	   sandbox$ curl -X POST -d param=postParam -d key=postKey "http://127.0.0.1:8080/test/foo?key=getKey&param=getParam"
	   `param`: getParam
	   `key`: postKey
	*/
}

func main() {
	// net_listen()
	// httpDemo()
	// pagesDemo()
	// servehttpDemo()
	// muxDemo()
	// serversDemo()
	get_paramsDemo()
}

func demoTemplate() {
	show("program started ...")
	show("end of program.")
}

func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}

// ts return current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}
