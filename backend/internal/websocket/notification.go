package websocket

import (
	"encoding/json"
	"log"
)

type NotificationSender struct {
	manager *Manager
}

func NewDBNotificationSender(manager *Manager) *NotificationSender {
	return &NotificationSender{
		manager: manager,
	}
}

func (s *NotificationSender) SendNotification(userID int64, data map[string]interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal notification payload: %v", err)
		return
	}
	s.manager.SendToUser(userID, payload)
}

func (s *NotificationSender) IsOnline(userID int64) bool {
	return s.manager.IsOnline(userID)
}
