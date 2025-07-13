package models

import (
	"math/rand"
	"time"
)

type UnscheduledEvent struct {
	Event
	Constraints Constraints `json:"constraints"`
}

type Constraints struct {
	Required    []FilterConstraint        `json:"required"`
	Preferences []ProbabilisticConstraint `json:"preferences"`
}

type FilterConstraint interface {
	CanSchedule(candidateEvent Event) bool
}

type ProbabilisticConstraint interface {
	ShouldSchedule(candidateEvent Event) bool
}

type SlotConstraint struct{ TeamScheduleGroup TeamScheduleGroup }

func (sc SlotConstraint) CanSchedule(candidateEvent Event) bool {
	// Get allowed slots for candidate event
	allowedSlots := sc.TeamScheduleGroup.AllowedSlots(candidateEvent.IsOnWeekend())

	// Check if any of the allowed slots matches the candidate event's slot
	for _, allowedSlot := range allowedSlots {
		if candidateEvent.Slot == allowedSlot {
			return true
		}
	}
	return false
}

type DayConstraint struct {
	NotBefore time.Time
	Before    time.Time
}

func (dc DayConstraint) CanSchedule(candidateEvent Event) bool {
	return !candidateEvent.Date.Before(dc.NotBefore) && candidateEvent.Date.Before(dc.Before)
}

type DayPreferenceConstraint struct {
	Probabilities map[time.Weekday]float64
}

func (dpc DayPreferenceConstraint) ShouldSchedule(candidateEvent Event) bool {
	probability := dpc.Probabilities[candidateEvent.Date.Weekday()]
	return rand.Float64() <= probability
}

func (ue UnscheduledEvent) MatchRequired(candidateEvent Event) bool {
	// Return false if candidate event fails to meet
	// *any* of the unscheduled event's required constraints.
	for _, fc := range ue.Constraints.Required {
		if !fc.CanSchedule(candidateEvent) {
			return false
		}
	}
	return true
}

func (ue UnscheduledEvent) MatchPreferences(candidateEvent Event) bool {
	// Return true if candidate event meets *any* of the unscheduled
	// event's preference constraints.
	for _, pc := range ue.Constraints.Preferences {
		if pc.ShouldSchedule(candidateEvent) {
			return true
		}
	}
	return false
}
