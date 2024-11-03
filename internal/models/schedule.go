package models

import "time"

type Schedule struct {
	startDate time.Time
	endDate   time.Time

	schedule map[time.Time]*Day

	iter time.Time
}

type Day struct {
	Date        time.Time
	DaytimeTeam *Team
	EveningTeam *Team
}

type ScheduleIterator time.Time

func NewSchedule(startDate, endDate time.Time) *Schedule {
	s := new(Schedule)
	s.startDate = startDate
	s.endDate = endDate
	s.schedule = map[time.Time]*Day{}

	return s
}

func (s *Schedule) ForDay(dt time.Time) *Day {
	// Initialize if needed
	if _, exists := s.schedule[dt]; !exists {
		s.schedule[dt] = &Day{
			Date: dt,
		}
	}

	return s.schedule[dt]
}

func (s *Schedule) ResetIterator() {
	s.iter = s.startDate
}

func (s *Schedule) Next() *Day {
	if s.iter.After(s.endDate) {
		return nil
	}

	next := s.schedule[s.iter]
	s.iter = s.iter.AddDate(0, 0, 1)
	return next
}

func (d *Day) HasDaytimeCapacity() bool {
	return d.DaytimeTeam == nil
}

func (d *Day) HasEveningCapacity() bool {
	return d.EveningTeam == nil
}

func (d *Day) HasCapacity() bool {
	return d.HasDaytimeCapacity() && d.HasEveningCapacity()
}
