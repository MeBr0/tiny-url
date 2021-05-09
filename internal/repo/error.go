package repo

import "errors"

var (
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrURLNotFound       = errors.New("url doesn't exists")
	ErrURLAlreadyExists  = errors.New("url already exists")
)
