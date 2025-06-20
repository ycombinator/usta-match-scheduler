package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFirstDayOfMatches(t *testing.T) {
	cases := map[string]struct {
		teams       map[string]SchedulingTeam
		expected    *time.Time
		expectedErr error
	}{
		"empty": {
			teams:       nil,
			expected:    nil,
			expectedErr: nil,
		},
		"one_date_per_team": {
			teams: map[string]SchedulingTeam{
				"w3.5": {Weeks: []string{"20241028"}},
				"m3.5": {Weeks: []string{"20241021"}},
				"m4.5": {Weeks: []string{"20241021"}},
				"w2.5": {Weeks: []string{"20241104"}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 10, 21, 0, 0, 0, 0, time.UTC)),
		},
		"multiple_date_per_team": {
			teams: map[string]SchedulingTeam{
				"w3.5": {Weeks: []string{"20241021", "20241028"}},
				"m3.5": {Weeks: []string{"20241021", "20241014"}},
				"m4.5": {Weeks: []string{"20241021", "20241104"}},
				"w2.5": {Weeks: []string{"20241104"}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC)),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			input := Input{Teams: test.teams}
			firstDayOfMatches, err := input.FirstDayOfMatches()

			require.Equal(t, test.expected, firstDayOfMatches)
			require.Equal(t, test.expectedErr, err)
		})
	}
}

func TestLastDayOfMatches(t *testing.T) {
	cases := map[string]struct {
		teams       map[string]SchedulingTeam
		expected    *time.Time
		expectedErr error
	}{
		"empty": {
			teams:       nil,
			expected:    nil,
			expectedErr: nil,
		},
		"one_date_per_team": {
			teams: map[string]SchedulingTeam{
				"w3.5": {Weeks: []string{"20241028"}},
				"m3.5": {Weeks: []string{"20241104"}},
				"m4.5": {Weeks: []string{"20241021"}},
				"w2.5": {Weeks: []string{"20241028"}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
		},
		"multiple_date_per_team": {
			teams: map[string]SchedulingTeam{
				"w3.5": {Weeks: []string{"20241021", "20241028"}},
				"m3.5": {Weeks: []string{"20241021", "20241111"}},
				"m4.5": {Weeks: []string{"20241021", "20241104"}},
				"w2.5": {Weeks: []string{"20241104"}},
			},
			expected: ptrTo[time.Time](time.Date(2024, 11, 17, 0, 0, 0, 0, time.UTC)),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			input := Input{Teams: test.teams}
			lastDayOfMatches, err := input.LastDayOfMatches()

			require.Equal(t, test.expected, lastDayOfMatches)
			require.Equal(t, test.expectedErr, err)
		})
	}
}

func ptrTo[T any](v T) *T {
	return &v
}
