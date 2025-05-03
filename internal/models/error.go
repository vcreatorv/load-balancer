package models

import (
	"errors"
	"fmt"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal server error")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

type Error struct {
	clientErr error
	msg       string
}

func NewError(clientErr error, msg string) error {
	return &Error{
		clientErr: clientErr,
		msg:       msg,
	}
}

func (e Error) Error() string {
	if e.clientErr != nil {
		return fmt.Sprint(e.msg)
	}
	return e.msg
}

func (e Error) ClientErr() error {
	return e.clientErr
}

func (e Error) Msg() string {
	return e.msg
}
