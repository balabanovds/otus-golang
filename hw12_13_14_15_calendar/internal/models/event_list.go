package models

import "time"

type EventsList struct {
	List []Event `json:"result"`
	Time int64   `json:"time"`
	Len  int     `json:"length"`
}

func NewEventsList(list []Event) EventsList {
	return EventsList{
		List: list,
		Time: time.Now().UnixNano(),
		Len:  len(list),
	}
}
