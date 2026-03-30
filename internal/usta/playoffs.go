package usta

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

type PlayoffTeamData struct {
	Rounds int       `json:"rounds"`
	Weeks  []time.Time `json:"weeks"`
}

// PlayoffWeeksData maps team ID (as string) to playoff data.
type PlayoffWeeksData map[string]PlayoffTeamData

func LoadPlayoffWeeks(filePath string) (PlayoffWeeksData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read playoff weeks file [%s]: %w", filePath, err)
	}

	var result PlayoffWeeksData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unable to parse playoff weeks file [%s]: %w", filePath, err)
	}

	return result, nil
}

// CorrectPlayoffDates replaces dummy dates on playoff matches with the correct
// "week of" dates from the playoff weeks data. The nth playoff match (in order
// of appearance) gets the nth week date, up to the configured number of rounds.
func CorrectPlayoffDates(matches []models.TeamMatch, teamID int, playoffData PlayoffWeeksData) []models.TeamMatch {
	teamData, ok := playoffData[strconv.Itoa(teamID)]
	if !ok {
		return matches
	}

	playoffIndex := 0
	for i := range matches {
		if !matches[i].IsPlayoff {
			continue
		}

		if playoffIndex < teamData.Rounds && playoffIndex < len(teamData.Weeks) {
			matches[i].Date = teamData.Weeks[playoffIndex]
		}
		playoffIndex++
	}

	return matches
}
