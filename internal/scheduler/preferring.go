package scheduler

import (
	"fmt"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

type Preferring struct {
	input *models.Input
}

func NewPreferring(input *models.Input) *Preferring {
	s := new(Preferring)
	s.input = input

	return s
}

func (p *Preferring) Run() (*models.Schedule, error) {
	// Track scheduled and unscheduled matches. Ideally, there should not be
	// any unscheduled matches at the end of this run.
	scheduledEvents := make([]models.Event, len(p.input.Events))
	copy(scheduledEvents, p.input.Events)
	fmt.Printf("input events: %v\n", p.input.Events)
	fmt.Printf("scheduled events: %v\n", scheduledEvents)
	unscheduledEvents := make([]models.UnscheduledEvent, 0)

	// Go week by week
	firstDayOfMatches := p.input.FirstDayOfMatches()
	lastDayOfMatches := p.input.LastDayOfMatches()
	fmt.Printf("first day of matches: [%s], last day of matches [%s]\n", firstDayOfMatches, lastDayOfMatches)
	for currentDay := *firstDayOfMatches; currentDay.Before(lastDayOfMatches.AddDate(0, 0, 1)); currentDay = currentDay.AddDate(0, 0, 7) {
		// Break out matches to be scheduled this week into
		// daytime and evening matches
		daytimeMatches := getMatchesForWeekAndScheduleGroup(p.input.Teams, currentDay, models.TeamScheduleGroupDaytime)
		eveningMatches := getMatchesForWeekAndScheduleGroup(p.input.Teams, currentDay, models.TeamScheduleGroupEvening)

		fmt.Printf("daytime matches for week [%s]: %v\n", currentDay, daytimeMatches)
		fmt.Printf("evening matches for week [%s]: %v\n", currentDay, eveningMatches)

		scheduledDaytimeMatches, unscheduledDaytimeMatches := scheduleMatches(daytimeMatches, models.TeamScheduleGroupDaytime, currentDay, scheduledEvents)
		scheduledEvents = append(scheduledEvents, scheduledDaytimeMatches...)
		unscheduledEvents = append(unscheduledEvents, unscheduledDaytimeMatches...)

		scheduledEveningMatches, unscheduledEveningMatches := scheduleMatches(eveningMatches, models.TeamScheduleGroupEvening, currentDay, scheduledEvents)
		scheduledEvents = append(scheduledEvents, scheduledEveningMatches...)
		unscheduledEvents = append(unscheduledEvents, unscheduledEveningMatches...)
	}

	s := models.Schedule{
		ScheduledEvents:   scheduledEvents,
		UnscheduledEvents: unscheduledEvents,
	}
	return &s, nil
}

func getMatchesForWeekAndScheduleGroup(teams []models.SchedulingTeam, week time.Time, scheduleGroup models.TeamScheduleGroup) []models.UnscheduledEvent {
	matches := make([]models.UnscheduledEvent, 0)
	for _, team := range teams {
		// If the team does not belong to the provided
		// schedule group (daytime, evening), skip to next
		// team.
		if team.ScheduleGroup != scheduleGroup {
			continue
		}

		// If the team does not have matches in the provided
		// week, skip to next team.
		teamHasMatchesThisWeek := false
		for _, currentWeek := range team.Weeks {
			if week == currentWeek {
				teamHasMatchesThisWeek = true
				break
			}
		}
		if !teamHasMatchesThisWeek {
			continue
		}

		// Schedule a match for the team this week.
		matches = append(matches, models.UnscheduledEvent{
			Event: models.Event{
				Title: team.DisplayName(),
				Type:  models.EventTypeMatch,
			},
			DayPreferences: team.DayPreferences,
		})
	}

	return matches
}

func scheduleMatches(matches []models.UnscheduledEvent, scheduleGroup models.TeamScheduleGroup, week time.Time, scheduledEvents []models.Event) ([]models.Event, []models.UnscheduledEvent) {
	scheduledMatches := make([]models.Event, len(scheduledEvents))
	copy(scheduledMatches, scheduledEvents)
	unscheduledMatches := make([]models.UnscheduledEvent, 0)

	// Randomize matches for the week. Go match by match.
	matches = randomizeSlice(matches)
	for _, match := range matches {
		// Find slot for match, considering slot availability, slot matching,
		// and match preference (look for 1st, then 2nd, etc.).
		scheduledMatch := findSlot(match, scheduleGroup, week, scheduledMatches, true)

		// If slot is not found, leave match unscheduled; it will be assigned
		// to an arbitrary available and matching slot later.
		if scheduledMatch == nil {
			unscheduledMatches = append(unscheduledMatches, match)
		} else {
			scheduledMatches = append(scheduledMatches, *scheduledMatch)
		}
	}

	// Loop over unscheduled matches and assign to arbitrary available and matching slot.
	finalUnscheduledMatches := make([]models.UnscheduledEvent, 0)
	for _, match := range unscheduledMatches {
		scheduledMatch := findSlot(match, scheduleGroup, week, scheduledMatches, false)
		if scheduledMatch == nil {
			finalUnscheduledMatches = append(finalUnscheduledMatches, match)
		} else {
			scheduledMatches = append(scheduledMatches, *scheduledMatch)
		}
	}

	return scheduledMatches, finalUnscheduledMatches
}

func findSlot(
	match models.UnscheduledEvent,
	scheduleGroup models.TeamScheduleGroup,
	week time.Time,
	scheduledEvents []models.Event,
	considerPreferred bool,
) *models.Event {
	candidateDates := make([]time.Time, 0)
	if considerPreferred && len(match.DayPreferences) > 0 {
		// If there are day preferences and we should consider them, create dates for
		// the week in order of day preferences.
		for _, preferredDay := range match.DayPreferences {
			preferredDayOffset := int(preferredDay) - 1
			if preferredDayOffset < 0 {
				preferredDayOffset += 7
			}

			candidateDates = append(candidateDates, week.AddDate(0, 0, preferredDayOffset))
		}
	} else {
		// Otherwise, randomize dates for the week so all matches don't fall on Mondays.
		for day := week; day.Before(week.AddDate(0, 0, 7)); day = day.AddDate(0, 0, 1) {
			candidateDates = append(candidateDates, day)
		}
		candidateDates = randomizeSlice(candidateDates)
	}

	for _, currentDay := range candidateDates {
		if !isWeekend(currentDay) {
			// If match is a daytime match, check if morning slot is available
			if scheduleGroup == models.TeamScheduleGroupDaytime && isSlotAvailable(scheduledEvents, currentDay, models.SlotMorning) {
				return &models.Event{
					Title: match.Title,
					Type:  match.Type,
					Slot:  models.SlotMorning,
					Date:  currentDay,
				}
			}

			// If match is an evening match, check if evening slot is available
			if scheduleGroup == models.TeamScheduleGroupEvening && isSlotAvailable(scheduledEvents, currentDay, models.SlotEvening) {
				return &models.Event{
					Title: match.Title,
					Type:  match.Type,
					Slot:  models.SlotEvening,
					Date:  currentDay,
				}
			}
		} else {
			// If match is a daytime match, it cannot be scheduled on a weekend
			if scheduleGroup == models.TeamScheduleGroupDaytime {
				continue
			}

			// Match is an evening match.

			// Check if morning slot is available.
			if isSlotAvailable(scheduledEvents, currentDay, models.SlotMorning) {
				return &models.Event{
					Title: match.Title,
					Type:  match.Type,
					Slot:  models.SlotMorning,
					Date:  currentDay,
				}
			}

			// Check if afternoon slot is available.
			if isSlotAvailable(scheduledEvents, currentDay, models.SlotAfternoon) {
				return &models.Event{
					Title: match.Title,
					Type:  match.Type,
					Slot:  models.SlotAfternoon,
					Date:  currentDay,
				}
			}

			// Check if evening slot is available.
			if isSlotAvailable(scheduledEvents, currentDay, models.SlotEvening) {
				return &models.Event{
					Title: match.Title,
					Type:  match.Type,
					Slot:  models.SlotEvening,
					Date:  currentDay,
				}
			}
		}
	}

	// No slots available for this match this week
	return nil
}

func isSlotAvailable(scheduledEvents []models.Event, currentDay time.Time, slot models.DaySlot) bool {
	for _, event := range scheduledEvents {
		if event.Date.Year() == currentDay.Year() &&
			event.Date.Month() == currentDay.Month() &&
			event.Date.Day() == currentDay.Day() &&
			event.Slot == slot {
			return false
		}
	}
	return true
}

//// Break out teams into daytime and evening teams
//daytimeTeams := filterTeamsBySchedulingType(p.input.Teams, "daytime")
//eveningTeams := filterTeamsBySchedulingType(p.input.Teams, "evening")
//
//// Group teams by week
//daytimeTeamsByWeek := mapTeamsByWeek(daytimeTeams)
//eveningTeamsByWeek := mapTeamsByWeek(eveningTeams)
//
//// Figure out first and last day of matches for schedule
//firstDayOfMatches := p.input.FirstDayOfMatches()
//lastDayOfMatches := p.input.LastDayOfMatches()

// Initialize schedule
//schedule := models.NewSchedule()

//teams := daytimeTeams
//teamsByWeek := daytimeTeamsByWeek
//
//var currentWeek string
//unscheduledTeams := make([]models.SchedulingTeam, 0)
//for currentDay := *firstDayOfMatches; !currentDay.After(*lastDayOfMatches); currentDay = currentDay.AddDate(0, 0, 1) {
//	newCurrentWeek := weekKey(currentDay)
//	if newCurrentWeek != currentWeek {
//		// We're changing to the next week...
//		// Verify that there are no unscheduled teams
//		if len(unscheduledTeams) > 0 {
//			return nil, fmt.Errorf("cannot schedule teams for week [%s] because there are unscheduled teams for the previous week", newCurrentWeek)
//		}
//
//		// Reset!
//		currentWeek = newCurrentWeek
//		unscheduledTeams = teamsByWeek[currentWeek]
//	}
//}

//return schedule, nil
//}

//	// Break out teams into daytime and evening teams
//	daytimeTeams := filterTeamsBySchedulingType(p.input.Teams, "daytime")
//	eveningTeams := filterTeamsBySchedulingType(p.input.Teams, "evening")
//
//	// Group teams by week
//	daytimeTeamsByWeek := mapTeamsByWeek(daytimeTeams)
//	eveningTeamsByWeek := mapTeamsByWeek(eveningTeams)
//
//	// Figure out first and last day of matches for schedule
//	firstDayOfMatches, err := p.input.FirstDayOfMatches()
//	if err != nil {
//		return nil, fmt.Errorf("failed to compute first day of matches from input schedule: %w", err)
//	}
//	lastDayOfMatches, err := p.input.LastDayOfMatches()
//	if err != nil {
//		return nil, fmt.Errorf("failed to compute last day of matches from input schedule: %w", err)
//	}
//
//	// Initialize schedule
//	schedule := models.NewSchedule(*firstDayOfMatches, *lastDayOfMatches)
//
//	// First pass:
//	// Loop over each day in schedule, keeping track of current week (by start date); for each day:
//	// - check if there's capacity to schedule matches on that day
//	// - if so, filter teams that prefer that day of the week
//	// - randomly pick a team from the list and assign it to that day
//	for currentDay := *firstDayOfMatches; !currentDay.After(*lastDayOfMatches); currentDay = currentDay.AddDate(0, 0, 1) {
//		currentWeek := weekKey(currentDay)
//		//fmt.Println("current week:", currentWeek)
//		//fmt.Println("current day:", currentDay.Format("01/02/2006"), currentDay.Weekday())
//
//		currentDaySchedule := schedule.ForDay(currentDay)
//		if !currentDaySchedule.HasCapacity() {
//			// Day is full; cannot schedule
//			continue
//		}
//
//		if isBlackoutDate(currentDay, p.input.BlackoutSlots) {
//			//fmt.Println(currentDay.Format("20060102"), "is a blackout day, skipping it")
//			continue
//		}
//
//		// Daytime
//		if currentDaySchedule.HasDaytimeCapacity() {
//			candidateTeams := daytimeTeamsByWeek[currentWeek]
//			candidateTeamsForDay, err := teamsThatPreferDay(candidateTeams, currentDay.Weekday())
//			if err != nil {
//				return nil, fmt.Errorf("cannot figure out teams that prefer to play on [%s]: %w", currentDay.Weekday().String(), err)
//			}
//
//			if len(candidateTeamsForDay) > 0 {
//				chosenTeamIdx := rand.Intn(len(candidateTeamsForDay))
//				chosenTeam := candidateTeams[chosenTeamIdx]
//
//				// Remove chosen team from teams by week, so it's not chosen again
//				daytimeTeamsByWeek[currentWeek] = removeTeam(candidateTeams, chosenTeam)
//
//				// Assign chosen team to schedule
//				currentDaySchedule.DaytimeTeam = &chosenTeam
//
//				//fmt.Printf("- Assigned [%s] to daytime\n", chosenTeam.Name)
//			}
//		}
//
//		// Evening
//		if currentDaySchedule.HasEveningCapacity() {
//			candidateTeams := eveningTeamsByWeek[currentWeek]
//			candidateTeamsForDay, err := teamsThatPreferDay(candidateTeams, currentDay.Weekday())
//			if err != nil {
//				return nil, fmt.Errorf("cannot figure out teams that prefer to play on [%s]: %w", currentDay.Weekday().String(), err)
//			}
//
//			if len(candidateTeamsForDay) > 0 {
//				chosenTeamIdx := rand.Intn(len(candidateTeamsForDay))
//				chosenTeam := candidateTeams[chosenTeamIdx]
//
//				//fmt.Println("chosen team:", chosenTeam.Name)
//
//				// Remove chosen team from teams by week, so it's not chosen again
//				eveningTeamsByWeek[currentWeek] = removeTeam(candidateTeams, chosenTeam)
//				//fmt.Println(eveningTeamsByWeek[currentWeek])
//
//				// Assign chosen team to schedule
//				currentDaySchedule.EveningTeam = &chosenTeam
//
//				//fmt.Printf("- Assigned [%s] to evening\n", chosenTeam.Name)
//			}
//		}
//	}
//
//	// Second pass: set schedule eagerly for unassigned teams
//	setScheduleEagerly(firstDayOfMatches, lastDayOfMatches, schedule, daytimeTeamsByWeek, eveningTeamsByWeek, p.input.BlackoutSlots)
//
//	return schedule, nil
//}

func randomizeSlice[T any](slice []T) []T {
	randomizedIndices := randSliceOfIntegers(len(slice))
	output := make([]T, len(slice))
	for idx, randomizedIdx := range randomizedIndices {
		output[randomizedIdx] = slice[idx]
	}
	return output
}
