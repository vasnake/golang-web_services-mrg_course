package gqlgen6

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type Resolver struct{}

// RatePhoto is the resolver for the ratePhoto field.
func (r *mutationResolver) RatePhoto(ctx context.Context, photoID string, direction string) (*Photo, error) {
	panic("not implemented")
}

// UploadPhoto is the resolver for the uploadPhoto field.
func (r *mutationResolver) UploadPhoto(ctx context.Context, comment string, file graphql.Upload) (*Photo, error) {
	panic("not implemented")
}

// ID is the resolver for the id field.
func (r *photoResolver) ID(ctx context.Context, obj *Photo) (string, error) {
	panic("not implemented")
}

// User is the resolver for the user field.
func (r *photoResolver) User(ctx context.Context, obj *Photo) (*User, error) {
	panic("not implemented")
}

// Timeline is the resolver for the timeline field.
func (r *queryResolver) Timeline(ctx context.Context) ([]*Photo, error) {
	panic("not implemented")
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, userID string) (*User, error) {
	panic("not implemented")
}

// Photos is the resolver for the photos field.
func (r *queryResolver) Photos(ctx context.Context, userID string) ([]*Photo, error) {
	panic("not implemented")
}

// Photos is the resolver for the photos field.
func (r *userResolver) Photos(ctx context.Context, obj *User, count int) ([]*Photo, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Photo returns PhotoResolver implementation.
func (r *Resolver) Photo() PhotoResolver { return &photoResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type photoResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
