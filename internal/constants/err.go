package constants

import "errors"

var (
	ErrNotTSV     = errors.New("not a tsv file")
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)
