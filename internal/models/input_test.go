package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFirstDayOfMatches(t *testing.T) {
	cases := map[string]struct {
		teams       []SchedulingTeam
		expected    *time.Time
		expectedErr error
	}{
		"empty": {
			teams:       nil,
			expected:    nil,
			expectedErr: nil,
		},
		"one_date_per_team": {
			teams: []SchedulingTeam{
				{Weeks: []time.Time{time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC)}},
				{Weeks: []time.Time{time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC)}},
				{Weeks: []time.Time{time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC)}},
				{Weeks: []time.Time{time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC)}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC)),
		},
		"multiple_dates_per_team": {
			teams: []SchedulingTeam{
				{Weeks: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC),
				}},
				{Weeks: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC),
				}},
				{Weeks: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC),
				}},
				{Weeks: []time.Time{
					time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC),
				}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC)),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			input := Input{Teams: test.teams}
			firstDayOfMatches := input.FirstDayOfMatches()
			require.Equal(t, test.expected, firstDayOfMatches)
		})
	}
}

func TestLastDayOfMatches(t *testing.T) {
	cases := map[string]struct {
		teams       []SchedulingTeam
		expected    *time.Time
		expectedErr error
	}{
		"empty": {
			teams:       nil,
			expected:    nil,
			expectedErr: nil,
		},
		"one_date_per_team": {
			teams: []SchedulingTeam{
				{Weeks: []time.Time{time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC)}},
				{Weeks: []time.Time{time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC)}},
				{Weeks: []time.Time{time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC)}},
				{Weeks: []time.Time{time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC)}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
		},
		"multiple_date_per_team": {
			teams: []SchedulingTeam{
				{Weeks: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC),
				}},
				{Weeks: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
				}},
				{Weeks: []time.Time{
					time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC),
				}},
				{Weeks: []time.Time{
					time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC),
				}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 11, 17, 0, 0, 0, 0, time.UTC)),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			input := Input{Teams: test.teams}
			lastDayOfMatches := input.LastDayOfMatches()
			require.Equal(t, test.expected, lastDayOfMatches)
		})
	}
}

func ptrTo[T any](v T) *T {
	return &v
}
