package gqlgen6

import (
	"log"
	"strconv"
)

type Photo struct {
	ID     uint `json:"id"`
	UserID uint `json:"-"` // для демонстрации превращения модель.айди в графкл.юзер
	// User     *User  `json:"user"`
	URL     string `json:"url"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
	Liked   bool   `json:"liked"`
}

func (ph *Photo) Id() string {
	log.Println("Photo.Id(): ", ph.ID)
	return strconv.Itoa(int(ph.ID))
}
