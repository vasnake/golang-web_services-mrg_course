package session

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"photolist/pkg/middleware"
	"photolist/pkg/utils/traceutils"
)

var (
	_ SessionManager = (*SessionsGRPC)(nil)
)

type SessionsGRPC struct {
	client AuthClient
}

func NewSessionsGRPC(addr string) (*SessionsGRPC, error) {
	grcpConn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cant connect to grpc")
	}
	return &SessionsGRPC{
		client: NewAuthClient(grcpConn),
	}, nil
}

func ctxWithTrace(ctx context.Context, opName string) (opentracing.Span, context.Context) {
	span, newCtx := opentracing.StartSpanFromContext(ctx, opName)

	md := metadata.Pairs("X-Request-ID", middleware.RequestIDFromContext(ctx))
	ext.Component.Set(span, "grpc-session")

	mdWriter := traceutils.MetadataReaderWriter{md}
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		mdWriter,
	)
	log.Println("----- ctxWithTrace md", md)
	return span, metadata.NewOutgoingContext(newCtx, md)
}

func (sm *SessionsGRPC) Check(ctx context.Context, r *http.Request) (*Session, error) {
	sp, grpcCtx := ctxWithTrace(ctx, "Session.Check")
	defer sp.Finish()

	sessionCookie, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	authSess, err := sm.client.Check(grpcCtx, &AuthCheckIn{SessKey: sessionCookie.Value})
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:     authSess.GetID(),
		UserID: authSess.GetUserID(),
	}, nil
}

func (sm *SessionsGRPC) Create(ctx context.Context, w http.ResponseWriter, user UserInterface) error {
	sp, grpcCtx := ctxWithTrace(ctx, "Session.Create")
	defer sp.Finish()

	authSess, err := sm.client.Create(grpcCtx, &AuthUserIn{
		UserID: user.GetID(),
		Ver:    user.GetVer(),
	})
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   authSess.GetID(),
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func (sm *SessionsGRPC) DestroyCurrent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sp, grpcCtx := ctxWithTrace(ctx, "Session.DestroyCurrent")
	defer sp.Finish()

	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	_, err := sm.client.DestroyCurrent(grpcCtx, &AuthSession{
		ID: cookie.Value,
	})
	return err
}

func (sm *SessionsGRPC) DestroyAll(ctx context.Context, w http.ResponseWriter, user UserInterface) error {
	sp, grpcCtx := ctxWithTrace(ctx, "Session.DestroyAll")
	defer sp.Finish()

	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	_, err := sm.client.DestroyAll(grpcCtx, &AuthUserIn{
		UserID: user.GetID(),
		Ver:    user.GetVer(),
	})
	return err
}
