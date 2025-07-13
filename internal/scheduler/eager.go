package scheduler

import (
	"github.com/ycombinator/usta-match-scheduler/internal/models"
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
	//// Break out teams into daytime and evening teams
	//daytimeTeams := filterTeamsBySchedulingType(s.input.Teams, models.TeamScheduleGroupDaytime)
	//eveningTeams := filterTeamsBySchedulingType(s.input.Teams, models.TeamScheduleGroupDaytime)
	//
	//// Group teams by week
	//daytimeTeamsByWeek := mapTeamsByWeek(daytimeTeams)
	//eveningTeamsByWeek := mapTeamsByWeek(eveningTeams)

	// Figure out first and last day of matches for schedule
	firstDayOfMatches := s.input.FirstDayOfMatches()
	lastDayOfMatches := s.input.LastDayOfMatches()

	// Initialize schedule
	schedule := models.NewSchedule()

	for currentDay := *firstDayOfMatches; !currentDay.After(*lastDayOfMatches); currentDay = currentDay.AddDate(0, 0, 7) {
		//currentWeek := weekKey(currentDay)
		//teams := daytimeTeamsByWeek[currentWeek]

	}

	//setScheduleEagerly(firstDayOfMatches, lastDayOfMatches, schedule, daytimeTeamsByWeek, eveningTeamsByWeek, s.input.BlackoutSlots)

	return schedule, nil
}

//func filterTeamsByDayPreference(teams []models.SchedulingTeam, day time.Weekday) []models.SchedulingTeam {
//	filteredTeams := make([]models.SchedulingTeam, 0)
//	for _, team := range teams {
//		if team.DayPreferences == day {
//			filteredTeams = append(filteredTeams, team)
//		}
//	}
//	return filteredTeams
//}

//func setScheduleEagerly(
//	firstDayOfMatches *time.Time, lastDayOfMatches *time.Time,
//	schedule *models.Schedule,
//	daytimeTeamsByWeek map[string][]models.Team, eveningTeamsByWeek map[string][]models.Team,
//	blackoutDates []string,
//) {
//	// Loop over each day in schedule, keeping track of the current week (by start date); for each day:
//	// - check if there's capacity to schedule matches on that day
//	// - if so, filter teams to those that have matches that week
//	// - randomly pick a team from the list and assign it to that day
//	for currentDay := *firstDayOfMatches; !currentDay.After(*lastDayOfMatches); currentDay = currentDay.AddDate(0, 0, 1) {
//		currentWeek := weekKey(currentDay)
//		//fmt.Println(currentDay.Format("01/02/2006"), currentDay.Weekday())
//
//		currentDaySchedule := schedule.ForDay(currentDay)
//		if !currentDaySchedule.HasCapacity() {
//			// Day is full; cannot schedule
//			continue
//		}
//
//		if isBlackoutDate(currentDay, blackoutDates) {
//			//fmt.Println(currentDay.Format("20060102"), "is a blackout day, skipping it")
//			continue
//		}
//
//		// Daytime
//		if currentDaySchedule.HasDaytimeCapacity() {
//			candidateTeams := daytimeTeamsByWeek[currentWeek]
//			if len(candidateTeams) > 0 {
//				chosenTeamIdx := rand.Intn(len(candidateTeams))
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
//			//fmt.Println("candidate teams:", candidateTeams)
//			if len(candidateTeams) > 0 {
//				chosenTeamIdx := rand.Intn(len(candidateTeams))
//				chosenTeam := candidateTeams[chosenTeamIdx]
//
//				// Remove chosen team from teams by week, so it's not chosen again
//				eveningTeamsByWeek[currentWeek] = removeTeam(candidateTeams, chosenTeam)
//
//				// Assign chosen team to schedule
//				currentDaySchedule.EveningTeam = &chosenTeam
//
//				//fmt.Printf("- Assigned [%s] to evening\n", chosenTeam.Name)
//			}
//		}
//	}
//}
