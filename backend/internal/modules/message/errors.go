package message

import "errors"

var (
	ErrNotFriend             = errors.New("not friend")
	ErrCannotMessageYourself = errors.New("cannot send message to yourself")
	ErrEmptyMessage          = errors.New("message must have content or at least one attachment")
	ErrAttachmentNotFound    = errors.New("attachment not uploaded yet")
	ErrTooManyAttachments    = errors.New("too many attachments")
	ErrMessageNotFound     = errors.New("message not found")
	ErrForbidden           = errors.New("forbidden")
	ErrAlreadyRevoked      = errors.New("message already revoked")
	ErrRevokeWindowExpired = errors.New("revoke window expired")
)
