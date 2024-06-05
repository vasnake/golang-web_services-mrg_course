package main

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
)

// auth user stuff

func CheckAuthorizedMiddleware(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	sessionRef, err := SessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	show("CheckAuthorizedMiddleware, session: ", sessionRef)
	return next(ctx)
}

func SessionFromContext(ctx context.Context) (*AppSession, error) {
	sess, isSession := ctx.Value(CONTEXT_SESSION_KEY).(*AppSession)
	if !isSession {
		return nil, ERROR_NO_USER
	} else {
		return sess, nil
	}
}

type AppSession interface{}

const (
	CONTEXT_SESSION_KEY = "SESSION_CONTEXT_KEY"
)

var (
	ERROR_NO_USER = errors.New("User not authorized")
)
