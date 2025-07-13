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

func TestMakeCandidateEvents(t *testing.T) {
	tests := map[string]struct {
		input    *models.Input
		expected []models.Event
	}{
		"no_blackout_events": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Events: []models.Event{},
			},
			expected: []models.Event{
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotAfternoon},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotAfternoon},
			},
		},
		"all_blackout_events": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Events: []models.Event{
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotAfternoon},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotAfternoon},
				},
			},
			expected: []models.Event{},
		},
		"some_blackout_events": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				Events: []models.Event{
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC), Slot: models.SlotAfternoon},
					{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				},
			},
			expected: []models.Event{
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC), Slot: models.SlotMorning},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotEvening},
				{Type: models.EventTypeMatch, Date: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Slot: models.SlotAfternoon},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := makeCandidateEvents(test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestMakeUnscheduledEvents(t *testing.T) {
	tests := map[string]struct {
		input       *models.Input
		expected    []models.UnscheduledEvent
		expectedErr string
	}{
		"one_team_one_week": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Team: models.Team{
							Name:          "M3.5 40+",
							ScheduleGroup: models.TeamScheduleGroupDaytime,
						},
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			expected: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+",
						Type:  models.EventTypeMatch,
					},
					Constraints: models.Constraints{
						Required: []models.FilterConstraint{
							models.DayConstraint{
								NotBefore: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
								Before:    time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC),
							},
							models.SlotConstraint{TeamScheduleGroup: models.TeamScheduleGroupDaytime},
						},
					},
				},
			},
		},
		"one_team_two_weeks": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Team: models.Team{
							Name:          "M3.5 40+",
							ScheduleGroup: models.TeamScheduleGroupDaytime,
						},
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
							time.Date(2025, 7, 28, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			expected: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+",
						Type:  models.EventTypeMatch,
					},
					Constraints: models.Constraints{
						Required: []models.FilterConstraint{
							models.DayConstraint{
								NotBefore: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
								Before:    time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC),
							},
							models.SlotConstraint{TeamScheduleGroup: models.TeamScheduleGroupDaytime},
						},
					},
				},
				{
					Event: models.Event{
						Title: "M3.5 40+",
						Type:  models.EventTypeMatch,
					},
					Constraints: models.Constraints{
						Required: []models.FilterConstraint{
							models.DayConstraint{
								NotBefore: time.Date(2025, 7, 28, 0, 0, 0, 0, time.UTC),
								Before:    time.Date(2025, 8, 4, 0, 0, 0, 0, time.UTC),
							},
							models.SlotConstraint{TeamScheduleGroup: models.TeamScheduleGroupDaytime},
						},
					},
				},
			},
		},
		"two_teams_same_week": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Team: models.Team{
							Name:          "M3.5 40+ DT",
							ScheduleGroup: models.TeamScheduleGroupDaytime,
						},
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
					},
					{
						Team: models.Team{
							Name:          "W3.5 40+",
							ScheduleGroup: models.TeamScheduleGroupEvening,
						},
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			expected: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+ DT",
						Type:  models.EventTypeMatch,
					},
					Constraints: models.Constraints{
						Required: []models.FilterConstraint{
							models.DayConstraint{
								NotBefore: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
								Before:    time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC),
							},
							models.SlotConstraint{TeamScheduleGroup: models.TeamScheduleGroupDaytime},
						},
					},
				},
				{
					Event: models.Event{
						Title: "W3.5 40+",
						Type:  models.EventTypeMatch,
					},
					Constraints: models.Constraints{
						Required: []models.FilterConstraint{
							models.DayConstraint{
								NotBefore: time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
								Before:    time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC),
							},
							models.SlotConstraint{TeamScheduleGroup: models.TeamScheduleGroupEvening},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := makeUnscheduledEvents(test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestRemoveFromEvents(t *testing.T) {
	tests := map[string]struct {
		unscheduledEvents []models.UnscheduledEvent
		idxToRemove       int
		expected          []models.UnscheduledEvent
	}{
		"nothing_to_remove": {
			unscheduledEvents: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+ DT",
						Type:  models.EventTypeMatch,
						Slot:  models.SlotMorning,
					},
				},
			},
			idxToRemove: -1,
			expected: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+ DT",
						Type:  models.EventTypeMatch,
						Slot:  models.SlotMorning,
					},
				},
			},
		},
		"no_unscheduled_events": {
			unscheduledEvents: []models.UnscheduledEvent{},
			idxToRemove:       3,
			expected:          []models.UnscheduledEvent{},
		},
		"remove_second": {
			unscheduledEvents: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+ DT",
						Type:  models.EventTypeMatch,
						Slot:  models.SlotMorning,
					},
				},
				{
					Event: models.Event{
						Title: "W3.5 40+",
						Type:  models.EventTypeMatch,
						Slot:  models.SlotEvening,
					},
				},
			},
			idxToRemove: 1,
			expected: []models.UnscheduledEvent{
				{
					Event: models.Event{
						Title: "M3.5 40+ DT",
						Type:  models.EventTypeMatch,
						Slot:  models.SlotMorning,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := removeFromEvents[models.UnscheduledEvent](test.unscheduledEvents, test.idxToRemove)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestConstraining_Run(t *testing.T) {
	tests := map[string]struct {
		input    *models.Input
		expected *models.Schedule
	}{
		"no_teams": {
			input: &models.Input{
				Teams: nil,
			},
			expected: &models.Schedule{
				ScheduledEvents:   []models.Event{},
				UnscheduledEvents: []models.UnscheduledEvent{},
			},
		},
		"one_team_one_week": {
			input: &models.Input{
				Teams: []models.SchedulingTeam{
					{
						Team: models.Team{
							Name:          "M3.5 40+",
							ScheduleGroup: models.TeamScheduleGroupEvening,
						},
						Weeks: []time.Time{
							time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
						},
						DayPreferences: []time.Weekday{
							time.Wednesday,
						},
					},
				},
			},
			expected: &models.Schedule{
				ScheduledEvents: []models.Event{
					{
						Title: "M3.5 40+",
						Type:  models.EventTypeMatch,
						Date:  time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
						Slot:  models.SlotEvening,
					},
				},
				UnscheduledEvents: []models.UnscheduledEvent{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewConstraining(test.input)
			schedule, err := c.Run()
			require.NoError(t, err)
			require.Equal(t, test.expected, schedule)
		})
	}
}
