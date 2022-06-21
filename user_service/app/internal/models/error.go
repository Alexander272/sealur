package models

import "errors"

var (
	ErrUserNotFound = errors.New("unable to find user")
	ErrPassword     = errors.New("passwords do not match")
	ErrUsersEmpty   = errors.New("user list is empty")
)
