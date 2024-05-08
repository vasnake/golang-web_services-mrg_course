package grpc_4_balanced

import (
	"fmt"
	// "gws/7/microservices/grpc/session"
	"math/rand"
	"sync"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/codes"

	"golang.org/x/net/context"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	host     string
}

func NewSessionManager(port string) *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[string]*Session{},
		host:     port,
	}
}

// mustEmbedUnimplementedAuthCheckerServer implements AuthCheckerServer.
func (sm *SessionManager) mustEmbedUnimplementedAuthCheckerServer() {
	panic("unimplemented")
}

func (sm *SessionManager) Create(ctx context.Context, in *Session) (*SessionID, error) {
	fmt.Println("call Create", in)
	id := &SessionID{ID: RandStringRunes(sessKeyLen)}

	sm.mu.Lock()
	sm.sessions[id.ID] = in
	sm.mu.Unlock()
	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *SessionID) (*Session, error) {
	fmt.Println("call Check", in)
	// между сервисами нет общения, возвращаем заглушку
	fakeLogin := sm.host + " " + in.GetID()
	return &Session{Login: fakeLogin}, nil

	// sm.mu.RLock()
	// defer sm.mu.RUnlock()
	// if sess, ok := sm.sessions[*in]; ok {
	// 	return sess, nil
	// }
	// return nil, grpc.Errorf(codes.NotFound, "session not found")
}

func (sm *SessionManager) Delete(ctx context.Context, in *SessionID) (*Nothing, error) {
	fmt.Println("call Delete", in)

	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, in.ID)
	return &Nothing{Dummy: true}, nil
}

const sessKeyLen = 10

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
