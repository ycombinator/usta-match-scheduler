package usta

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

const organizationURL = "https://leagues.ustanorcal.com/organization.asp?id=%d"

func GetOrganizationTeams(id int) ([]models.Team, error) {
	fmt.Printf("Getting teams for organization [%d]...\n", id)
	u := fmt.Sprintf(organizationURL, id)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("unable to get organization page from URL [%s]: %w", u, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization page from URL [%s]: %w", u, err)
	}

	var teamIDs []int

	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		u, exists := sel.Attr("href")
		if !exists {
			return
		}

		if strings.HasPrefix(u, "teaminfo.asp?") {
			teamID, err := parseTeamID(u)
			if err != nil {
				return
			}

			teamIDs = append(teamIDs, teamID)
		}
	})

	var teams = make([]models.Team, 0, len(teamIDs))

	var wg sync.WaitGroup
	for _, teamID := range teamIDs {
		wg.Add(1)
		go func(teamID int) {
			t, _ := GetTeam(teamID)
			teams = append(teams, t)
			wg.Done()
		}(teamID)
	}

	wg.Wait()
	return teams, nil
}
