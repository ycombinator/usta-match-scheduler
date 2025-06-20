package usta

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

//go:embed testdata/organization_225.html
var organization225Html []byte

func getOrganizationUrl(id int) string {
	return fmt.Sprintf(baseURL+"/organization.asp?id=%d", id)
}

func WithFilterIsTeamSeasonUpcoming(u bool) TeamsFilterOpt {
	return func(f *TeamsFilter) {
		f.isSeasonUpcoming = &u
	}
}

type TeamsFilterOpt = func(f *TeamsFilter)

type TeamsFilter struct {
	isSeasonUpcoming *bool
}

func GetOrganizationTeams(id int, opts ...TeamsFilterOpt) ([]models.Team, error) {
	var f TeamsFilter
	for _, opt := range opts {
		opt(&f)
	}

	u := getOrganizationUrl(id)

	var body io.ReadCloser
	if useMockData() {
		body = io.NopCloser(bytes.NewReader(organization225Html))
		if f.isSeasonUpcoming != nil && *f.isSeasonUpcoming {
			f.isSeasonUpcoming = ptrTo[bool](false)
		}

	} else {
		fmt.Printf("Getting teams for organization [%d] from url [%s]...\n", id, u)

		resp, err := http.Get(u)
		if err != nil {
			return nil, fmt.Errorf("unable to get organization page from URL [%s]: %w", u, err)
		}
		defer resp.Body.Close()

		body = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization page from URL [%s]: %w", u, err)
	}

	var teams = make([]models.Team, 0)
	doc.Find("table tbody tr td table tr").Each(func(i int, sel *goquery.Selection) {
		// Parse team ID
		aSel := sel.Find("td").Next().Find("a")
		l, exists := aSel.Attr("href")
		if !exists {
			return
		}

		if !strings.Contains(l, "teaminfo.asp?") {
			return
		}

		teamID, err := parseTeamID(l)
		if err != nil {
			return
		}

		// Parse start date
		sdSel := sel.Find("td").Get(5)
		startDate, err := time.ParseInLocation("01/02/2006", sdSel.FirstChild.Data, time.Local)
		if err != nil {
			return
		}

		// Filter by isTeamSeasonUpcoming if that filter is set
		if f.isSeasonUpcoming != nil {
			if *f.isSeasonUpcoming {
				if startDate.Before(time.Now()) {
					return // Skip teams with a start date in the past
				}
			} else {
				if startDate.After(time.Now()) {
					return // Skip teams with a start date in the future
				}
			}
		}

		// Parse team name
		nSel := sel.Find("td").Get(1).FirstChild
		name := nSel.FirstChild.Data

		// Parse schedule type
		scheduleType := models.TeamScheduleGroupEvening // Default to evening
		if strings.HasSuffix(name, "-DT") {
			scheduleType = models.TeamScheduleGroupDaytime
		}

		// Parse captain
		cSel := sel.Find("td").Get(3)
		captain := cSel.FirstChild.Data

		team := models.Team{
			ID:            teamID,
			Captain:       captain,
			StartDate:     startDate,
			ScheduleGroup: scheduleType,
		}

		if err := team.SetRawName(name); err != nil {
			return
		}

		teams = append(teams, team)
	})

	return teams, nil
}
