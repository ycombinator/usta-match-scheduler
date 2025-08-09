package usta

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

const teamURL = "https://leagues.ustanorcal.com/teaminfo.asp?id=%d"

var (
	tz, _ = time.LoadLocation("America/Los_Angeles")
)

//go:embed testdata/team_1.html
var team1Html []byte

//go:embed testdata/team_2.html
var team2Html []byte

// TODO: memoize when we actual get team information over the network
func GetTeam(id int) (models.Team, error) {
	//fmt.Printf("Getting team [%d]...\n", id)
	return models.Team{
		ID: id,
	}, nil
}

func WithFilterIsMatchScheduled(s bool) TeamMatchesFilterOpt {
	return func(f *TeamMatchesFilter) {
		f.isScheduled = &s
	}
}

func WithFilterMatchLocation(l models.MatchLocation) TeamMatchesFilterOpt {
	return func(f *TeamMatchesFilter) {
		f.location = &l
	}
}

func WithFilterAfter(a time.Time) TeamMatchesFilterOpt {
	return func(f *TeamMatchesFilter) {
		f.after = &a
	}
}

func WithFilterBefore(b time.Time) TeamMatchesFilterOpt {
	return func(f *TeamMatchesFilter) {
		f.before = &b
	}
}

type TeamMatchesFilterOpt = func(f *TeamMatchesFilter)

type TeamMatchesFilter struct {
	isScheduled *bool
	location    *models.MatchLocation
	after       *time.Time
	before      *time.Time
}

func GetTeamMatches(t models.Team, opts ...TeamMatchesFilterOpt) ([]models.TeamMatch, error) {
	var f TeamMatchesFilter
	for _, opt := range opts {
		opt(&f)
	}

	u := fmt.Sprintf(teamURL, t.ID)

	var body io.ReadCloser
	if useMockData() {
		//body = io.NopCloser(bytes.NewReader(team1Html))
		var htmlSource []byte
		switch rand.Intn(2) {
		case 0:
			htmlSource = team1Html
		case 1:
			htmlSource = team2Html
		}
		body = io.NopCloser(bytes.NewReader(htmlSource))
	} else {
		//fmt.Printf("Getting matches for team [%d]...\n", t.ID)

		resp, err := http.Get(u)
		if err != nil {
			return nil, fmt.Errorf("unable to get team page from URL [%s]: %w", u, err)
		}
		defer resp.Body.Close()

		body = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("unable to read team page from URL [%s]: %w", u, err)
	}

	matches := make([]models.Match, 0)

	doc.Find("table tbody tr td table tbody tr").Each(func(i int, sel *goquery.Selection) {
		bgcolor, exists := sel.Attr("bgcolor")
		if !exists {
			return
		}

		if bgcolor != "white" && bgcolor != "#D2D2FF" {
			return
		}

		// Parse match date
		cells := sel.Find("td")
		if cells.Length() < 2 {
			return
		}

		c := cells.Get(2).FirstChild
		if c.NextSibling != nil {
			c = c.NextSibling.FirstChild
		}

		v := strings.TrimSpace(c.Data)
		dt, err := time.ParseInLocation("01/02/06", v, tz)
		if err != nil {
			return
		}

		// Parse match times
		c = cells.Get(4).FirstChild

		isScheduled := false
		v = strings.TrimSpace(c.Data)
		hour, minute, err := parseMatchTime(v)
		if err != nil {
			return
		}
		if hour > 0 {
			isScheduled = true
			dt = time.Date(dt.Year(), dt.Month(), dt.Day(), hour, minute, 0, 0, dt.Location())
		}

		// Filter by isScheduled if that filter is set
		if f.isScheduled != nil {
			if *f.isScheduled != isScheduled {
				return
			}
		}

		// Parse opposing team ID
		v = cells.Get(5).FirstChild.Attr[0].Val
		teamID, err := parseTeamID(v)
		if err != nil {
			return
		}

		o, err := GetTeam(teamID)
		if err != nil {
			return
		}

		// Parse location (home or away)
		location := sel.Find("td").Get(6).FirstChild.Data

		var homeTeam, visitingTeam models.Team
		if location == "Home" {
			homeTeam = t
			visitingTeam = o
		} else {
			homeTeam = o
			visitingTeam = t
		}

		m := models.Match{
			Date:         dt,
			HomeTeam:     homeTeam,
			VisitingTeam: visitingTeam,
			IsScheduled:  isScheduled,
		}

		// Parse outcome
		v = sel.Find("td").Get(7).FirstChild.Data
		verb, winnerPoints, loserPoints, err := parseMatchOutcome(v)
		if err != nil {
			return
		}

		if verb != "" {
			var winningTeam models.Team
			if verb == "Won" {
				winningTeam = t
			} else {
				winningTeam = o
			}

			outcome := models.MatchOutcome{
				WinningTeam:  winningTeam,
				WinnerPoints: winnerPoints,
				LoserPoints:  loserPoints,
			}

			m.Outcome = outcome
		}

		matches = append(matches, m)
	})

	// Convert Matches to TeamMatches
	teamMatches := make([]models.TeamMatch, 0, len(matches))
	for _, m := range matches {
		location, err := m.LocationFor(t)
		if err != nil {
			return nil, fmt.Errorf("unable to determine match location for team: %w", err)
		}

		// Filter by match location if that filter is set
		if f.location != nil {
			if *f.location != location {
				continue
			}
		}

		// Filter by match date if after filter is set
		if f.after != nil {
			if m.Date.Before(*f.after) {
				continue
			}
		}

		// Filter by match date if before filter is set
		if f.before != nil {
			if m.Date.After(*f.before) {
				continue
			}
		}

		tm := models.TeamMatch{
			Match:    m,
			Location: location,
		}
		teamMatches = append(teamMatches, tm)
	}

	return teamMatches, nil
}

func parseTeamID(u string) (int, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return 0, fmt.Errorf("could not parse team URL: %w", err)
	}

	v := pu.Query().Get("id")
	teamID, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("could not parse team ID from team URL: %w", err)
	}

	return int(teamID), nil
}
