package scheduler

import (
	"fmt"
	"slices"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

type Constraining struct {
	unscheduledEvents []models.UnscheduledEvent
	candidateEvents   []models.Event
}

func NewConstraining(input *models.Input) (*Constraining, error) {
	candidateEvents, err := makeCandidateEvents(input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate list of candidate events: %w", err)
	}

	unscheduledEvents, err := makeUnscheduledEvents(input)
	if err != nil {
		return nil, fmt.Errorf("failed to generated list of unscheduled events: %w", err)
	}

	return &Constraining{
		unscheduledEvents: unscheduledEvents,
		candidateEvents:   candidateEvents,
	}, nil
}

func makeCandidateEvents(input *models.Input) ([]models.Event, error) {
	firstDayOfMatches := input.FirstDayOfMatches()
	lastDayOfMatches := input.LastDayOfMatches()

	// Go from first to last day of matches (inclusive) and come up with list of candidate events, excluding
	// blackout events
	blackoutEvents := input.Events
	candidateEvents := make([]models.Event, 0)
	for currentDay := *firstDayOfMatches; currentDay.Before(lastDayOfMatches.AddDate(0, 0, 1)); currentDay = currentDay.AddDate(0, 0, 1) {
		// All days have at least two candidate events - morning and evening slots
		morningEvent := makeCandidateEvent(currentDay, models.SlotMorning)
		if !isEventBlackedOut(morningEvent, blackoutEvents) {
			candidateEvents = append(candidateEvents, morningEvent)
		}

		eveningEvent := makeCandidateEvent(currentDay, models.SlotEvening)
		if !isEventBlackedOut(eveningEvent, blackoutEvents) {
			candidateEvents = append(candidateEvents, eveningEvent)
		}

		// Weekends can have an extra candidate event - afternoon slot
		if isWeekend(currentDay) {
			afternoonEvent := makeCandidateEvent(currentDay, models.SlotAfternoon)
			if !isEventBlackedOut(afternoonEvent, blackoutEvents) {
				candidateEvents = append(candidateEvents, afternoonEvent)
			}
		}
	}
	return candidateEvents, nil
}

func makeUnscheduledEvents(input *models.Input) ([]models.UnscheduledEvent, error) {
	// Make list of match events with constraints in each
	events := make([]models.UnscheduledEvent, 0)
	for _, team := range input.Teams {
		for _, week := range team.Weeks {
			event := models.UnscheduledEvent{
				Event: models.Event{
					Title: team.Name,
					Type:  models.EventTypeMatch,
				},
				Constraints: models.Constraints{},
			}

			// Add required constraints
			dayConstraint := models.DayConstraint{NotBefore: week, Before: week.AddDate(0, 0, 7)}
			slotConstraint := models.SlotConstraint{TeamScheduleGroup: team.ScheduleGroup}
			event.Constraints.Required = append(event.Constraints.Required, dayConstraint, slotConstraint)

			// TODO: add preference constraints

			events = append(events, event)
		}
	}

	return events, nil
}

func (c *Constraining) Run() (*models.Schedule, error) {
	// Loop over candidate events, checking if each one fits any of
	// the unscheduled events. If it does, schedule it.
	scheduledEvents := make([]models.Event, 0)
	unscheduledEvents := c.unscheduledEvents
	for _, candidateEvent := range c.candidateEvents {

		removeScheduled(unscheduledEvents, scheduledEvents)
		for _, unscheduledEvent := range c.unscheduledEvents {
			// Make sure ALL required constraints are met
			areAllRequiredConstraintsMet := true
			for _, requiredConstraint := range unscheduledEvent.Constraints.Required {
				if !requiredConstraint.CanSchedule(candidateEvent) {
					areAllRequiredConstraintsMet = false
					break
				}
			}

			if !areAllRequiredConstraintsMet {
				// Candidate event is not a fit for this unscheduled event; move
				// on to next unscheduled event.
				continue
			}

			// TODO: check preference constraints

			// Candidate event is a good fit for this unscheduled event, so let's
			// schedule it. Since candidate event has been scheduled, we can no longer
			// use it to match against the remaining unscheduled events so we break out
			// and start over with the next candidate event and an updated list of
			// unscheduled events.
			scheduledEvents = append(scheduledEvents, candidateEvent)
			break
		}
	}

	// Return scheduled events and also any unscheduled events
	s := models.Schedule{
		ScheduledEvents:   scheduledEvents,
		UnscheduledEvents: unscheduledEvents,
	}
	return &s, nil
}

func isWeekend(day time.Time) bool {
	return day.Weekday() == time.Saturday || day.Weekday() == time.Sunday
}

func makeCandidateEvent(day time.Time, slot models.DaySlot) models.Event {
	return models.Event{
		Type: models.EventTypeMatch,
		Slot: slot,
		Date: day,
	}
}

// TODO: optimize if needed
func isEventBlackedOut(candidateEvent models.Event, blackoutEvents []models.Event) bool {
	for _, blackoutEvent := range blackoutEvents {
		if candidateEvent.OverlapsWith(blackoutEvent) {
			return true
		}
	}

	return false
}

// removeScheduled removes any scheduled events from the unscheduled events list, modifying
// it in place.
func removeScheduled(unscheduledEvents []models.UnscheduledEvent, scheduledEvents []models.Event) {
	for _, scheduledEvent := range scheduledEvents {
		for idx, unscheduledEvent := range unscheduledEvents {
			if unscheduledEvent.Event == scheduledEvent {
				unscheduledEvents = slices.Delete(unscheduledEvents, idx, idx+1)
				break
			}
		}
	}
}
