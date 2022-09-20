package gqlgen3

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	PhotosData map[string]*Photo
	Users      map[uint]*User
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Photo() PhotoResolver {
	return &photoResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RatePhoto(ctx context.Context, id string, direction string) (*Photo, error) {
	log.Println("call mutationResolver.RatePhoto method with id", id, direction)
	rate := 1
	if direction != "up" {
		rate = -1
	}
	ph, ok := r.PhotosData[id]
	if !ok {
		return nil, fmt.Errorf("no photo %v found", id)
	}
	ph.Rating += rate
	return ph, nil
}

type photoResolver struct{ *Resolver }

func (r *photoResolver) ID(ctx context.Context, obj *Photo) (string, error) {
	return obj.Id(), nil
}

// get user data loader from context, call loader, wait for result
func (r *photoResolver) User(ctx context.Context, obj *Photo) (*User, error) {
	// return r.Users[obj.UserID], nil
	log.Println("call photoResolver.User", obj.UserID)
	start := time.Now()
	// call and wait
	user, err := ctx.Value("userLoaderKey").(*UserLoader).Load(obj.UserID)
	// got batch or cached value
	log.Println("get photoResolver.User", obj.UserID, "from UserLoader, time ", time.Since(start))
	return user, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Timeline(ctx context.Context) ([]*Photo, error) {
	log.Println("call queryResolver.Timeline with ctx.userID", ctx.Value("userID"))
	items := []*Photo{}
	for _, ph := range r.PhotosData {
		items = append(items, ph)
	}
	return items, nil
}

func (r *queryResolver) User(ctx context.Context, userID string) (*User, error) {
	log.Println("call queryResolver.User for", userID)
	id, _ := strconv.Atoi(userID)
	return r.Users[uint(id)], nil
}

func (r *queryResolver) Photos(ctx context.Context, userID string) ([]*Photo, error) {
	log.Println("call queryResolver.Photos")
	id, _ := strconv.Atoi(userID)
	items := []*Photo{}
	for _, ph := range r.PhotosData {
		if ph.UserID != uint(id) {
			continue
		}
		items = append(items, ph)
	}
	return items, nil
}
