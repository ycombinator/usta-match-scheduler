package usta

import (
	"regexp"
	"strconv"
	"strings"
)

func parseMatchOutcome(outcome string) (string, int, int, error) {
	outcome = strings.TrimSpace(outcome)
	parts := strings.Split(outcome, " ")
	if len(parts) != 2 {
		return "", 0, 0, nil
	}

	verb := parts[0]
	points := strings.Split(parts[1], "-")

	points1, err := strconv.ParseInt(points[0], 10, 0)
	if err != nil {
		return "", 0, 0, err
	}

	points2, err := strconv.ParseInt(points[1], 10, 0)
	if err != nil {
		return "", 0, 0, err
	}

	var winnerPoints, loserPoints int64
	if points1 > points2 {
		winnerPoints = points1
		loserPoints = points2
	} else {
		winnerPoints = points2
		loserPoints = points1
	}

	return verb, int(winnerPoints), int(loserPoints), nil
}

func parseMatchTime(u string) (int, int, error) {
	u = strings.TrimSpace(u)
	if u == "" {
		return 0, 0, nil
	}

	regex, err := regexp.Compile(`at[^\d]+(\d+):(\d\d)\s+([aApP]M)`)
	if err != nil {
		return 0, 0, err
	}

	parts := regex.FindStringSubmatch(u)
	if len(parts) < 4 {
		return 0, 0, nil
	}
	hour, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return 0, 0, err
	}

	minute, err := strconv.Atoi(string(parts[2]))
	if err != nil {
		return 0, 0, err
	}

	if strings.ToLower(string(parts[3])) == "pm" && hour < 12 {
		hour += 12
	}

	return hour, minute, nil
}
