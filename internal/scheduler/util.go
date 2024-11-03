package scheduler

import (
	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"slices"
	"time"
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
