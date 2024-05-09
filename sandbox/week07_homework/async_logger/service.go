package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	// sync "sync"
	"time"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type AclListType map[string][]string

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {
	// show("StartMyMicroservice, addr: ", listenAddr)

	// parse params
	var acl = make(AclListType, 16)
	err := json.Unmarshal([]byte(ACLData), &acl)
	if err != nil {
		show("StartMyMicroservice failed, invalid ACL json: ", err)
		return err
	}
	// show("StartMyMicroservice, ACL parsed: ", acl)

	// setup transport
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		show("StartMyMicroservice failed, net.Listen error: ", err)
		return err
	}
	// register server
	server := grpc.NewServer()
	RegisterBizServer(server, NewBizServerImpl())
	RegisterAdminServer(server, NewAdminServerImpl())
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

/*

type BizClient interface {
	Test(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	Check(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	Add(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
}

message Nothing {
    bool dummy = 1;
}

service Biz {
    rpc Check(Nothing) returns(Nothing) {}
    rpc Add(Nothing) returns(Nothing) {}
    rpc Test(Nothing) returns(Nothing) {}
}

*/

type BizServerImpl struct {
	// acl AclListType
	// mutex sync.RWMutex
	// data  map[string]string
	UnimplementedBizServer
}

func NewBizServerImpl() *BizServerImpl {
	return &BizServerImpl{
		// acl: make(AclListType, 16),
		// mutex: sync.RWMutex{},
		// data:  make(map[string]string, 16),
	}
}

func (bs *BizServerImpl) Check(ctx context.Context, in *Nothing) (*Nothing, error) {
	// return nil, status.Errorf(codes.NotFound, "session not found")
	return &Nothing{}, nil
}

// Add implements BizServer.
func (bs *BizServerImpl) Add(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

// Test implements BizServer.
func (bs *BizServerImpl) Test(ctx context.Context, _ *Nothing) (*Nothing, error) {
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return nil, status.Errorf(codes.InvalidArgument, "context w/o metadata")
	}
	consumer := strings.Join(md.Get("consumer"), "")
	show("biz test, consumer: ", consumer)
	if consumer == "biz_admin" {
		return &Nothing{}, nil
	}

	return nil, status.Errorf(codes.Unauthenticated, "access denied")
}

/*
	service Admin {
	    rpc Logging (Nothing) returns (stream Event) {}
	    rpc Statistics (StatInterval) returns (stream Stat) {}
	}
*/

type AdminServerImpl struct {
	subscribers subscribersDB
	UnimplementedAdminServer
}

func NewAdminServerImpl() *AdminServerImpl {
	var as = AdminServerImpl{
		subscribers: *NewSubscribersDB(),
	}

	return &as
}

// Logging implements AdminServer.
func (as *AdminServerImpl) Logging(_ *Nothing, srv Admin_LoggingServer) error {
	var ctx = srv.Context()
	show("Logging, context: ", ctx)
	return status.Errorf(codes.Unauthenticated, "Statistics not implemented yet")

	// subscriberId, events := as.subscribers.AddSubscriber()
	// defer as.subscribers.RemoveSubscriber(subscriberId)
	// for e := range events {
	// 	if err := srv.Send(e); err != nil {
	// 		return err
	// 	}
	// }
	// return nil
}

// Statistics implements AdminServer.
func (as *AdminServerImpl) Statistics(*StatInterval, Admin_StatisticsServer) error {
	return status.Errorf(codes.Unauthenticated, "Statistics not implemented yet")
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

func (db *subscribersDB) RemoveSubscriber(id uint64) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if ch, exist := db.channels[id]; exist {
		close(ch)
		delete(db.channels, id)
	}
}

func (db *subscribersDB) Push(e *Event) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, ch := range db.channels {
		ch <- e
	}
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
