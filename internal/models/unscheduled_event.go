package models

import "time"

type UnscheduledEvent struct {
	Event
	Constraints Constraints `json:"constraints"`
}

type Constraints struct {
	Required    []FilterConstraint  `json:"required"`
	Preferences []ScoringConstraint `json:"preferences"`
}

type FilterConstraint interface {
	CanSchedule(candidateEvent Event) bool
}

type ScoringConstraint interface {
	Score(candidateEvent Event) float64
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

// TODO: what about broader considerations, e.g. all UnscheduledEvents or all
// candidate UnscheduledEvents?
func (ue UnscheduledEvent) Match(candidateEvent Event) float64 {
	// Return a score of zero if candidate event fails to meet
	// *any* of the unscheduled event's required constraints.
	for _, fc := range ue.Constraints.Required {
		if !fc.CanSchedule(candidateEvent) {
			return 0.0
		}
	}

	var finalScore float64
	for _, sc := range ue.Constraints.Preferences {
		score := sc.Score(candidateEvent)

		// TODO: think about how to combine scores from scoring constraints
		finalScore = score
	}

	return finalScore
}
