package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnscheduledEvent_Match(t *testing.T) {
	friday := time.Date(2025, 7, 11, 0, 0, 0, 0, time.UTC)
	saturday := time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC)
	sunday := time.Date(2025, 7, 13, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		unscheduledEvent UnscheduledEvent
		candidateEvent   Event
		expected         float64
	}{
		"date_too_early": {
			unscheduledEvent: UnscheduledEvent{
				Constraints: Constraints{
					Required: []FilterConstraint{
						DayConstraint{
							NotBefore: saturday,
							Before:    sunday,
						},
					},
				},
			},
			candidateEvent: Event{
				Date: friday,
			},
			expected: 0,
		},
		"date_too_late": {
			unscheduledEvent: UnscheduledEvent{
				Constraints: Constraints{
					Required: []FilterConstraint{
						DayConstraint{
							NotBefore: friday,
							Before:    sunday,
						},
					},
				},
			},
			candidateEvent: Event{
				Date: sunday,
			},
			expected: 0,
		},
		"slot_weekday_mismatch": {
			unscheduledEvent: UnscheduledEvent{
				Constraints: Constraints{
					Required: []FilterConstraint{
						SlotConstraint{
							TeamScheduleGroup: TeamScheduleGroupEvening,
						},
					},
				},
			},
			candidateEvent: Event{
				Slot: SlotAfternoon,
				Date: friday,
			},
			expected: 0,
		},
		"slot_weekend_mismatch": {
			unscheduledEvent: UnscheduledEvent{
				Constraints: Constraints{
					Required: []FilterConstraint{
						SlotConstraint{
							TeamScheduleGroup: TeamScheduleGroupDaytime,
						},
					},
				},
			},
			candidateEvent: Event{
				Slot: SlotAfternoon,
				Date: friday,
			},
			expected: 0,
		},
		"all_match": {
			unscheduledEvent: UnscheduledEvent{
				Constraints: Constraints{
					Required: []FilterConstraint{
						DayConstraint{
							NotBefore: friday,
							Before:    sunday,
						},
						SlotConstraint{
							TeamScheduleGroup: TeamScheduleGroupDaytime,
						},
					},
				},
			},
			candidateEvent: Event{
				Slot: SlotMorning,
				Date: friday,
			},
			expected: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			score := test.unscheduledEvent.Match(test.candidateEvent)
			require.Equal(t, test.expected, score)
		})
	}
}
