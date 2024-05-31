package user

import (
	"strconv"
)

type User struct {
	ID       uint32
	Login    string `gqlgen:"name"`
	Email    string
	Ver      int32
	Followed *bool
}

func (u *User) GetID() uint32 {
	return u.ID
}

func (u *User) GetVer() int32 {
	return u.Ver
}

func (u *User) Id() string {
	return strconv.Itoa(int(u.ID))
}

func (u *User) Name() string {
	return u.Login
}

func (u *User) Avatar() string {
	return "https://via.placeholder.com/80"
}
