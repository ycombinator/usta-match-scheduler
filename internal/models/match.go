package models

import (
	"fmt"
	"strings"
	"time"
)

type MatchLocation int

const (
	MatchLocationHome MatchLocation = iota
	MatchLocationAway
)

func MatchLocationFromString(location string) (MatchLocation, error) {
	switch strings.ToLower(location) {
	case "home":
		return MatchLocationHome, nil
	case "away":
		return MatchLocationAway, nil
	default:
		return 0, fmt.Errorf("expected match location to be either home or away; got [%s] instead", location)
	}
}

type Match struct {
	Date        time.Time `json:"date"`
	IsScheduled bool      `json:"is_scheduled"`

	HomeTeam     Team `json:"home_team"`
	VisitingTeam Team `json:"visiting_team"`

	Outcome MatchOutcome `json:"outcome"`
}

func (m *Match) LocationFor(t Team) (MatchLocation, error) {
	if m.HomeTeam.ID == t.ID {
		return MatchLocationHome, nil
	}

	if m.VisitingTeam.ID == t.ID {
		return MatchLocationAway, nil
	}

	return 0, fmt.Errorf("team [%d] is not involved in this match", t.ID)
}

type MatchOutcome struct {
	WinningTeam  Team `json:"winning_team"`
	WinnerPoints int  `json:"winner_points"`
	LoserPoints  int  `json:"loser_points"`
}
