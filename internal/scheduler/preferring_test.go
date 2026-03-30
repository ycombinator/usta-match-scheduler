package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

func TestGetMatchesForWeekAndScheduleGroup(t *testing.T) {
	week := time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC)
	otherWeek := time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC)

	t.Run("regular_match", func(t *testing.T) {
		teams := []models.SchedulingTeam{
			{
				Team: models.Team{
					Type:          models.TeamTypeMixed,
					MinAge:        40,
					Gender:        models.GenderMale,
					Level:         "3.5",
					Captain:       "Smith",
					ScheduleGroup: models.TeamScheduleGroupEvening,
				},
				Weeks: []time.Time{week},
			},
		}

		matches := getMatchesForWeekAndScheduleGroup(teams, week, models.TeamScheduleGroupEvening)
		require.Len(t, matches, 1)
		require.False(t, matches[0].IsPlayoff)
	})

	t.Run("playoff_match", func(t *testing.T) {
		teams := []models.SchedulingTeam{
			{
				Team: models.Team{
					Type:          models.TeamTypeMixed,
					MinAge:        40,
					Gender:        models.GenderMale,
					Level:         "3.5",
					Captain:       "Smith",
					ScheduleGroup: models.TeamScheduleGroupEvening,
				},
				Weeks:        []time.Time{week},
				PlayoffWeeks: map[time.Time]bool{week: true},
			},
		}

		matches := getMatchesForWeekAndScheduleGroup(teams, week, models.TeamScheduleGroupEvening)
		require.Len(t, matches, 1)
		require.True(t, matches[0].IsPlayoff)
	})

	t.Run("playoff_week_nil_map", func(t *testing.T) {
		teams := []models.SchedulingTeam{
			{
				Team: models.Team{
					Type:          models.TeamTypeMixed,
					MinAge:        40,
					Gender:        models.GenderMale,
					Level:         "3.5",
					Captain:       "Smith",
					ScheduleGroup: models.TeamScheduleGroupEvening,
				},
				Weeks:        []time.Time{week},
				PlayoffWeeks: nil,
			},
		}

		matches := getMatchesForWeekAndScheduleGroup(teams, week, models.TeamScheduleGroupEvening)
		require.Len(t, matches, 1)
		require.False(t, matches[0].IsPlayoff, "nil PlayoffWeeks should not mark match as playoff")
	})

	t.Run("wrong_week", func(t *testing.T) {
		teams := []models.SchedulingTeam{
			{
				Team: models.Team{
					Type:          models.TeamTypeMixed,
					MinAge:        40,
					Gender:        models.GenderMale,
					Level:         "3.5",
					Captain:       "Smith",
					ScheduleGroup: models.TeamScheduleGroupEvening,
				},
				Weeks:        []time.Time{otherWeek},
				PlayoffWeeks: map[time.Time]bool{otherWeek: true},
			},
		}

		matches := getMatchesForWeekAndScheduleGroup(teams, week, models.TeamScheduleGroupEvening)
		require.Empty(t, matches)
	})

	t.Run("wrong_schedule_group", func(t *testing.T) {
		teams := []models.SchedulingTeam{
			{
				Team: models.Team{
					Type:          models.TeamTypeMixed,
					MinAge:        40,
					Gender:        models.GenderMale,
					Level:         "3.5",
					Captain:       "Smith",
					ScheduleGroup: models.TeamScheduleGroupDaytime,
				},
				Weeks:        []time.Time{week},
				PlayoffWeeks: map[time.Time]bool{week: true},
			},
		}

		matches := getMatchesForWeekAndScheduleGroup(teams, week, models.TeamScheduleGroupEvening)
		require.Empty(t, matches)
	})

	t.Run("mixed_regular_and_playoff_teams", func(t *testing.T) {
		teams := []models.SchedulingTeam{
			{
				Team: models.Team{
					Type:          models.TeamTypeMixed,
					MinAge:        40,
					Gender:        models.GenderMale,
					Level:         "3.5",
					Captain:       "Smith",
					ScheduleGroup: models.TeamScheduleGroupEvening,
				},
				Weeks: []time.Time{week},
			},
			{
				Team: models.Team{
					Type:          models.TeamTypeAdult,
					MinAge:        40,
					Gender:        models.GenderFemale,
					Level:         "3.5",
					Captain:       "Jones",
					ScheduleGroup: models.TeamScheduleGroupEvening,
				},
				Weeks:        []time.Time{week},
				PlayoffWeeks: map[time.Time]bool{week: true},
			},
		}

		matches := getMatchesForWeekAndScheduleGroup(teams, week, models.TeamScheduleGroupEvening)
		require.Len(t, matches, 2)
		require.False(t, matches[0].IsPlayoff, "first team is regular")
		require.True(t, matches[1].IsPlayoff, "second team is playoff")
	})
}

func TestFindSlot_EveningPlayoff(t *testing.T) {
	// Week starts on Monday 2025-07-14
	week := time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC)

	eveningPlayoffMatch := models.UnscheduledEvent{
		Event: models.Event{
			Title: "40+ Mixed 3.5 (Smith)",
			Type:  models.EventTypeMatch,
		},
		IsPlayoff: true,
	}

	regularEveningMatch := models.UnscheduledEvent{
		Event: models.Event{
			Title: "40+ Mixed 3.5 (Jones)",
			Type:  models.EventTypeMatch,
		},
		IsPlayoff: false,
	}

	t.Run("evening_playoff_skips_mon_tue_wed", func(t *testing.T) {
		// With Thursday preference, playoff should land on Thursday
		eveningPlayoffWithPref := eveningPlayoffMatch
		eveningPlayoffWithPref.DayPreferences = []time.Weekday{
			time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		}

		result := findSlot(eveningPlayoffWithPref, models.TeamScheduleGroupEvening, week, nil, true)
		require.NotNil(t, result)
		require.Equal(t, time.Thursday, result.Date.Weekday(),
			"evening playoff should skip Mon/Tue/Wed and land on Thursday")
		require.Equal(t, models.SlotEvening, result.Slot)
	})

	t.Run("evening_playoff_only_mon_tue_wed_prefs", func(t *testing.T) {
		// If only Mon-Wed preferences, preferred pass should fail
		eveningPlayoffMTW := eveningPlayoffMatch
		eveningPlayoffMTW.DayPreferences = []time.Weekday{
			time.Monday, time.Tuesday, time.Wednesday,
		}

		result := findSlot(eveningPlayoffMTW, models.TeamScheduleGroupEvening, week, nil, true)
		require.Nil(t, result,
			"evening playoff with only Mon-Wed preferences should not find a slot in preferred pass")
	})

	t.Run("evening_playoff_falls_back_to_thu_sun", func(t *testing.T) {
		// Without preferences (non-preferred pass), should find Thu-Sun
		result := findSlot(eveningPlayoffMatch, models.TeamScheduleGroupEvening, week, nil, false)
		require.NotNil(t, result)
		day := result.Date.Weekday()
		require.True(t,
			day == time.Thursday || day == time.Friday || day == time.Saturday || day == time.Sunday,
			"evening playoff should only be scheduled Thu-Sun, got %s", day)
	})

	t.Run("regular_evening_allows_all_weekdays", func(t *testing.T) {
		// Regular evening match with Monday preference should work
		regularWithPref := regularEveningMatch
		regularWithPref.DayPreferences = []time.Weekday{time.Monday}

		result := findSlot(regularWithPref, models.TeamScheduleGroupEvening, week, nil, true)
		require.NotNil(t, result)
		require.Equal(t, time.Monday, result.Date.Weekday(),
			"regular evening match should be allowed on Monday")
	})

	t.Run("evening_playoff_all_thu_sun_full", func(t *testing.T) {
		// Fill all Thu-Sun slots
		thursday := time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC)
		friday := time.Date(2025, 7, 18, 0, 0, 0, 0, time.UTC)
		saturday := time.Date(2025, 7, 19, 0, 0, 0, 0, time.UTC)
		sunday := time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC)

		scheduledEvents := []models.Event{
			{Date: thursday, Slot: models.SlotEvening},
			{Date: friday, Slot: models.SlotEvening},
			// Weekend slots: morning, afternoon, evening
			{Date: saturday, Slot: models.SlotMorning},
			{Date: saturday, Slot: models.SlotAfternoon},
			{Date: saturday, Slot: models.SlotEvening},
			{Date: sunday, Slot: models.SlotMorning},
			{Date: sunday, Slot: models.SlotAfternoon},
			{Date: sunday, Slot: models.SlotEvening},
		}

		result := findSlot(eveningPlayoffMatch, models.TeamScheduleGroupEvening, week, scheduledEvents, false)
		require.Nil(t, result,
			"should return nil when all Thu-Sun slots are full for evening playoff")
	})
}

func TestFindSlot_DaytimePlayoff(t *testing.T) {
	week := time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC)

	daytimePlayoffMatch := models.UnscheduledEvent{
		Event: models.Event{
			Title: "40+ Women's 3.5 Daytime (Smith)",
			Type:  models.EventTypeMatch,
		},
		IsPlayoff: true,
	}

	t.Run("daytime_playoff_allows_weekdays", func(t *testing.T) {
		// Daytime playoff should be allowed Mon-Fri (same as regular daytime)
		daytimeWithPref := daytimePlayoffMatch
		daytimeWithPref.DayPreferences = []time.Weekday{time.Monday}

		result := findSlot(daytimeWithPref, models.TeamScheduleGroupDaytime, week, nil, true)
		require.NotNil(t, result)
		require.Equal(t, time.Monday, result.Date.Weekday(),
			"daytime playoff should be allowed on Monday")
		require.Equal(t, models.SlotMorning, result.Slot)
	})

	t.Run("daytime_playoff_not_on_weekend", func(t *testing.T) {
		// Daytime playoff should not be scheduled on weekends (same as regular)
		daytimeWeekendPref := daytimePlayoffMatch
		daytimeWeekendPref.DayPreferences = []time.Weekday{time.Saturday, time.Sunday}

		result := findSlot(daytimeWeekendPref, models.TeamScheduleGroupDaytime, week, nil, true)
		require.Nil(t, result,
			"daytime playoff should not be scheduled on weekend")
	})
}

func TestPreferring_Run_EveningPlayoff(t *testing.T) {
	// Week starts Monday 2025-07-14
	week := time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC)

	t.Run("evening_playoff_scheduled_thu_sun", func(t *testing.T) {
		input := &models.Input{
			Teams: []models.SchedulingTeam{
				{
					Team: models.Team{
						Type:          models.TeamTypeMixed,
						MinAge:        40,
						Gender:        models.GenderMale,
						Level:         "3.5",
						Captain:       "Smith",
						ScheduleGroup: models.TeamScheduleGroupEvening,
					},
					Weeks:        []time.Time{week},
					PlayoffWeeks: map[time.Time]bool{week: true},
					DayPreferences: []time.Weekday{
						time.Thursday,
					},
				},
			},
		}

		p := NewPreferring(input)
		schedule, err := p.Run()
		require.NoError(t, err)
		require.Len(t, schedule.ScheduledEvents, 1)
		require.Empty(t, schedule.UnscheduledEvents)

		event := schedule.ScheduledEvents[0]
		day := event.Date.Weekday()
		require.True(t,
			day == time.Thursday || day == time.Friday || day == time.Saturday || day == time.Sunday,
			"evening playoff should be scheduled Thu-Sun, got %s", day)
	})

	t.Run("regular_and_playoff_same_week", func(t *testing.T) {
		input := &models.Input{
			Teams: []models.SchedulingTeam{
				{
					Team: models.Team{
						Type:          models.TeamTypeMixed,
						MinAge:        40,
						Gender:        models.GenderMale,
						Level:         "3.5",
						Captain:       "Smith",
						ScheduleGroup: models.TeamScheduleGroupEvening,
					},
					Weeks: []time.Time{week},
					DayPreferences: []time.Weekday{
						time.Tuesday,
					},
				},
				{
					Team: models.Team{
						Type:          models.TeamTypeAdult,
						MinAge:        40,
						Gender:        models.GenderFemale,
						Level:         "3.5",
						Captain:       "Jones",
						ScheduleGroup: models.TeamScheduleGroupEvening,
					},
					Weeks:        []time.Time{week},
					PlayoffWeeks: map[time.Time]bool{week: true},
					DayPreferences: []time.Weekday{
						time.Friday,
					},
				},
			},
		}

		p := NewPreferring(input)
		schedule, err := p.Run()
		require.NoError(t, err)
		require.Len(t, schedule.ScheduledEvents, 2)
		require.Empty(t, schedule.UnscheduledEvents)

		// Find regular and playoff events
		for _, event := range schedule.ScheduledEvents {
			if event.Title == "40+ Women's 3.5 (Jones)" {
				// Playoff match
				day := event.Date.Weekday()
				require.True(t,
					day == time.Thursday || day == time.Friday || day == time.Saturday || day == time.Sunday,
					"playoff match should be Thu-Sun, got %s", day)
			}
			// Regular match can be any day - no restriction
		}
	})

	t.Run("daytime_playoff_any_weekday", func(t *testing.T) {
		input := &models.Input{
			Teams: []models.SchedulingTeam{
				{
					Team: models.Team{
						Type:          models.TeamTypeAdult,
						MinAge:        40,
						Gender:        models.GenderFemale,
						Level:         "3.5",
						Captain:       "Smith",
						ScheduleGroup: models.TeamScheduleGroupDaytime,
					},
					Weeks:        []time.Time{week},
					PlayoffWeeks: map[time.Time]bool{week: true},
					DayPreferences: []time.Weekday{
						time.Monday,
					},
				},
			},
		}

		p := NewPreferring(input)
		schedule, err := p.Run()
		require.NoError(t, err)
		require.Len(t, schedule.ScheduledEvents, 1)
		require.Empty(t, schedule.UnscheduledEvents)

		event := schedule.ScheduledEvents[0]
		require.Equal(t, time.Monday, event.Date.Weekday(),
			"daytime playoff should be allowed on Monday")
		require.Equal(t, models.SlotMorning, event.Slot)
	})
}
