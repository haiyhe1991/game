package logic

import "github.com/pkg/errors"

const(
	UserStatusDelOff = 0
	UserStatusDelOn =  1

	UserBanStatusOn = 1
	UserBanStatusOff = 0
)

var(
	ErrLoginFailed =  errors.New("user name or password wrong")
)