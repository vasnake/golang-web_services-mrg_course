package gqlgen2

import (
	"log"
	"strconv"
)

type Photo struct {
	ID     uint `json:"id"`
	UserID uint `json:"-"`
	// DB have not user object in Photo record, resolver must be used
	// User     *User  `json:"user"`
	URL      string `json:"url"`
	Comment  string `json:"comment"`
	Rating   int    `json:"rating"`
	Liked    bool   `json:"liked"`
	Followed bool   `json:"followed"`
}

// Photo.Id getter, gql need string, have uint in db
func (ph *Photo) Id() string {
	log.Println("call Photo.Id method", ph.ID)
	return strconv.Itoa(int(ph.ID))
}
