package usta

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
)

func getOrganizationUrl(id int) string {
	return fmt.Sprintf(baseURL+"/organization.asp?id=%d", id)
}

func GetOrganizationTeams(id int) ([]models.Team, error) {
	u := getOrganizationUrl(id)
	fmt.Printf("Getting teams for organization [%d] from url [%s]...\n", id, u)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("unable to get organization page from URL [%s]: %w", u, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read organization page from URL [%s]: %w", u, err)
	}

	var teams = make([]models.Team, 0)
	doc.Find("table tbody tr td table tr").Each(func(i int, sel *goquery.Selection) {
		aSel := sel.Find("td").Next().Find("a")
		l, exists := aSel.Attr("href")
		if !exists {
			return
		}

		if !strings.Contains(l, "teaminfo.asp?") {
			return
		}

		// Parse team ID
		teamID, err := parseTeamID(l)
		if err != nil {
			return
		}

		// Parse name
		nSel := sel.Find("td").Get(1).FirstChild
		name := nSel.FirstChild.Data

		// Parse captain
		cSel := sel.Find("td").Get(3)
		captain := cSel.FirstChild.Data

		// Parse start date
		sdSel := sel.Find("td").Get(5)
		startDate, err := time.ParseInLocation("01/02/2006", sdSel.FirstChild.Data, time.Local)
		if err != nil {
			return
		}

		team := models.Team{
			ID:        teamID,
			Name:      name,
			Captain:   captain,
			StartDate: startDate,
		}

		teams = append(teams, team)
	})

	return teams, nil
}
