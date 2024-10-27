package models

import "time"

type Event struct {
	title string

	startTime time.Time
	endTime   time.Time
}

func (event Event) OverlapsWith(anotherEvent Event) bool {
	if event.Includes(anotherEvent.startTime) {
		return true
	}

	if event.Includes(anotherEvent.endTime) {
		return true
	}

	if anotherEvent.Includes(event.startTime) {
		return true
	}

	if anotherEvent.Includes(event.endTime) {
		return true
	}

	return false
}

func (event Event) Includes(t time.Time) bool {
	if t.Before(event.startTime) {
		return false
	}

	if t.After(event.endTime) {
		return false
	}

	return true
}
