package models

import "time"

type MQNotification struct {
	EventID int       `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  int       `json:"user_id"`
}

func NewMQNotification(e Event) MQNotification {
	return MQNotification{
		EventID: e.ID,
		Title:   e.Title,
		Date:    e.StartAt,
		UserID:  e.UserID,
	}
}
