// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen6

type Mutation struct {
}

type Query struct {
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Followed bool   `json:"followed"`
	// возвращает фотограции данного пользователя
	Photos []*Photo `json:"photos"`
}
