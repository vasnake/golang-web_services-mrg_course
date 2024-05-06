package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

func main() {
	// sessionServiceBefore()
	// sessionServiceAfter()

	// netRpcClientServerSessions()
	netRpcJsonServerSessions()
}

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func netRpcJsonServerSessions() {
	show("netRpcJsonServerSessions: SERVER started ...")
	/*
	   {
	      "jsonrpc":"2.0",
	      "id":1,
	      "method":"SessionManager.Create",
	      "params":[
	         {
	            "login":"rvasily",
	            "useragent":"chrome"
	         }
	      ]
	   }

	   curl -v -X POST -H "Content-Type: application/json" -H "X-Auth: 123" -d '{"jsonrpc":"2.0", "id": 1, "method": "SessionManager_Server.Create", "params": [{"login":"rvasily", "useragent": "chrome"}]}' http://localhost:8081/rpc

	   curl -v -X POST -H "Content-Type: application/json" -H "X-Auth: 123" -d '{"jsonrpc":"2.0", "id": 2, "method": "SessionManager_Server.Check", "params": [{"id":"XVlBzgbaiC"}]}' http://localhost:8081/rpc

	*/
	_ = `
	Note: Unnecessary use of -X or --request, POST is already inferred.
	*   Trying 127.0.0.1:8081...
	* Connected to localhost (127.0.0.1) port 8081 (#0)
	> POST /rpc HTTP/1.1
	> Host: localhost:8081
	> User-Agent: curl/7.81.0
	> Accept: */*
	> Content-Type: application/json
	> X-Auth: 123
	> Content-Length: 124
	>
	rpc auth:  123
	call Create &{rvasily chrome}
	2024/05/06 15:50:44 http: superfluous response.WriteHeader call from main.(*SessionManagerHttpRpcHandler).ServeHTTP (main.go:95)
	* Mark bundle as not supporting multiuse
	< HTTP/1.1 200 OK
	< Content-Type: application/json
	< Date: Mon, 06 May 2024 12:50:44 GMT
	< Content-Length: 51
	<
	{"id":1,"result":{"ID":"hXGFMaeTAE"},"error":null}
	* Connection #0 to host localhost left intact
	`

	sessManager := NewSessManager_Server()
	server := rpc.NewServer()
	server.Register(sessManager)

	sessionHandler := &SessionManagerHttpRpcHandler{
		rpcServer: server,
	}
	http.Handle("/rpc", sessionHandler)

	fmt.Println("starting server at :8081")
	err := http.ListenAndServe(":8081", nil)

	show("end of program. ", err)
}

type SessionManagerHttpRpcHandler struct {
	rpcServer *rpc.Server
}

func (h *SessionManagerHttpRpcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rpc auth: ", r.Header.Get("X-Auth"))

	serverCodec := jsonrpc.NewServerCodec(&ReadWriteSkipClose{
		in:  r.Body,
		out: w,
	})

	w.Header().Set("Content-type", "application/json")
	err := h.rpcServer.ServeRequest(serverCodec)
	if err != nil {
		log.Printf("Error while serving JSON request: %v", err)
		http.Error(w, `{"error":"cant serve request"}`, 500)
	} else {
		w.WriteHeader(200)
	}
}

type ReadWriteSkipClose struct {
	in  io.Reader
	out io.Writer
}

func (c *ReadWriteSkipClose) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *ReadWriteSkipClose) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *ReadWriteSkipClose) Close() error                      { return nil }

func netRpcClientServerSessions() {
	show("netRpcClientServerSessions: program started ...")

	var server = func() {
		sessManager := NewSessManager_Server()
		rpc.Register(sessManager)
		rpc.HandleHTTP()

		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal("listen error:", err)
		}

		fmt.Println("starting server at :8081")
		http.Serve(listener, nil)
	}

	var client = func() {
		var sessManager SessionManagerI_Client

		sessManager = NewSessManager()

		// создаем сессию
		sessId, err := sessManager.Create(
			&Session{
				Login:     "baz",
				Useragent: "chrome",
			})
		fmt.Println("sessId", sessId, err)

		// проеряем сессию
		sess := sessManager.Check(
			&SessionID{
				ID: sessId.ID,
			})
		fmt.Println("sess", sess)

		// удаляем сессию
		sessManager.Delete(
			&SessionID{
				ID: sessId.ID,
			})

		// проверяем еще раз
		sess = sessManager.Check(
			&SessionID{
				ID: sessId.ID,
			})
		fmt.Println("sess", sess)
	}

	go server()
	time.Sleep(987 * time.Millisecond)

	client()
	show("end of program")
}

type SessionID struct {
	ID string
}
type Session struct {
	Login     string
	Useragent string
}

type SessionManagerI_Client interface {
	Create(*Session) (*SessionID, error)
	Check(*SessionID) *Session
	Delete(*SessionID)
}

type SessionManager_Client struct {
	client *rpc.Client
}

func NewSessManager() *SessionManager_Client {
	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	return &SessionManager_Client{
		client: client,
	}
}
func (sm *SessionManager_Client) Create(in *Session) (*SessionID, error) {
	id := new(SessionID)
	err := sm.client.Call("SessionManager_Server.Create", in, id)
	if err != nil {
		fmt.Println("SessionManager_Server.Create error:", err)
		return nil, nil
	}
	return id, nil
}

func (sm *SessionManager_Client) Check(in *SessionID) *Session {
	sess := new(Session)
	err := sm.client.Call("SessionManager_Server.Check", in, sess)
	if err != nil {
		fmt.Println("SessionManager_Server.Check error:", err)
		return nil
	}
	return sess
}

func (sm *SessionManager_Client) Delete(in *SessionID) {
	var reply int
	err := sm.client.Call("SessionManager_Server.Delete", in, &reply)
	if err != nil {
		fmt.Println("SessionManager_Server.Delete error:", err)
	}
}

type SessionManager_Server struct {
	mu       sync.RWMutex
	sessions map[SessionID]*Session
}

func NewSessManager_Server() *SessionManager_Server {
	return &SessionManager_Server{
		mu:       sync.RWMutex{},
		sessions: map[SessionID]*Session{},
	}
}
func (sm *SessionManager_Server) Create(in *Session, out *SessionID) error {
	const sessKeyLen = 10

	fmt.Println("call Create", in)
	id := &SessionID{RandStringRunes(sessKeyLen)}
	sm.mu.Lock()
	sm.sessions[*id] = in
	sm.mu.Unlock()
	*out = *id
	return nil
}

func (sm *SessionManager_Server) Check(in *SessionID, out *Session) error {
	fmt.Println("call Check", in)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		*out = *sess
	}
	return nil
}

func (sm *SessionManager_Server) Delete(in *SessionID, out *int) error {
	fmt.Println("call Delete", in)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
	*out = 1
	return nil
}

func sessionServiceAfter() {
	show("sessionServiceAfter: program started ...")

	var sessManager SessionManagerI_demo2
	sessManager = NewSessManager_demo2()

	// создаем сессию
	sessId, err := sessManager.Create(
		&Session_after{
			Login:     "bar",
			Useragent: "chrome",
		})
	show("sessId: ", sessId, err)

	// проеряем сессию
	sess := sessManager.Check(
		&SessionID_after{
			ID: sessId.ID,
		})
	show("sess: ", sess)

	// удаляем сессию
	sessManager.Delete(
		&SessionID_after{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = sessManager.Check(
		&SessionID_after{
			ID: sessId.ID,
		})
	show("sess: ", sess)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	// err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

type SessionManagerI_demo2 interface {
	Create(*Session_after) (*SessionID_after, error)
	Check(*SessionID_after) *Session_after
	Delete(*SessionID_after)
}

type SessionID_after struct {
	ID string
}
type Session_after struct {
	Login     string
	Useragent string
}
type SessionManager_demo2 struct {
	mu       *sync.RWMutex
	sessions map[SessionID_after]*Session_after
}

func NewSessManager_demo2() *SessionManager_demo2 {
	return &SessionManager_demo2{
		mu:       &sync.RWMutex{},
		sessions: map[SessionID_after]*Session_after{},
	}
}
func (sm *SessionManager_demo2) Create(in *Session_after) (*SessionID_after, error) {
	const sessKeyLen = 10

	sm.mu.Lock()
	id := SessionID_after{RandStringRunes(sessKeyLen)}
	sm.mu.Unlock()
	sm.sessions[id] = in
	return &id, nil
}

func (sm *SessionManager_demo2) Check(in *SessionID_after) *Session_after {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		return sess
	}
	return nil
}

func (sm *SessionManager_demo2) Delete(in *SessionID_after) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
}

func sessionServiceBefore() {
	show("sessionServiceBefore: program started ...")

	// CLI

	// создаем сессию
	sessId, err := AuthCreateSession_before(
		&Session_before{
			Login:     "foo",
			Useragent: "chrome",
		})
	show("sessId: ", sessId, err)

	// проеряем сессию
	sess := AuthCheckSession_before(
		&SessionID_before{
			ID: sessId.ID,
		})
	show("sess: ", sess)

	// удаляем сессию
	AuthSessionDelete_before(
		&SessionID_before{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = AuthCheckSession_before(
		&SessionID_before{
			ID: sessId.ID,
		})
	show("sess: ", sess)

	// WEB

	http.HandleFunc("/", innerPage_demo1)
	http.HandleFunc("/login", loginPage_demo1)
	http.HandleFunc("/logout", logoutPage_demo1)

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func checkSession_demo1(r *http.Request) (*Session_before, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sess := AuthCheckSession_before(&SessionID_before{
		ID: cookieSessionID.Value,
	})
	return sess, nil
}

func innerPage_demo1(w http.ResponseWriter, r *http.Request) {
	var loginFormTmpl_demo1 = []byte(`
<html>
	<body>
	<form action="/login" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`)

	sess, err := checkSession_demo1(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if sess == nil {
		w.Write(loginFormTmpl_demo1)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Welcome, "+sess.Login+" <br />")
	fmt.Fprintln(w, "Session ua: "+sess.Useragent+" <br />")
	fmt.Fprintln(w, `<a href="/logout">logout</a>`)
}

func loginPage_demo1(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue("login")
	expiration := time.Now().Add(365 * 24 * time.Hour)

	sess, err := AuthCreateSession_before(&Session_before{
		Login:     inputLogin,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage_demo1(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	AuthSessionDelete_before(&SessionID_before{
		ID: session.Value,
	})

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

func AuthCheckSession_before(in *SessionID_before) *Session_before {
	mu_demo1.RLock()
	defer mu_demo1.RUnlock()
	if sess, ok := sessions_demo1[*in]; ok {
		return sess
	}
	return nil
}

func AuthCreateSession_before(in *Session_before) (*SessionID_before, error) {
	const sessKeyLen = 10

	mu_demo1.Lock()
	id := SessionID_before{RandStringRunes(sessKeyLen)}
	mu_demo1.Unlock()
	sessions_demo1[id] = in
	return &id, nil
}

func AuthSessionDelete_before(in *SessionID_before) {
	mu_demo1.Lock()
	defer mu_demo1.Unlock()
	delete(sessions_demo1, *in)
}

type SessionID_before struct {
	ID string
}

var (
	sessions_demo1 = map[SessionID_before]*Session_before{}
	mu_demo1       = &sync.RWMutex{}
)

type Session_before struct {
	Login     string
	Useragent string
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
