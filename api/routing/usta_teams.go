package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"github.com/ycombinator/usta-match-scheduler/internal/usta"
)

func GetUSTAOrganizationTeams(w http.ResponseWriter, r *http.Request) {
	orgId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, fmt.Errorf("invalid organization ID: %w", err), http.StatusBadRequest)
		return
	}

	// Parse query string parameters and setup filter options
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
	}

	filters := make([]usta.TeamsFilterOpt, 0)
	qUpcomingStr := q.Get("upcoming")
	if qUpcomingStr != "" {
		qUpcoming, err := strconv.ParseBool(qUpcomingStr)
		if err != nil {
			err = fmt.Errorf("expected upcoming to be true or false, got [%s] instead: %w", qUpcomingStr, err)
			handleError(w, err, http.StatusBadRequest)
			return
		}
		filters = append(filters, usta.WithFilterIsTeamSeasonUpcoming(qUpcoming))
	}

	t, err := usta.GetOrganizationTeams(orgId, filters...)
	if err != nil {
		handleError(w, fmt.Errorf("unable to get teams for organization: %w", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, ContentTypeApplicationJson)

	var teams struct {
		Teams []models.Team `json:"teams"`
	}
	teams.Teams = t
	j, err := json.Marshal(teams)
	if err != nil {
		handleError(w, fmt.Errorf("unable to marshal teams to JSON: %w", err), http.StatusInternalServerError)
		return
	}
	w.Write(j)
}
