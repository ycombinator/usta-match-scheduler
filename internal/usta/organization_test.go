package usta

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"

	"github.com/stretchr/testify/require"
)

func TestGetOrganizationTeams(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(organization225Html)
	}))
	defer testServer.Close()

	baseURL = testServer.URL

	teams, err := GetOrganizationTeams(225)
	require.NoError(t, err)
	require.Len(t, teams, 38)

	expectedFirstTeam := models.Team{
		ID:            106665,
		Name:          "ALMADEN SR 40MX6.0A",
		Captain:       "Bui-Quang, Phu",
		StartDate:     time.Date(2025, 6, 9, 0, 0, 0, 0, time.Local),
		ScheduleGroup: models.TeamScheduleGroupEvening,
	}
	require.Equal(t, expectedFirstTeam, teams[0])
}

func TestHTTPGet(t *testing.T) {
	u := getOrganizationUrl(225)
	resp, err := http.Get(u)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
