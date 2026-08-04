package friend

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrAlreadyFriend     = errors.New("already friend")
	ErrRequestExists     = errors.New("friend request already exists")
	ErrCannotAddYourself = errors.New("cannot add yourself")
	ErrRequestNotFound   = errors.New("friend request not found")
	ErrForbidden         = errors.New("forbidden")
	ErrNotFriend         = errors.New("not friend")
)
