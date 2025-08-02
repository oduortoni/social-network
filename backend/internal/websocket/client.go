package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID        int64
	Nickname  string
	Avatar    string
	Conn      *websocket.Conn
	Send      chan []byte
	Connected time.Time
}

func NewClient(id int64, nickname string, avatar string, conn *websocket.Conn) *Client {
	return &Client{
		ID:        id,
		Nickname:  nickname,
		Avatar:    avatar,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Connected: time.Now(),
	}
}
