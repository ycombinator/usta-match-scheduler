package models

import "time"

type Schedule struct {
	ScheduledEvents   []Event            `json:"scheduled_events"`
	UnscheduledEvents []UnscheduledEvent `json:"unscheduled_events"`
}

//type Schedule struct {
//	startDate time.Time
//	endDate   time.Time
//
//	schedule map[time.Time]*Day
//
//	iter time.Time
//}

type Day struct {
	Date        time.Time
	DaytimeTeam *Team
	EveningTeam *Team
}

type ScheduleIterator time.Time

func NewSchedule() *Schedule {
	s := new(Schedule)
	s.ScheduledEvents = make([]Event, 0)
	s.UnscheduledEvents = make([]UnscheduledEvent, 0)
	return s
}
