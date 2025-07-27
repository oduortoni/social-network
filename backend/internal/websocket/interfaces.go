package websocket

import "net/http"

// SessionResolver resolves the authenticated user ID, nickname, and avatar from the HTTP request.
type SessionResolver interface {
	GetUserFromRequest(r *http.Request) (int64, string, string, error)
}

// GroupMemberFetcher fetches group member IDs from DB
// to allow broadcasting messages to the group.
type GroupMemberFetcher interface {
	GetGroupMemberIDs(groupID string) ([]int64, error)
}

// MessagePersister stores chat messages for retrieval and persistence.
type MessagePersister interface {
	SaveMessage(senderID int64, msg *Message) error
}
