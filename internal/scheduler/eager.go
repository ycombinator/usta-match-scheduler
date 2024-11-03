package scheduler

import (
	"fmt"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"math/rand"
)

type Eager struct {
	input models.Input
}

func NewEager(input models.Input) (*Eager, error) {
	s := new(Eager)
	s.input = input

	return s, nil
}

func (s *Eager) Run() (*models.Schedule, error) {
	// Break out teams into daytime and evening teams
	daytimeTeams := filterTeamsBySchedulingType(s.input.Teams, "daytime")
	eveningTeams := filterTeamsBySchedulingType(s.input.Teams, "evening")

	// Group teams by week
	daytimeTeamsByWeek := mapTeamsByWeek(daytimeTeams)
	eveningTeamsByWeek := mapTeamsByWeek(eveningTeams)

	// Figure out first and last day of matches for schedule
	firstDayOfMatches, err := s.input.FirstDayOfMatches()
	if err != nil {
		return nil, fmt.Errorf("failed to compute first day of matches from input schedule: %w", err)
	}
	lastDayOfMatches, err := s.input.LastDayOfMatches()
	if err != nil {
		return nil, fmt.Errorf("failed to compute last day of matches from input schedule: %w", err)
	}

	// Initialize schedule
	schedule := models.NewSchedule(*firstDayOfMatches, *lastDayOfMatches)

	// Loop over each day in schedule, keeping track of the current week (by start date); for each day:
	// - check if there's capacity to schedule matches on that day
	// - if so, filter teams to those that have matches that week
	// - randomly pick a team from the list and assign it to that day
	for currentDay := *firstDayOfMatches; !currentDay.After(*lastDayOfMatches); currentDay = currentDay.AddDate(0, 0, 1) {
		currentWeek := weekKey(currentDay)
		//fmt.Println(currentDay.Format("01/02/2006"), currentDay.Weekday())

		currentDaySchedule := schedule.ForDay(currentDay)
		if !currentDaySchedule.HasCapacity() {
			// Day is full; cannot schedule
			continue
		}

		// Daytime
		if currentDaySchedule.HasDaytimeCapacity() {
			candidateTeams := daytimeTeamsByWeek[currentWeek]
			if len(candidateTeams) > 0 {
				chosenTeamIdx := rand.Intn(len(candidateTeams))
				chosenTeam := candidateTeams[chosenTeamIdx]

				// Remove chosen team from teams by week, so it's not chosen again
				newTeams := append(candidateTeams[:chosenTeamIdx], candidateTeams[chosenTeamIdx+1:]...)
				daytimeTeamsByWeek[currentWeek] = newTeams

				// Assign chosen team to schedule
				currentDaySchedule.DaytimeTeam = &chosenTeam

				//fmt.Printf("- Assigned [%s] to daytime\n", chosenTeam.Title)
			}
		}

		// Evening
		if currentDaySchedule.HasEveningCapacity() {
			candidateTeams := eveningTeamsByWeek[currentWeek]
			if len(candidateTeams) > 0 {
				chosenTeamIdx := rand.Intn(len(candidateTeams))
				chosenTeam := candidateTeams[chosenTeamIdx]

				// Remove chosen team from teams by week, so it's not chosen again
				newTeams := append(candidateTeams[:chosenTeamIdx], candidateTeams[chosenTeamIdx+1:]...)
				eveningTeamsByWeek[currentWeek] = newTeams

				// Assign chosen team to schedule
				currentDaySchedule.EveningTeam = &chosenTeam

				//fmt.Printf("- Assigned [%s] to evening\n", chosenTeam.Title)
			}
		}
	}

	return schedule, nil
}
