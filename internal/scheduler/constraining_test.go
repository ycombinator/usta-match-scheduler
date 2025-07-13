package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

func TestIsEventBlackedOut(t *testing.T) {
	today := time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC)
	tomorrow := time.Date(2025, 7, 13, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		candidateEvent models.Event
		blackoutEvents []models.Event
		expected       bool
	}{
		"no_blackout_events": {
			candidateEvent: makeCandidateEvent(today, models.SlotMorning),
			blackoutEvents: nil,
			expected:       false,
		},
		"no_matching_blackout_events_same_day_different_slot": {
			candidateEvent: makeCandidateEvent(today, models.SlotMorning),
			blackoutEvents: []models.Event{makeCandidateEvent(today, models.SlotEvening)},
			expected:       false,
		},
		"no_matching_blackout_events_different_day_same_slot": {
			candidateEvent: makeCandidateEvent(today, models.SlotMorning),
			blackoutEvents: []models.Event{makeCandidateEvent(tomorrow, models.SlotMorning)},
			expected:       false,
		},
		"no_matching_blackout_events_different_day_different_slot": {
			candidateEvent: makeCandidateEvent(today, models.SlotMorning),
			blackoutEvents: []models.Event{makeCandidateEvent(tomorrow, models.SlotEvening)},
			expected:       false,
		},
		"matching_blackout_events": {
			candidateEvent: makeCandidateEvent(today, models.SlotMorning),
			blackoutEvents: []models.Event{
				makeCandidateEvent(today, models.SlotMorning),
				makeCandidateEvent(today, models.SlotEvening),
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := isEventBlackedOut(test.candidateEvent, test.blackoutEvents)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestIsWeekend(t *testing.T) {
	tests := map[string]struct {
		day      time.Time
		expected bool
	}{
		"monday": {
			day:      time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		"tuesday": {
			day:      time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		"wednesday": {
			day:      time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		"thursday": {
			day:      time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		"friday": {
			day:      time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		"saturday": {
			day:      time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		"sunday": {
			day:      time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := isWeekend(test.day)
			require.Equal(t, test.expected, actual)
		})
	}

}
