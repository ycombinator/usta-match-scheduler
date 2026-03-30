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

	originalBaseURL := baseURL
	baseURL = testServer.URL
	defer func() { baseURL = originalBaseURL }()

	teams, err := GetOrganizationTeams(225)
	require.NoError(t, err)
	require.Len(t, teams, 38)

	expectedFirstTeam := models.Team{
		ID:            106665,
		RawName:       "ALMADEN SR 40MX6.0A (Mixed Drinkers)",
		OrgName:       "ALMADEN SR",
		MinAge:        40,
		Type:          models.TeamTypeMixed,
		Level:         "6.0",
		Captain:       "Phu Bui-Quang",
		Name:          "Mixed Drinkers",
		StartDate:     time.Date(2026, 3, 2, 0, 0, 0, 0, time.Local),
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
