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
