package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/usta"

	"github.com/ycombinator/usta-match-scheduler/internal/scheduler"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

func ScheduleMatches(w http.ResponseWriter, r *http.Request) {
	var input models.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	teamMatchesHomeFilter := usta.WithFilterMatchLocation(models.MatchLocationHome)
	for idx, team := range input.Teams {
		matches, err := usta.GetTeamMatches(team.Team, teamMatchesHomeFilter)
		fmt.Printf("team: [%s], matches: %v\n", team.DisplayName(), matches)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}

		var weeks = make([]time.Time, 0)
		for _, match := range matches {
			weeks = append(weeks, match.Date)
		}

		fmt.Printf("%s %v\n", team.DisplayName(), weeks)
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
