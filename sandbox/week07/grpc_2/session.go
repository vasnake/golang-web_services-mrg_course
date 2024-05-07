package grpc_2

import (
	"fmt"
	// "gws/7/microservices/grpc/session"
	"math/rand"
	"sync"
	"week07/grpc_1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

const sessKeyLen = 10

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*grpc_1.Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[string]*grpc_1.Session{},
	}
}

// mustEmbedUnimplementedAuthCheckerServer implements grpc_1.AuthCheckerServer.
func (sm *SessionManager) mustEmbedUnimplementedAuthCheckerServer() {
	// panic("unimplemented")
}

func (sm *SessionManager) Create(ctx context.Context, in *grpc_1.Session) (*grpc_1.SessionID, error) {
	fmt.Println("call Create", in)

	// metadata demo: prefix, suffix
	header := metadata.Pairs("header-key", "42")
	grpc.SendHeader(ctx, header)
	trailer := metadata.Pairs("trailer-key", "3.14")
	grpc.SetTrailer(ctx, trailer)

	// do work
	id := &grpc_1.SessionID{ID: RandStringRunes(sessKeyLen)}
	sm.mu.Lock()
	sm.sessions[id.ID] = in
	sm.mu.Unlock()
	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *grpc_1.SessionID) (*grpc_1.Session, error) {
	fmt.Println("call Check", in)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[in.ID]; ok {
		return sess, nil
	}
	return nil, grpc.Errorf(codes.NotFound, "session not found")
}

func (sm *SessionManager) Delete(ctx context.Context, in *grpc_1.SessionID) (*grpc_1.Nothing, error) {
	fmt.Println("call Delete", in)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, in.ID)
	return &grpc_1.Nothing{Dummy: true}, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
