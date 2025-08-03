package models

import "time"

type EventType string

const (
	EventTypeMatch    EventType = "match"
	EventTypeBlackout           = "blackout"
)

type Event struct {
	Title string    `json:"title"`
	Type  EventType `json:"type"`
	Slot  DaySlot   `json:"slot"`

	Date time.Time `json:"date"`
	//StartTime time.Time
	//EndTime   time.Time
}

func (event Event) OverlapsWith(anotherEvent Event) bool {
	return event.Date == anotherEvent.Date &&
		event.Slot == anotherEvent.Slot
}

func (event Event) IsOnWeekend() bool {
	return event.Date.Weekday() == time.Saturday || event.Date.Weekday() == time.Sunday
}

type ComparableEvent interface {
	ID() string
	IsEqualTo(another ComparableEvent) bool
}

func (event Event) ID() string {
	return event.Date.Format("2006-01-02") + "-" + string(event.Slot)
}

func (event Event) IsEqualTo(another ComparableEvent) bool {
	return event.ID() == another.ID()
}
