package graphql

//go:generate go run github.com/99designs/gqlgen -v

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"photolist/pkg/photos"
	"photolist/pkg/session"
	"photolist/pkg/user"
)

type Resolver struct {
	UsersRepo  *user.UserRepository
	PhotosRepo *photos.StDb
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Photo() PhotoResolver {
	return &photoResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RatePhoto(ctx context.Context, idStr string, direction string) (*photos.Photo, error) {
	// log.Println("call mutationResolver.RatePhoto method with id", id, direction)
	rate := 1
	if direction != "up" {
		rate = -1
	}

	sess, _ := session.SessionFromContext(ctx)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("bad id")
	}

	err = r.PhotosRepo.Rate(uint32(id), sess.UserID, rate)
	if err != nil {
		log.Println("PhotosRepo.Rate err:", err)
		return nil, fmt.Errorf("db err")
	}

	return r.PhotosRepo.GetByID(uint32(id), sess.UserID)
}

func (r *mutationResolver) FollowUser(ctx context.Context, userIDStr string, direction string) (*user.User, error) {
	sess, _ := session.SessionFromContext(ctx)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("bad id")
	}

	folUser, err := r.UsersRepo.GetByID(uint32(userID))
	if user.IsErrUserNotFound(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	rate := 1
	if direction == "down" {
		rate = -1
	}

	err = r.UsersRepo.Follow(folUser.ID, sess.UserID, rate)
	if err == nil {
		return nil, err
	}
	return folUser, nil
}

func (r *mutationResolver) UploadPhoto(ctx context.Context, comment string, file graphql.Upload) (*photos.Photo, error) {
	sess, _ := session.SessionFromContext(ctx)

	md5Sum, err := photos.SaveFile(file.File)
	if err != nil {
		return nil, err
	}
	realFile := "./images/" + md5Sum + ".jpg"
	err = photos.MakeThumbnails(realFile, md5Sum)
	if err != nil {
		return nil, err
	}

	ph := &photos.Photo{
		UserID:  sess.UserID,
		Comment: comment,
		URL:     md5Sum,
	}
	ph.ID, err = r.PhotosRepo.Add(ph)
	if err != nil {
		return nil, fmt.Errorf("db err")
	}
	return ph, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Photos(ctx context.Context, obj *user.User, count *int) ([]*photos.Photo, error) {
	sess, _ := session.SessionFromContext(ctx)
	// TODO handle count param
	return r.PhotosRepo.GetPhotos(obj.ID, sess.UserID)
}

func (r *userResolver) Followed(ctx context.Context, obj *user.User) (bool, error) {
	if obj.Followed != nil {
		return *obj.Followed, nil
	}
	sess, _ := session.SessionFromContext(ctx)
	return r.UsersRepo.IsFollowed(obj.ID, sess.UserID)
}

func (r *userResolver) FollowedUsers(ctx context.Context, obj *user.User, count *int) ([]*user.User, error) {
	// sess, _ := session.SessionFromContext(ctx)
	// TODO handle count param
	return r.UsersRepo.GetFollowedUsers(obj.ID)
}

func (r *userResolver) RecomendedUsers(ctx context.Context, obj *user.User, count *int) ([]*user.User, error) {
	// sess, _ := session.SessionFromContext(ctx)
	// TODO handle count param
	return r.UsersRepo.GetRecomendedUsers(obj.ID)
}

type photoResolver struct{ *Resolver }

func (r *photoResolver) User(ctx context.Context, obj *photos.Photo) (*user.User, error) {
	return ctx.Value("userLoaderKey").(*user.UserLoader).Load(obj.UserID)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Timeline(ctx context.Context) ([]*photos.Photo, error) {
	sess, _ := session.SessionFromContext(ctx)
	// TODO handle count param
	return r.PhotosRepo.GetPhotos(sess.UserID, sess.UserID)
}

func (r *queryResolver) User(ctx context.Context, userIDStr string) (*user.User, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("bad id")
	}
	return r.UsersRepo.GetByID(uint32(userID))
}

func (r *queryResolver) Me(ctx context.Context) (*user.User, error) {
	sess, _ := session.SessionFromContext(ctx)
	return r.UsersRepo.GetByID(sess.UserID)
}

func (r *queryResolver) Photo(ctx context.Context, photoIDStr string) (*photos.Photo, error) {
	sess, _ := session.SessionFromContext(ctx)
	id, err := strconv.Atoi(photoIDStr)
	if err != nil {
		return nil, fmt.Errorf("bad id")
	}
	return r.PhotosRepo.GetByID(uint32(id), sess.UserID)
}

func (r *queryResolver) Photos(ctx context.Context, userIDStr string) ([]*photos.Photo, error) {
	sess, _ := session.SessionFromContext(ctx)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("bad id")
	}
	return r.PhotosRepo.GetPhotos(uint32(userID), sess.UserID)
}
