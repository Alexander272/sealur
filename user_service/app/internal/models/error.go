package models

import "errors"

var (
	ErrPassword   = errors.New("passwords do not match")
	ErrUsersEmpty = errors.New("user list is empty")
	ErrUserExist  = errors.New("user already exists")

	ErrUserNotFound = errors.New("user is not found")
)
