package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	ws "github.com/tajjjjr/social-network/backend/internal/websocket"
)

type WebSocketHandler struct {
	Manager  *ws.Manager
	Upgrader websocket.Upgrader
}

func NewWebSocketHandler(m *ws.Manager) *WebSocketHandler {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &WebSocketHandler{
		Manager:  m,
		Upgrader: upgrader,
	}
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade", http.StatusBadRequest)
		return
	}

	userID, nickname, avatar, err := h.Manager.Resolver.GetUserFromRequest(r)
	if err != nil {
		conn.Close()
		return
	}

	client := ws.NewClient(userID, nickname, avatar, conn)
	h.Manager.Register(client)
	defer h.Manager.Unregister(userID)
	defer conn.Close()

	go h.Manager.WritePump(client)
	h.Manager.ReadPump(client)
}
