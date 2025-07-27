package websocket

import "encoding/json"

type Message struct {
	Type      string `json:"type"`
	To        int64  `json:"to,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func parseMessage(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return &m, err
}
