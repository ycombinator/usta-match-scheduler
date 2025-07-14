package scheduler

import (
	"math/rand"
	"slices"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

type Constraining struct {
	unscheduledEvents []models.UnscheduledEvent
	candidateEvents   []models.Event
}

func NewConstraining(input *models.Input) *Constraining {
	candidateEvents := makeCandidateEvents(input)
	unscheduledEvents := makeUnscheduledEvents(input)

	return &Constraining{
		unscheduledEvents: unscheduledEvents,
		candidateEvents:   candidateEvents,
	}
}

func makeCandidateEvents(input *models.Input) []models.Event {
	firstDayOfMatches := input.FirstDayOfMatches()
	lastDayOfMatches := input.LastDayOfMatches()

	if firstDayOfMatches == nil || lastDayOfMatches == nil {
		return []models.Event{}
	}

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
	return candidateEvents
}

func makeUnscheduledEvents(input *models.Input) []models.UnscheduledEvent {
	dayPreferenceProbabilities := computeDayPreferenceProbabilities(input.Teams)

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

			// Add preference constraints
			dayPreferenceConstraint := models.DayPreferenceConstraint{
				Probabilities: dayPreferenceProbabilities,
				PreferredDays: team.DayPreferences,
			}
			event.Constraints.Preferences = append(event.Constraints.Preferences, dayPreferenceConstraint)

			events = append(events, event)
		}
	}

	return events
}

// Run generates a schedule in three passes:
//   - In the first pass, it tries to schedule each unscheduled event with a candidate event
//     that matches its required constraints and preferences, with the latter being based on
//     probabilities..
//   - In the second pass, it tries to schedule any remaining unscheduled events with the
//     remaining candidate events, again matching required constraints and preferences, but this
//     time, for the latter, it does not use probabilities, just the first candidate event that
//     matches.
//   - The first and second passes ensure that every preference is satisfied at least once. The
//     third pass finally schedules any remaining unscheduled events with the remaining candidate
//     events, matching only required constraints.
func (c *Constraining) Run() (*models.Schedule, error) {
	// First pass: Loop over candidate events, checking if each one fits any of
	// the unscheduled events. If it does, schedule it. Also, keep track
	// of any unscheduled events that could not be scheduled.
	scheduledEvents := make([]models.Event, 0)
	pass1UnscheduledEvents := make([]models.UnscheduledEvent, 0)

	candidateEvents := make([]models.Event, len(c.candidateEvents))
	copy(candidateEvents, c.candidateEvents)

	for _, unscheduledEvent := range c.unscheduledEvents {
		candidateIdxToRemove := -1
		for candidateIdx, candidateEvent := range candidateEvents {
			if !unscheduledEvent.MatchRequired(candidateEvent) {
				// Candidate event is not a fit for this unscheduled event; move
				// on to next candidate event.
				continue
			}

			if !unscheduledEvent.MatchPreferences(candidateEvent, true) {
				// Candidate event is not fit for this unscheduled event; move
				// on to next candidate event.
				continue
			}

			// Candidate event is a good fit for this unscheduled event, so let's
			// schedule it. Since candidate event has been scheduled, we can no longer
			// use it to match against the remaining unscheduled events. So we break out
			// and remove the scheduled candidate event from the list of candidate events.
			candidateEvent.Title = unscheduledEvent.Event.Title
			scheduledEvents = append(scheduledEvents, candidateEvent)
			candidateIdxToRemove = candidateIdx
			break
		}

		if candidateIdxToRemove > -1 {
			// Event was scheduled, so we remove the corresponding candidate event.
			candidateEvents = removeFromEvents[models.Event](candidateEvents, candidateIdxToRemove)
		} else {
			// Event did not get scheduled, so we add it to the list of unscheduled events.
			pass1UnscheduledEvents = append(pass1UnscheduledEvents, unscheduledEvent)
		}
	}

	// Second pass: At this point, we may still have some unscheduled events so let's go
	// ahead and schedule them taking their required constraints into account and letting the
	// first candidate event that matches the unscheduled event's preference "win".
	// Note that we only consider any remaining candidate events for scheduling. Also, we
	// randomize these remaining candidate events so scheduling isn't front-heavy.
	candidateEvents = randomizeSlice(candidateEvents)

	pass2UnscheduledEvents := make([]models.UnscheduledEvent, 0)
	for _, unscheduledEvent := range pass1UnscheduledEvents {
		candidateIdxToRemove := -1
		for candidateIdx, candidateEvent := range candidateEvents {
			if !unscheduledEvent.MatchRequired(candidateEvent) {
				// Candidate event is not a fit for this unscheduled event; move
				// on to next unscheduled event.
				continue
			}

			if !unscheduledEvent.MatchPreferences(candidateEvent, false) {
				// Candidate event is not fit for this unscheduled event; move
				// on to next candidate event.
				continue
			}

			// Candidate event is a good fit for this unscheduled event, so let's
			// schedule it.
			candidateEvent.Title = unscheduledEvent.Event.Title
			scheduledEvents = append(scheduledEvents, candidateEvent)
			candidateIdxToRemove = candidateIdx
			break
		}

		if candidateIdxToRemove > -1 {
			// Event was scheduled, so we remove the corresponding candidate event.
			candidateEvents = removeFromEvents[models.Event](candidateEvents, candidateIdxToRemove)
		} else {
			// Event was not scheduled
			pass2UnscheduledEvents = append(pass2UnscheduledEvents, unscheduledEvent)
		}
	}

	// Third pass: At this point, we may still have some unscheduled events so let's go
	// ahead and schedule them taking only their required constraints into account.
	pass3UnscheduledEvents := make([]models.UnscheduledEvent, 0)
	for _, unscheduledEvent := range pass2UnscheduledEvents {
		candidateIdxToRemove := -1
		for candidateIdx, candidateEvent := range candidateEvents {
			if !unscheduledEvent.MatchRequired(candidateEvent) {
				// Candidate event is not a fit for this unscheduled event; move
				// on to next unscheduled event.
				continue
			}

			// Candidate event is a good fit for this unscheduled event, so let's
			// schedule it.
			candidateEvent.Title = unscheduledEvent.Event.Title
			scheduledEvents = append(scheduledEvents, candidateEvent)
			candidateIdxToRemove = candidateIdx
			break
		}

		if candidateIdxToRemove > -1 {
			// Event was scheduled, so we remove the corresponding candidate event.
			candidateEvents = removeFromEvents[models.Event](candidateEvents, candidateIdxToRemove)
		} else {
			// Event was not scheduled
			pass3UnscheduledEvents = append(pass3UnscheduledEvents, unscheduledEvent)
		}
	}

	// Return scheduled events and also any unscheduled events (there should be none
	// of these at this point but better to return any than silently dropping them).
	s := models.Schedule{
		ScheduledEvents:   scheduledEvents,
		UnscheduledEvents: pass3UnscheduledEvents,
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

// removeFromEvents removes the event with the given index from the
// given events list and returns the updated unscheduled events list
func removeFromEvents[T any](events []T, idxToRemove int) []T {
	if idxToRemove < 0 {
		return events
	}
	if len(events) == 0 {
		return events
	}
	return slices.Delete(events, idxToRemove, idxToRemove+1)
}

// computeDayPreferenceProbabilities returns a map of weekdays to a [0,1) probability. The probability depends
// on the number of teams that prefer that particular weekday.  The more teams that prefer a weekday, the lower
// the probability associated with it.
func computeDayPreferenceProbabilities(teams []models.SchedulingTeam) map[time.Weekday]float64 {
	// Compute day preference frequencies
	weekdayPreferenceFrequencies := map[time.Weekday]int{}
	for _, team := range teams {
		for _, preferredDay := range team.DayPreferences {
			weekdayPreferenceFrequencies[preferredDay]++
		}
	}

	// Invert day preference frequencies to arrive at probabilities
	probabilities := map[time.Weekday]float64{}
	for preferredDay, frequency := range weekdayPreferenceFrequencies {
		if frequency == 0 {
			// No team prefers this day, leave probability as zero
			continue
		}

		probabilities[preferredDay] = 1 / float64(frequency)
	}

	return probabilities
}

func randSliceOfIntegers(size int) []int {
	seen := make(map[int]struct{}, size)

	output := make([]int, 0)
	for len(output) < size {
		value := rand.Intn(size)
		if _, exists := seen[value]; exists {
			// Value is already in output slice; try again
			continue
		}

		seen[value] = struct{}{}
		output = append(output, value)
	}

	return output
}

func randomizeSlice[T any](slice []T) []T {
	randomizedIndices := randSliceOfIntegers(len(slice))
	output := make([]T, len(slice))
	for idx, randomizedIdx := range randomizedIndices {
		output[randomizedIdx] = slice[idx]
	}
	return output
}
