package scheduler

import (
	"fmt"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"math/rand"
)

type Preferring struct {
	input models.Input
}

func NewPreferring(input models.Input) (*Preferring, error) {
	s := new(Preferring)
	s.input = input

	return s, nil
}

func (p *Preferring) Run() (*models.Schedule, error) {
	// Break out teams into daytime and evening teams
	daytimeTeams := filterTeamsBySchedulingType(p.input.Teams, "daytime")
	eveningTeams := filterTeamsBySchedulingType(p.input.Teams, "evening")

	// Group teams by week
	daytimeTeamsByWeek := mapTeamsByWeek(daytimeTeams)
	eveningTeamsByWeek := mapTeamsByWeek(eveningTeams)

	// Figure out first and last day of matches for schedule
	firstDayOfMatches, err := p.input.FirstDayOfMatches()
	if err != nil {
		return nil, fmt.Errorf("failed to compute first day of matches from input schedule: %w", err)
	}
	lastDayOfMatches, err := p.input.LastDayOfMatches()
	if err != nil {
		return nil, fmt.Errorf("failed to compute last day of matches from input schedule: %w", err)
	}

	// Initialize schedule
	schedule := models.NewSchedule(*firstDayOfMatches, *lastDayOfMatches)

	// First pass:
	// Loop over each day in schedule, keeping track of current week (by start date); for each day:
	// - check if there's capacity to schedule matches on that day
	// - if so, filter teams that prefer that day of the week
	// - randomly pick a team from the list and assign it to that day
	for currentDay := *firstDayOfMatches; !currentDay.After(*lastDayOfMatches); currentDay = currentDay.AddDate(0, 0, 1) {
		currentWeek := weekKey(currentDay)
		//fmt.Println("current week:", currentWeek)
		//fmt.Println("current day:", currentDay.Format("01/02/2006"), currentDay.Weekday())

		currentDaySchedule := schedule.ForDay(currentDay)
		if !currentDaySchedule.HasCapacity() {
			// Day is full; cannot schedule
			continue
		}

		if isBlackoutDate(currentDay, p.input.BlackoutDates) {
			//fmt.Println(currentDay.Format("20060102"), "is a blackout day, skipping it")
			continue
		}

		// Daytime
		if currentDaySchedule.HasDaytimeCapacity() {
			candidateTeams := daytimeTeamsByWeek[currentWeek]
			candidateTeamsForDay, err := teamsThatPreferDay(candidateTeams, currentDay.Weekday())
			if err != nil {
				return nil, fmt.Errorf("cannot figure out teams that prefer to play on [%s]: %w", currentDay.Weekday().String(), err)
			}

			if len(candidateTeamsForDay) > 0 {
				chosenTeamIdx := rand.Intn(len(candidateTeamsForDay))
				chosenTeam := candidateTeams[chosenTeamIdx]

				// Remove chosen team from teams by week, so it's not chosen again
				daytimeTeamsByWeek[currentWeek] = removeTeam(candidateTeams, chosenTeam)

				// Assign chosen team to schedule
				currentDaySchedule.DaytimeTeam = &chosenTeam

				//fmt.Printf("- Assigned [%s] to daytime\n", chosenTeam.Title)
			}
		}

		// Evening
		if currentDaySchedule.HasEveningCapacity() {
			candidateTeams := eveningTeamsByWeek[currentWeek]
			candidateTeamsForDay, err := teamsThatPreferDay(candidateTeams, currentDay.Weekday())
			if err != nil {
				return nil, fmt.Errorf("cannot figure out teams that prefer to play on [%s]: %w", currentDay.Weekday().String(), err)
			}

			if len(candidateTeamsForDay) > 0 {
				chosenTeamIdx := rand.Intn(len(candidateTeamsForDay))
				chosenTeam := candidateTeams[chosenTeamIdx]

				//fmt.Println("chosen team:", chosenTeam.Title)

				// Remove chosen team from teams by week, so it's not chosen again
				eveningTeamsByWeek[currentWeek] = removeTeam(candidateTeams, chosenTeam)
				//fmt.Println(eveningTeamsByWeek[currentWeek])

				// Assign chosen team to schedule
				currentDaySchedule.EveningTeam = &chosenTeam

				//fmt.Printf("- Assigned [%s] to evening\n", chosenTeam.Title)
			}
		}
	}

	// Second pass: set schedule eagerly for unassigned teams
	setScheduleEagerly(firstDayOfMatches, lastDayOfMatches, schedule, daytimeTeamsByWeek, eveningTeamsByWeek, p.input.BlackoutDates)

	return schedule, nil
}
