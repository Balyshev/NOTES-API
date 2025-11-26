package models

import "errors"

var (
	ErrUsernameRequired   = errors.New("username is required")
	ErrUsernameTooLong    = errors.New("username must be at most 50 characters")
	ErrUsernameTooShort   = errors.New("username must be at least 3 characters")
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUnauthorized       = errors.New("unauthorized")
)

var (
	ErrTitleRequired   = errors.New("title is required")
	ErrTitleTooLong    = errors.New("title must be at most 255 characters")
	ErrContentRequired = errors.New("content is required")
	ErrNoteNotFound    = errors.New("note not found")
	ErrForbidden       = errors.New("you don't have permission to access this note")
)
