package models

import "errors"

var (
	ErrStandAlreadyExists  = errors.New("stand with such title already exists")
	ErrFlangeAlreadyExists = errors.New("flange with such title or short already exists")

	ErrSessionEmpty     = errors.New("user session not found")
	ErrClientIPNotFound = errors.New("client ip not found")
	ErrToken            = errors.New("tokens do not match")
)
