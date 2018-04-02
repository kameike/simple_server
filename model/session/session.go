package session

import (
	"github.com/kameike/simple_server/model/user"
)

type Session string

func Create(_u user.User) Session {
	return "hashLw"
}

func (_session Session) Check() bool {
	return false
}
