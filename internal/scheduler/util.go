package scheduler

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

func mapTeamsByWeek(teams []models.Team) map[string][]models.Team {
	teamsByWeek := map[string][]models.Team{}
	for _, team := range teams {
		weeks := team.Weeks
		if len(weeks) == 0 {
			continue
		}

		slices.Sort(weeks)
		for _, week := range weeks {
			// Initialize if needed
			if _, exists := teamsByWeek[week]; !exists {
				teamsByWeek[week] = []models.Team{}
			}

			teamsByWeek[week] = append(teamsByWeek[week], team)
		}
	}

	return teamsByWeek
}

func teamsThatPreferDay(teams []models.Team, day time.Weekday) ([]models.Team, error) {
	teamsThatPreferDay := make([]models.Team, 0)
	for _, team := range teams {
		for _, preferredDay := range team.DayPreferences {
			weekday, err := weekdayFromStr(preferredDay)
			if err != nil {
				return nil, fmt.Errorf("cannot parse preferred day [%s] for team [%s]: %w", preferredDay, team.Name, err)
			}

			if weekday == day {
				teamsThatPreferDay = append(teamsThatPreferDay, team)
			}
		}
	}

	return teamsThatPreferDay, nil
}

func weekdayFromStr(day string) (time.Weekday, error) {
	day = strings.ToLower(strings.TrimSpace(day))
	if len(day) < 3 {
		return 0, fmt.Errorf("cannot parse given day [%s] as weekday", day)
	}

	dayAbbr := day[0:3]
	switch dayAbbr {
	case "sun":
		return time.Sunday, nil
	case "mon":
		return time.Monday, nil
	case "tue":
		return time.Tuesday, nil
	case "wed":
		return time.Wednesday, nil
	case "thu":
		return time.Thursday, nil
	case "fri":
		return time.Friday, nil
	case "sat":
		return time.Saturday, nil
	default:
		return 0, fmt.Errorf("cannot parse given day [%s] as weekday", day)
	}
}

// Golang weeks start on Sundays
// USTA schedule weeks start on Mondays
func weekKey(dt time.Time) string {
	dayOfWeek := dt.Weekday()
	diff := int(dayOfWeek - time.Monday)
	// Special case for Sunday
	if diff == -1 {
		diff = 6
	}

	mondayDt := dt.AddDate(0, 0, -diff)
	return mondayDt.Format("20060102")
}

func filterTeamsBySchedulingType(teams map[string]models.Team, schedulingType string) []models.Team {
	filtered := make([]models.Team, 0)
	for _, team := range teams {
		if team.SchedulingType == schedulingType {
			filtered = append(filtered, team)
		}
	}

	return filtered
}

func findTeamIdx(candidateTeams []models.Team, candidateTeam models.Team) int {
	teamIdx := -1
	for idx, team := range candidateTeams {
		if team.Name == candidateTeam.Name {
			return idx
		}
	}

	return teamIdx
}

func removeTeam(teams []models.Team, team models.Team) []models.Team {
	chosenTeamIdx := findTeamIdx(teams, team)
	//fmt.Println("chosen team idx: ", chosenTeamIdx)
	if chosenTeamIdx == -1 {
		// Team not found; nothing to remove
		return teams
	}

	newTeams := append(teams[:chosenTeamIdx], teams[chosenTeamIdx+1:]...)
	return newTeams
}

func isBlackoutDate(candidate time.Time, blackoutDates []string) bool {
	candidateStr := candidate.Format("20060102")
	fmt.Println(candidateStr, blackoutDates)
	for _, blackoutDate := range blackoutDates {
		if blackoutDate == candidateStr {
			return true
		}
	}

	return false
}
