package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"net"
	sync "sync"
	"time"

	grpc "google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type AclListType map[string][]string

func StartMyMicroservice(ctx context.Context, listenAddr, ACLData string) error {
	show("StartMyMicroservice, addr: ", listenAddr)

	// TODO: parse acl json: map[login:str]listUrls:listStr
	var acl = make(AclListType, 16)
	err := json.Unmarshal([]byte(ACLData), &acl)
	if err != nil {
		show("StartMyMicroservice failed, invalid ACL json: ", err)
		return err
	}
	// show("StartMyMicroservice, ACL parsed: ", acl)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		show("StartMyMicroservice failed, net.Listen error: ", err)
		return err
	}

	server := grpc.NewServer()
	RegisterBizServer(server, NewBizManager())
	show("starting gRPC Biz server at " + listenAddr)

	go server.Serve(listener)

	return nil
}

/*
add implementation
type BizClient interface {
	Test(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	Check(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	Add(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
}

*/

type BizManager struct {
	acl   AclListType
	mutex sync.RWMutex
	data  map[string]string
}

func NewBizManager() *BizManager {
	return &BizManager{
		acl:   make(AclListType, 16),
		mutex: sync.RWMutex{},
		data:  make(map[string]string, 16),
	}
}

func (sm *BizManager) Check(ctx context.Context, in *Nothing) (*Nothing, error) {
	// return nil, status.Errorf(codes.NotFound, "session not found")
	return nil, nil
}

// Add implements BizServer.
func (sm *BizManager) Add(context.Context, *Nothing) (*Nothing, error) {
	return nil, nil
}

// Test implements BizServer.
func (sm *BizManager) Test(context.Context, *Nothing) (*Nothing, error) {
	return nil, nil
}

// mustEmbedUnimplementedBizServer implements BizServer.
func (sm *BizManager) mustEmbedUnimplementedBizServer() {
	panic("unimplemented")
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
