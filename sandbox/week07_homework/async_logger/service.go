package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"

	// sync "sync"
	"time"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type AclListType map[string][]string

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {
	// show("StartMyMicroservice, addr: ", listenAddr)

	// parse params
	var aclData = make(AclListType, 16)
	err := json.Unmarshal([]byte(ACLData), &aclData)
	if err != nil {
		show("StartMyMicroservice failed, invalid ACL json: ", err)
		return err
	}

	// setup server
	adminRef := NewAdminServerImpl().SetAuth(aclData)
	server := grpc.NewServer()
	RegisterBizServer(server, NewBizServerImpl(adminRef))
	RegisterAdminServer(server, adminRef)

	// setup transport
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		show("StartMyMicroservice failed, net.Listen error: ", err)
		return err
	}

	// start server
	show("StartMyMicroservice, starting gRPC server at ", listenAddr)
	go server.Serve(listener)

	// stop server
	go func() {
		var stopSignal = <-ctx.Done()
		show("StartMyMicroservice, stopping gRPC server at ", listenAddr, stopSignal)
		// subs.RemoveAll()
		server.GracefulStop()
	}()

	return nil
}

type BizServerImpl struct {
	// acl AclListType
	// mutex sync.RWMutex
	// data  map[string]string
	adminRef *AdminServerImpl
	UnimplementedBizServer
}

func NewBizServerImpl(admin *AdminServerImpl) *BizServerImpl {
	return &BizServerImpl{
		// acl: make(AclListType, 16),
		// mutex: sync.RWMutex{},
		// data:  make(map[string]string, 16),
		adminRef: admin,
	}
}

func (bs *BizServerImpl) Check(ctx context.Context, in *Nothing) (*Nothing, error) {
	ctx = context.WithValue(ctx, "method", "/main.Biz/Check") // TODO: get value from grpc.*ServerInfo.FullMethod

	// event
	bs.adminRef.pushEvent(ctx)

	return &Nothing{}, nil
}

func (bs *BizServerImpl) Add(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (bs *BizServerImpl) Test(ctx context.Context, _ *Nothing) (*Nothing, error) {
	ctx = context.WithValue(ctx, "method", "/main.Biz/Test") // TODO: get value from grpc.*ServerInfo.FullMethod

	// event
	bs.adminRef.pushEvent(ctx)

	// auth
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return nil, status.Errorf(codes.InvalidArgument, "context w/o metadata")
	}

	consumer := strings.Join(md.Get("consumer"), "")
	// show("biz test, consumer: ", consumer)

	if consumer != "biz_admin" {
		return nil, status.Errorf(codes.Unauthenticated, "access denied")

	}

	// business logic?

	return &Nothing{}, nil
}

type AdminServerImpl struct {
	subscribers subscribersDB
	authSvc     authService
	UnimplementedAdminServer
}

func NewAdminServerImpl() *AdminServerImpl {
	return &AdminServerImpl{
		subscribers: *NewSubscribersDB(),
		authSvc:     *NewAuthService(),
	}
}
func (as *AdminServerImpl) SetAuth(db AclListType) *AdminServerImpl {
	as.authSvc.SetAuth(db)
	return as
}

// Logging implements AdminServer. Take nothing, produce stream Event
func (as *AdminServerImpl) Logging(_ *Nothing, outStream Admin_LoggingServer) error {
	/*
		service Admin ...
			rpc Logging (Nothing) returns (stream Event) {}
		...
		message Event {
			int64  timestamp = 1;
			string consumer  = 2;
			string method    = 3;
			string host      = 4; // читайте это поле как remote_addr
		}
	*/
	var ctx = outStream.Context()
	ctx = context.WithValue(ctx, "method", "/main.Admin/Logging") // TODO: get value from grpc.*ServerInfo.FullMethod

	// event
	as.pushEvent(ctx)

	// auth
	err := as.auth(ctx)
	if err != nil {
		show("Admin Logging, access denied: ", err)
		return err
	}
	show("Admin Logging, access granted")

	// serve log
	err = as.sendLog(outStream)
	if err != nil {
		show("Admin Logging, sendLog failed: ", err)
		return err
	}

	return nil
}

func (as *AdminServerImpl) Statistics(*StatInterval, Admin_StatisticsServer) error {
	return status.Errorf(codes.Unauthenticated, "Statistics not implemented yet")
}

func (as *AdminServerImpl) pushEvent(ctx context.Context) {
	h := "unknown"
	p, ok := peer.FromContext(ctx)
	if ok {
		h = p.Addr.String()
	}

	evt := Event{
		Method:    ctx.Value("method").(string),
		Consumer:  getConsumer(ctx),
		Timestamp: time.Now().Unix(),
		Host:      h,
	}
	show("event: ", evt.Consumer, evt.Method)

	as.subscribers.Push(&evt)
}

func (as *AdminServerImpl) auth(ctx context.Context) error {
	var consumer string = getConsumer(ctx)
	var method = ctx.Value("method").(string)
	var isAllowed bool = as.authSvc.IsAllowed(consumer, method)
	show("Admin auth; consumer, method, allowed: ", consumer, method, isAllowed)
	if isAllowed {
		return nil
	}
	return status.Errorf(codes.Unauthenticated, "access denied")
}

func (as *AdminServerImpl) sendLog(outStream Admin_LoggingServer) error {
	subscriberId, events := as.subscribers.AddSubscriber()
	defer as.subscribers.RemoveSubscriber(subscriberId)

	for e := range events {
		if err := outStream.Send(e); err != nil {
			show("Admin sendLog; stream.Send failed: ", err)
			return err
		}
	}

	return nil
}

type authService struct {
	authDB AclListType // user: list of methods
	mutex  *sync.RWMutex
}

func NewAuthService() *authService {
	return &authService{
		authDB: make(AclListType, 16),
		mutex:  &sync.RWMutex{},
	}
}
func (as *authService) SetAuth(db AclListType) *authService {
	as.lock()
	defer as.unlock()

	as.authDB = db
	return as
}
func (as *authService) lock() {
	as.mutex.Lock()

}
func (as *authService) unlock() {
	as.mutex.Unlock()
}
func (as *authService) lockRead() {
	as.mutex.RLock()

}
func (as *authService) unlockRead() {
	as.mutex.RUnlock()
}

func (as *authService) IsAllowed(user, method string) bool {
	as.lockRead()
	defer as.unlockRead()

	methods, exist := as.authDB[user]
	if !exist {
		return false
	}

	return slices.ContainsFunc(methods, func(x string) bool {
		if strings.HasSuffix(x, method) {
			return true
		}
		return false
	})
}

type subscribersDB struct {
	mutex    sync.RWMutex
	lastId   uint64 // TODO: should be pool of available id's
	channels map[uint64]chan *Event
}

func NewSubscribersDB() *subscribersDB {
	var db = subscribersDB{}
	return db.Clear()
}

func (db *subscribersDB) Clear() *subscribersDB {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, ch := range db.channels {
		close(ch)
	}
	db.lastId = 0
	db.channels = make(map[uint64]chan *Event, 16)

	return db
}

func (db *subscribersDB) AddSubscriber() (newId uint64, events chan *Event) {
	const queueSize = 1 // TODO: queue size should be configurable
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.lastId++
	db.channels[db.lastId] = make(chan *Event, queueSize)

	return db.lastId, db.channels[db.lastId]
}

func (db *subscribersDB) RemoveSubscriber(id uint64) *subscribersDB {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if ch, exist := db.channels[id]; exist {
		close(ch)
		delete(db.channels, id)
	}

	return db
}

func (db *subscribersDB) Push(e *Event) *subscribersDB {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, ch := range db.channels {
		ch <- e
	}

	return db
}

// GetContextMetaDataValue returns value under given key from context metadata.
// If no data in context or value is empty, return `dflt` value
func GetContextMetaDataValue(ctx context.Context, key, dflt string) string {
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		show("GetContextMetaDataValue failed, context w/o metadata")
		return dflt
	}
	x := strings.Join(md.Get("consumer"), "")
	if x == "" {
		return dflt
	}
	return x
}

func getConsumer(ctx context.Context) string {
	return GetContextMetaDataValue(ctx, "consumer", "")
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
