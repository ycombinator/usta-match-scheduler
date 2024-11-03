package scheduler

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWeekKey(t *testing.T) {
	cases := map[string]struct {
		input    time.Time
		expected string
	}{
		"monday": {
			input:    time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC),
			expected: "20241028",
		},
		"tuesday": {
			input:    time.Date(2024, 10, 29, 0, 0, 0, 0, time.UTC),
			expected: "20241028",
		},
		"sunday": {
			input:    time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC),
			expected: "20241028",
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actual := weekKey(test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestWeekdayFromStr(t *testing.T) {
	cases := map[string]struct {
		expectedWeekday time.Weekday
		expectedErr     string
	}{
		"sun":      {time.Sunday, ""},
		"Sun":      {time.Sunday, ""},
		"sunday":   {time.Sunday, ""},
		" Sunday ": {time.Sunday, ""},
		"mon":      {time.Monday, ""},
		"tue":      {time.Tuesday, ""},
		"wed":      {time.Wednesday, ""},
		"thu":      {time.Thursday, ""},
		"fri":      {time.Friday, ""},
		"sat":      {time.Saturday, ""},
		"t":        {0, "cannot parse given day [t] as weekday"},
		"doomsday": {0, "cannot parse given day [doomsday] as weekday"},
	}

	for input, test := range cases {
		t.Run(input, func(t *testing.T) {
			weekday, err := weekdayFromStr(input)
			if test.expectedErr == "" {
				require.Equal(t, test.expectedWeekday, weekday)
				require.NoError(t, err)
			} else {
				require.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}
