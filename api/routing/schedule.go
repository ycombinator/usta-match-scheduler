package routing

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/logging"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"github.com/ycombinator/usta-match-scheduler/internal/scheduler"
	"github.com/ycombinator/usta-match-scheduler/internal/usta"
)

func ScheduleMatches(w http.ResponseWriter, r *http.Request) {
	logger := logging.NewLogger()
	var input models.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	logger.Info("Getting team matches", "teams_count", len(input.Teams))
	teamMatchesHomeFilter := usta.WithFilterMatchLocation(models.MatchLocationHome)
	for idx, team := range input.Teams {
		matches, err := usta.GetTeamMatches(team.Team, teamMatchesHomeFilter)
		logger.Debug("Getting matches for team", "team_name", team.DisplayName(), "home_matches_count", len(matches))
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}

		var weeks = make([]time.Time, 0)
		for _, match := range matches {
			weeks = append(weeks, match.Date)
		}

		logger.Debug("Team matches retrieved", "team_name", team.DisplayName(), "home_weeks", weeks)
		input.Teams[idx].Weeks = weeks
	}

	s := scheduler.NewPreferring(&input)
	output, err := s.Run()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(output)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set(HeaderContentType, ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
