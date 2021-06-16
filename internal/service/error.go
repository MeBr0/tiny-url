package service

import "errors"

var (
	ErrNoPossibleAliasEncoding = errors.New("cannot encode url to alias")
	ErrURLExpired              = errors.New("url expired")
	ErrURLLimit                = errors.New("cannot create more urls")
)
