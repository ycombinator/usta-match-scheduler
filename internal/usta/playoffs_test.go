package usta

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

func TestLoadPlayoffWeeks(t *testing.T) {
	t.Run("valid_file", func(t *testing.T) {
		// Write a temp JSON file
		dir := t.TempDir()
		path := filepath.Join(dir, "playoff-weeks.json")
		content := `{
			"12345": {
				"rounds": 2,
				"weeks": [
					"2026-03-30T00:00:00-07:00",
					"2026-04-06T00:00:00-07:00",
					"2026-04-13T00:00:00-07:00"
				]
			}
		}`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		data, err := LoadPlayoffWeeks(path)
		require.NoError(t, err)
		require.Contains(t, data, "12345")
		require.Equal(t, 2, data["12345"].Rounds)
		require.Len(t, data["12345"].Weeks, 3)
	})

	t.Run("missing_file", func(t *testing.T) {
		_, err := LoadPlayoffWeeks("/nonexistent/path.json")
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to read playoff weeks file")
	})

	t.Run("invalid_json", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.json")
		err := os.WriteFile(path, []byte("not json"), 0644)
		require.NoError(t, err)

		_, err = LoadPlayoffWeeks(path)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to parse playoff weeks file")
	})

	t.Run("multiple_teams", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "multi.json")
		content := `{
			"111": {"rounds": 1, "weeks": ["2026-03-30T00:00:00-07:00"]},
			"222": {"rounds": 3, "weeks": ["2026-03-30T00:00:00-07:00", "2026-04-06T00:00:00-07:00", "2026-04-13T00:00:00-07:00"]}
		}`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		data, err := LoadPlayoffWeeks(path)
		require.NoError(t, err)
		require.Len(t, data, 2)
		require.Equal(t, 1, data["111"].Rounds)
		require.Equal(t, 3, data["222"].Rounds)
	})
}

func TestCorrectPlayoffDates(t *testing.T) {
	week1 := time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC)
	week2 := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)
	week3 := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)
	dummyDate := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	playoffData := PlayoffWeeksData{
		"100": {
			Rounds: 2,
			Weeks:  []time.Time{week1, week2, week3},
		},
	}

	t.Run("team_not_in_data", func(t *testing.T) {
		matches := []models.TeamMatch{
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
		}
		result := CorrectPlayoffDates(matches, 999, playoffData)
		// Date should be unchanged
		require.Equal(t, dummyDate, result[0].Date)
	})

	t.Run("no_playoff_matches", func(t *testing.T) {
		regularDate := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
		matches := []models.TeamMatch{
			{Match: models.Match{Date: regularDate, IsPlayoff: false}},
		}
		result := CorrectPlayoffDates(matches, 100, playoffData)
		require.Equal(t, regularDate, result[0].Date)
	})

	t.Run("corrects_up_to_rounds", func(t *testing.T) {
		// Team 100 has rounds=2, so only first 2 playoff matches get corrected
		matches := []models.TeamMatch{
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
		}
		result := CorrectPlayoffDates(matches, 100, playoffData)
		require.Equal(t, week1, result[0].Date)
		require.Equal(t, week2, result[1].Date)
		// Third playoff match should keep dummy date (beyond rounds limit)
		require.Equal(t, dummyDate, result[2].Date)
	})

	t.Run("skips_regular_matches", func(t *testing.T) {
		regularDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		matches := []models.TeamMatch{
			{Match: models.Match{Date: regularDate, IsPlayoff: false}},
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
		}
		result := CorrectPlayoffDates(matches, 100, playoffData)
		require.Equal(t, regularDate, result[0].Date)
		require.Equal(t, week1, result[1].Date)
		require.Equal(t, week2, result[2].Date)
	})

	t.Run("single_round", func(t *testing.T) {
		singleRoundData := PlayoffWeeksData{
			"200": {
				Rounds: 1,
				Weeks:  []time.Time{week1, week2},
			},
		}
		matches := []models.TeamMatch{
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
		}
		result := CorrectPlayoffDates(matches, 200, singleRoundData)
		require.Equal(t, week1, result[0].Date)
		// Second playoff match is beyond rounds=1, not corrected
		require.Equal(t, dummyDate, result[1].Date)
	})

	t.Run("empty_matches", func(t *testing.T) {
		result := CorrectPlayoffDates(nil, 100, playoffData)
		require.Nil(t, result)
	})

	t.Run("empty_playoff_data", func(t *testing.T) {
		matches := []models.TeamMatch{
			{Match: models.Match{Date: dummyDate, IsPlayoff: true}},
		}
		result := CorrectPlayoffDates(matches, 100, PlayoffWeeksData{})
		require.Equal(t, dummyDate, result[0].Date)
	})
}
