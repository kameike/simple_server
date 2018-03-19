package session

import (
	"github.com/kameike/simple_server/model/user"
)

type Session string

func Create(u user.User) Session {
	u.Name
	return "hashLw"
}

func (session Session) Check() bool {
}
