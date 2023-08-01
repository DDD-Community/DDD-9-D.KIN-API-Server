package user

import "errors"

var (
	ErrInvalidNickname       = errors.New("invalid nickname")
	ErrNicknameAlreadyExists = errors.New("nickname already exists")
	ErrUserNotFound          = errors.New("user not found")
)
