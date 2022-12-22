package models

import "errors"

var (
	ErrStandAlreadyExists  = errors.New("stand with such title already exists")
	ErrFlangeAlreadyExists = errors.New("flange with such title or short already exists")
)
