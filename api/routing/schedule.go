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

	// Load playoff weeks data
	playoffData, err := usta.LoadPlayoffWeeks("data/playoff-weeks-225.json")
	if err != nil {
		logger.Info("No playoff weeks data loaded", "error", err)
		playoffData = make(usta.PlayoffWeeksData)
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

		// Correct playoff match dates using the playoff weeks data
		matches = usta.CorrectPlayoffDates(matches, team.ID, playoffData)

		var weeks = make([]time.Time, 0)
		playoffWeeks := make(map[time.Time]bool)
		for _, match := range matches {
			if match.IsScheduled {
				// Already-scheduled matches get placed directly as events,
				// bypassing the scheduling algorithm.
				slot := models.SlotFromHour(match.Date.Hour())
				event := models.Event{
					Title: team.DisplayName(),
					Type:  models.EventTypeMatch,
					Slot:  slot,
					Date:  match.Date,
				}
				input.Events = append(input.Events, event)
				logger.Debug("Pre-scheduled match", "team_name", team.DisplayName(), "date", match.Date.Format("2006-01-02 15:04"), "slot", slot)
			} else {
				// Unscheduled matches go to the algorithm for scheduling.
				weeks = append(weeks, match.Date)
				if match.IsPlayoff {
					playoffWeeks[match.Date] = true
				}
			}
		}

		logger.Debug("Team matches retrieved", "team_name", team.DisplayName(), "home_weeks", weeks)
		input.Teams[idx].Weeks = weeks
		input.Teams[idx].PlayoffWeeks = playoffWeeks
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
