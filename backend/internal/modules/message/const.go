package message

import "time"

const (
	defaultHistoryLimit      = 20
	maxHistoryLimit          = 100
	maxAttachmentsPerMessage = 10
	presignGetExpiry         = 15 * time.Minute
	revokeWindow             = 15 * time.Minute
	eventNewMessage          = "message.new"
	eventMessageDeleted      = "message.deleted"
	eventMessageRevoked      = "message.revoked"
)
