package service

import "errors"

var (
	ErrNoPossibleAliasEncoding = errors.New("cannot encode url to alias")
	ErrURLLimit                = errors.New("cannot create more urls")
	ErrURLForbidden            = errors.New("url cannot be accessed")
)
