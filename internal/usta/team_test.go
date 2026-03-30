package usta

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/team_playoff.html
var teamPlayoffHtml []byte

func TestParsePlayoffMatches(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(teamPlayoffHtml))
	require.NoError(t, err)

	type parsedMatch struct {
		date      time.Time
		isPlayoff bool
		location  string
	}

	var matches []parsedMatch

	doc.Find("table tbody tr td table tbody tr").Each(func(i int, sel *goquery.Selection) {
		bgcolor, exists := sel.Attr("bgcolor")
		if !exists {
			return
		}
		if bgcolor != "white" && bgcolor != "#D2D2FF" {
			return
		}

		cells := sel.Find("td")
		if cells.Length() < 2 {
			return
		}

		// Parse date (cell 2)
		c := cells.Get(2).FirstChild
		if c.NextSibling != nil {
			c = c.NextSibling.FirstChild
		}
		v := strings.TrimSpace(c.Data)
		dt, err := time.ParseInLocation("01/02/06", v, tz)
		if err != nil {
			return
		}

		// Detect playoff (cell 1)
		isPlayoff := false
		if cells.Length() > 1 {
			roundText := strings.TrimSpace(cells.Eq(1).Text())
			if strings.Contains(roundText, "PlayOff") {
				isPlayoff = true
			}
		}

		// Parse location (cell 6)
		location := cells.Get(6).FirstChild.Data

		matches = append(matches, parsedMatch{
			date:      dt,
			isPlayoff: isPlayoff,
			location:  location,
		})
	})

	// Verify we parsed some matches
	require.NotEmpty(t, matches, "should parse at least one match from fixture")

	// Count playoff matches
	playoffCount := 0
	regularCount := 0
	for _, m := range matches {
		if m.isPlayoff {
			playoffCount++
		} else {
			regularCount++
		}
	}

	require.Greater(t, playoffCount, 0, "should detect at least one playoff match")
	require.Greater(t, regularCount, 0, "should have regular matches too")

	// Verify the last match is the playoff match
	lastMatch := matches[len(matches)-1]
	require.True(t, lastMatch.isPlayoff, "last match should be a playoff match")
	require.Equal(t, "Home", lastMatch.location, "playoff match should be Home")
}

func TestParsePlayoffMatchNotFalsePositive(t *testing.T) {
	// Verify that regular matches in team_1.html are not marked as playoff
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(team1Html))
	require.NoError(t, err)

	doc.Find("table tbody tr td table tbody tr").Each(func(i int, sel *goquery.Selection) {
		bgcolor, exists := sel.Attr("bgcolor")
		if !exists {
			return
		}
		if bgcolor != "white" && bgcolor != "#D2D2FF" {
			return
		}

		cells := sel.Find("td")
		if cells.Length() > 1 {
			roundText := strings.TrimSpace(cells.Eq(1).Text())
			require.False(t, strings.Contains(roundText, "PlayOff"),
				"regular team fixture should not contain PlayOff matches")
		}
	})
}
