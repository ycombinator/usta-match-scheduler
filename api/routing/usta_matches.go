package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"

	"github.com/ycombinator/usta-match-scheduler/internal/usta"
)

func GetUSTAOrganizationMatches(w http.ResponseWriter, r *http.Request) {
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

	filters := make([]usta.TeamMatchesFilterOpt, 0)
	qIsScheduledStr := q.Get("is_scheduled")
	if qIsScheduledStr != "" {
		qIsScheduled, err := strconv.ParseBool(qIsScheduledStr)
		if err != nil {
			err = fmt.Errorf("expected is_scheduled to be true or false, got [%s] instead: %w", qIsScheduledStr, err)
			handleError(w, err, http.StatusBadRequest)
			return
		}
		filters = append(filters, usta.WithFilterIsMatchScheduled(qIsScheduled))
	}

	qLocationStr := q.Get("location")
	if qLocationStr != "" {
		qLocation, err := models.MatchLocationFromString(qLocationStr)
		if err != nil {
			err = fmt.Errorf("expected location to be home or away, got [%s] instead: %w", qLocationStr, err)
			handleError(w, err, http.StatusBadRequest)
			return
		}
		filters = append(filters, usta.WithFilterMatchLocation(qLocation))
	}

	qAfterStr := q.Get("after")
	if qAfterStr != "" {
		qAfter, err := time.Parse(time.RFC3339, qAfterStr)
		if err != nil {
			err = fmt.Errorf("could not parse after time [%s]: %w", qAfterStr, err)
			handleError(w, err, http.StatusBadRequest)
			return
		}
		filters = append(filters, usta.WithFilterAfter(qAfter))
	}

	qBeforeStr := q.Get("before")
	if qBeforeStr != "" {
		qBefore, err := time.Parse(time.RFC3339, qBeforeStr)
		if err != nil {
			err = fmt.Errorf("could not parse before time [%s]: %w", qBeforeStr, err)
			handleError(w, err, http.StatusBadRequest)
			return
		}
		filters = append(filters, usta.WithFilterBefore(qBefore))
	}

	// Get teams for organization
	teams, err := usta.GetOrganizationTeams(orgId)

	w.Header().Set(HeaderContentType, ContentTypeApplicationJson)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Get matches for teams
	matches := make([]models.TeamMatch, 0)
	var wg sync.WaitGroup
	for _, team := range teams {
		wg.Add(1)
		go func(t models.Team) {
			m, _ := usta.GetTeamMatches(
				t,
				filters...,
			)
			matches = append(matches, m...)
			wg.Done()
		}(team)
	}

	wg.Wait()

	// Create response and send it
	var resp struct {
		Matches []models.TeamMatch `json:"matches"`
	}
	resp.Matches = matches

	respJson, err := json.Marshal(resp)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Write(respJson)
}
