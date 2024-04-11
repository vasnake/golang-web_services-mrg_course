package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"
)

func main() {
	net_listen()
	// httpDemo()
	// pagesDemo()
	// servehttpDemo()
	// muxDemo()
	// serversDemo()
	// get_paramsDemo()
	// postFormDemo()
	// cookiesDemo()
	// headersDemo()
	// staticServeDemo()
	// file_uploadDemo()
	// requestDemo()
	// inlineTemplate()
	// fileTemplate()
	// methodCallFromTemplate()
	// funcCallFromTemplate()
}

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

func postFormDemo() {
	show("program started ...")

	var loginFormTmpl = []byte(`
	<html> <body> <form action="/" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login please ...">
	</form> </body> </html>
	`)

	var loginFormHandler = func(w http.ResponseWriter, r *http.Request) {
		// Show form only when requested with GET (not POST)
		if r.Method != http.MethodPost {
			w.Write(loginFormTmpl)
			return
		}

		// if requested with POST

		// parse form data explicitly
		// r.ParseForm()
		// loginFormValue := r.Form["login"][0]

		// or implicitly
		loginFormValue := r.FormValue("login")

		fmt.Fprintln(w, "User login: ", loginFormValue)
	}

	http.HandleFunc("/", loginFormHandler)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show("Open url http://localhost:8080/")
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func cookiesDemo() {
	show("program started ...")

	var mainPage = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		sessionCookieRef, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "You need to login: ")
			fmt.Fprintln(w, `<a href="/login">login</a>`)
		} else {
			fmt.Fprintln(w, "Welcome, `"+sessionCookieRef.Value+"`. Now logout: ")
			fmt.Fprintln(w, `<a href="/logout">logout</a>`)
		}
	}

	var loginPage = func(w http.ResponseWriter, r *http.Request) {
		// Imagine that you check user credentials already, now set session cookie:
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   "Foo Bar",
			Expires: time.Now().Add(10 * time.Minute),
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	var logoutPage = func(w http.ResponseWriter, r *http.Request) {
		sessionCookieRef, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			sessionCookieRef.Expires = time.Now().AddDate(0, 0, -1) // expire yesterday
			http.SetCookie(w, sessionCookieRef)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", mainPage)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func headersDemo() {
	show("program started ...")

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("RequestID", "d41d8cd98f00b204")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "You browser UA is:", r.UserAgent())
		fmt.Fprintln(w, "You `Accept`:", r.Header.Get("Accept"))
	}

	http.HandleFunc("/", handler)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func staticServeDemo() {
	show("program started ...")

	var rootHandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			Hello World! <br />
			<img src="/data/img/gopher.png" />
		`))
	}

	// TLDR: `/data/img/gopher.png` => `./week04/static/img/gopher.png`
	// when looking for `/data/img/gopher.png`:
	// - strip `/data/`
	// - read `img/gopher.png` in `./week04/static` directory
	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./week04/static")),
	)

	http.HandleFunc("/", rootHandlerFunc)
	http.Handle("/data/", staticHandler)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func file_uploadDemo() {
	show("program started ...")

	var uploadFormTmpl = []byte(`
	<html> <body> <form action="/upload" method="post" enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload selected file ...">
	</form> </body> </html>
	`)

	var mainPage = func(w http.ResponseWriter, r *http.Request) {
		w.Write(uploadFormTmpl)
	}

	var uploadFormHandler = func(w http.ResponseWriter, r *http.Request) {
		// method: post

		// process first 5 MB
		if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
			fmt.Println(err)
			return
		}

		file, fileHeaderRef, err := r.FormFile("my_file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Fprintf(w, "Filename: %v\n", fileHeaderRef.Filename)
		fmt.Fprintf(w, "Header: %#v\n", fileHeaderRef.Header)

		// some file processing
		hash := md5.New()
		io.Copy(hash, file)
		fmt.Fprintf(w, "md5: %x\n", hash.Sum(nil))
	}

	var uploadRawContentHandler = func(w http.ResponseWriter, r *http.Request) {
		// curl -v -X POST -H "Content-Type: application/json" -d '{"id": 42, "user": "Foo Bar"}' http://localhost:8080/raw_body

		bodyBytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		// some bytes processing
		type Params struct {
			ID   int
			User string
		}
		decodedContent := &Params{}
		if err = json.Unmarshal(bodyBytes, decodedContent); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, "content-type %#v\n", r.Header.Get("Content-Type"))
		fmt.Fprintf(w, "params %#v\n", decodedContent)
	}

	// render upload form
	http.HandleFunc("/", mainPage)
	// upload file using form
	http.HandleFunc("/upload", uploadFormHandler)
	// upload raw data using POST
	http.HandleFunc("/raw_body", uploadRawContentHandler)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func requestDemo() {
	show("program started ...")

	// async, run the server with 2 endpoints
	var runHttpServer = func() {
		// register server root handler
		http.HandleFunc(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "root handler, incoming request r: %#v\n", r)
				fmt.Fprintf(w, "root handler, r.Url: %#v\n", r.URL)
			},
		)

		// register server `raw_body` handler
		http.HandleFunc(
			"/raw_body",
			func(w http.ResponseWriter, r *http.Request) {
				bodyBytes, err := io.ReadAll(r.Body)
				defer r.Body.Close() // baware of leaks
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				fmt.Fprintf(w, "`raw_body` handler, raw body %s\n", string(bodyBytes))
			},
		)

		// run
		const addrStr = ":8080"
		show("Starting server at: ", addrStr)
		show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
		err := http.ListenAndServe(addrStr, nil)
		show("end of server. ", err)
	}

	var execGetQuery = func(uri string) {
		if len(uri) < 1 {
			uri = "http://127.0.0.1:8080/?param1=123&param2=test"
		}

		resp, err := http.Get(uri)
		if err != nil {
			show("While doing http.Get, got error: ", err)
			return
		}
		defer resp.Body.Close() // beware of leaks

		respBody, err := io.ReadAll(resp.Body)
		show("http.Get, response.Body: ", string(respBody), err)
	}

	var execLowLevelGetQuery = func(uri string) {
		if len(uri) < 1 {
			uri = "http://127.0.0.1:8080/?id=42"
		}

		// prepare request
		req := &http.Request{
			Method: http.MethodGet,
			Header: http.Header{
				"User-Agent": {"coursera/golang"},
			},
		}
		req.URL, _ = url.Parse(uri) // id=42
		req.URL.Query().Set("user", "Foo")

		// ask service
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			show("While doing http.DefaultClient.Do, got error: ", err)
			return
		}
		defer resp.Body.Close() // beware of leaks

		respBody, err := io.ReadAll(resp.Body)
		show("Low-level GET (http.DefaultClient.Do), responce.Body: ", string(respBody), err)
	}

	var execLowLowLevelPostQuery = func(uri string) {
		if len(uri) < 1 {
			uri = "http://127.0.0.1:8080/raw_body"
		}

		// define request using low-level API

		transportRef := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		clientRef := &http.Client{
			Timeout:   time.Second * 10,
			Transport: transportRef,
		}

		bodyData := `{"id": 24, "user": "Bar"}`
		bodyBytesRef := bytes.NewBufferString(bodyData)

		reqRef, _ := http.NewRequest(http.MethodPost, uri, bodyBytesRef)
		reqRef.Header.Add("Content-Type", "application/json")
		reqRef.Header.Add("Content-Length", strconv.Itoa(len(bodyData)))

		// ask service

		respRef, err := clientRef.Do(reqRef)
		if err != nil {
			show("While doing http.Client.Do, got error: ", err)
			return
		}
		defer respRef.Body.Close() // beware of leaks

		respBodyBytes, err := io.ReadAll(respRef.Body)
		show("LowLow-level POST (http.Client.Do), responce.Body: ", string(respBodyBytes), err)
	}

	go runHttpServer() // async service
	time.Sleep(100 * time.Millisecond)

	execGetQuery("http://127.0.0.1:8080/?param1=123&param2=test") // ask service using get
	execLowLevelGetQuery("http://127.0.0.1:8080/?id=42")          // ask service with low-level get API
	execLowLowLevelPostQuery("http://127.0.0.1:8080/raw_body")    // ask service using post and very low-level API

	show("end of program.")
	/*
	   2023-12-11T14:12:28.795Z: program started ...
	   2023-12-11T14:12:28.796Z: Starting server at: string(:8080);
	   2023-12-11T14:12:28.796Z: Open url http://localhost:8080/

	   2023-12-11T14:12:28.897Z: http.Get, response.Body: string(root handler, incoming request r:
	   	&http.Request{Method:"GET", URL:(*url.URL)(0xc0001521b0), Proto:"HTTP/1.1", ProtoMajor:1, ProtoMinor:1,
	   	Header:http.Header{"Accept-Encoding":[]string{"gzip"}, "User-Agent":[]string{"Go-http-client/1.1"}},
	   	Body:http.noBody{}, GetBody:(func() (io.ReadCloser, error))(nil), ContentLength:0,
	   	TransferEncoding:[]string(nil), Close:false, Host:"127.0.0.1:8080",
	   	Form:url.Values(nil), PostForm:url.Values(nil), MultipartForm:(*multipart.Form)(nil), Trailer:http.Header(nil),
	   	RemoteAddr:"127.0.0.1:47746",
	   	RequestURI:"/?param1=123&param2=test", TLS:(*tls.ConnectionState)(nil), Cancel:(<-chan struct {})(nil),
	   	Response:(*http.Response)(nil), ctx:(*context.cancelCtx)(0xc00007e230)}
	   root handler, r.Url: &url.URL{Scheme:"", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/", RawPath:"",
	   OmitHost:false, ForceQuery:false,
	   RawQuery:"param1=123&param2=test", Fragment:"", RawFragment:""});
	   <nil>(<nil>);

	   2023-12-11T14:12:28.897Z: Low-level GET (http.DefaultClient.Do), responce.Body: string(root handler, incoming request r:
	   	&http.Request{Method:"GET", URL:(*url.URL)(0xc0001801b0), Proto:"HTTP/1.1", ProtoMajor:1, ProtoMinor:1,
	   	Header:http.Header{"Accept-Encoding":[]string{"gzip"}, "User-Agent":[]string{"coursera/golang"}},
	   	Body:http.noBody{}, GetBody:(func() (io.ReadCloser, error))(nil), ContentLength:0,
	   	TransferEncoding:[]string(nil), Close:false, Host:"127.0.0.1:8080", Form:url.Values(nil), PostForm:url.Values(nil),
	   	MultipartForm:(*multipart.Form)(nil), Trailer:http.Header(nil), RemoteAddr:"127.0.0.1:47746",
	   	RequestURI:"/?id=42", TLS:(*tls.ConnectionState)(nil), Cancel:(<-chan struct {})(nil),
	   	Response:(*http.Response)(nil), ctx:(*context.cancelCtx)(0xc0001a2050)}
	   root handler, r.Url: &url.URL{Scheme:"", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/", RawPath:"", OmitHost:false, ForceQuery:false,
	   RawQuery:"id=42", Fragment:"", RawFragment:""});
	   <nil>(<nil>);

	   2023-12-11T14:12:28.898Z: LowLow-level POST (http.Client.Do), responce.Body: string(`raw_body` handler, raw body {"id": 24, "user": "Bar"});
	   <nil>(<nil>);

	   2023-12-11T14:12:28.898Z: end of program.
	*/
}

func inlineTemplate() {
	show("program started ...")

	type templateParamsStruct struct {
		URL     string
		Browser string
	}

	const EXAMPLE_TEMPLATE = `
	Browser {{.Browser}}
	
	you at {{.URL}}
	`

	var rootPageHandler = func(w http.ResponseWriter, r *http.Request) {
		templateParams := templateParamsStruct{
			URL:     r.URL.String(),
			Browser: r.UserAgent(),
		}

		templateRef := template.New(`example`)
		templateRef, _ = templateRef.Parse(EXAMPLE_TEMPLATE)
		templateRef.Execute(w, templateParams)
	}

	http.HandleFunc("/", rootPageHandler)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func fileTemplate() {
	show("program started ...")

	type User struct {
		ID     int
		Name   string
		Active bool
	}

	users := []User{
		{1, "Foo", true},
		{2, "<i>Bar</i>", false},
		{3, "Baz", true},
	}

	templateRef := template.Must(template.ParseFiles("week04/static/users.html"))
	/*
		   <html>
		   <body>
		   	<h1>Users</h1>
		   	{{range .Users}}
		   		{{.ID}}
				<b>{{.Name}}</b>
		   		{{if .Active}}active{{end}}
		   		<br />
		   	{{end}}
		   </body>
		   </html>
	*/

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			templateRef.Execute(w,
				struct { // type
					Users []User
				}{ // data
					users,
				})
		},
	)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err := http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func methodCallFromTemplate() {
	show("program started ...")

	users := []User_Template{
		{1, "Foo", true},
		{2, "Bar", false},
		{3, "Baz", true},
	}

	templateRef, err := template.New("method.html").ParseFiles("week04/static/method.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			if err := templateRef.ExecuteTemplate(
				w,
				"method.html", // name of the file, always
				struct {
					Users []User_Template
				}{
					users,
				},
			); err != nil {
				panic(err)
			}
		},
	)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err = http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
}

func funcCallFromTemplate() {
	show("program started ...")

	type User struct {
		ID     int
		Name   string
		Active bool
	}

	var IsUserOdd = func(user *User) bool {
		// function to call from template
		return (user.ID % 2) != 0
	}

	// register functions
	templateFuncs := template.FuncMap{
		"OddUser": IsUserOdd,
	}

	// add funcs before parsing
	templateRef, err := template.New("func.html").Funcs(templateFuncs).ParseFiles("week04/static/func.html")
	if err != nil {
		panic(err)
	}

	users := []User{
		{1, "Foo", true},
		{2, "Bar", false},
		{3, "Baz", true},
	}

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			if err := templateRef.ExecuteTemplate(w, "func.html",
				struct {
					Users []User
				}{
					users,
				}); err != nil {
				panic(err)
			}
		},
	)

	const addrStr = ":8080"
	show("Starting server at: ", addrStr)
	show(fmt.Sprintf("Open url http://localhost%s/", addrStr))
	err = http.ListenAndServe(addrStr, nil)
	show("end of program. ", err)
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
