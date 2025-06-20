package models

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type TeamScheduleGroup string

const (
	TeamScheduleGroupDaytime TeamScheduleGroup = "Daytime"
	TeamScheduleGroupEvening TeamScheduleGroup = "Evening"
)

type TeamType string

const (
	TeamTypeAdult TeamType = "Adult"
	TeamTypeMixed TeamType = "Mixed"
	TeamTypeCombo TeamType = "Combo"
)

type Gender string

const (
	GenderMale   Gender = "Men's"
	GenderFemale Gender = "Women's"
)

type Team struct {
	ID            int               `json:"id"`
	RawName       string            `yaml:"raw_name" json:"raw_name"`
	OrgName       string            `yaml:"org_name" json:"org_name"`
	MinAge        int               `yaml:"min_age" json:"min_age"`
	Type          TeamType          `yaml:"type" json:"type"`
	Level         string            `yaml:"level" json:"level"`
	Gender        Gender            `yaml:"gender" json:"gender"`
	Name          string            `yaml:"name" json:"name"`
	Captain       string            `yaml:"captain" json:"captain"`
	StartDate     time.Time         `yaml:"start_date" json:"start_date"`
	ScheduleGroup TeamScheduleGroup `yaml:"schedule_group" json:"schedule_group"`
}

type TeamMatch struct {
	Match
	Location MatchLocation `json:"location"`
}

type SchedulingTeam struct {
	Team
	DayPreferences []string `yaml:"day_preferences"`
	Weeks          []string `yaml:"weeks"`
}

func (t *Team) SetRawName(rawName string) error {
	// Example raw names -> OrgName, MinAge, Type, Level, Gender, Name, ScheduleGroup
	// ALMADEN SR 18MX6.0A (Skirts & Balls) -> ALMADEN SR, 18, TypeMixed, 6.0,, Skirts & Balls, Evening
	// ALMADEN SR 40AW2.5+A-DT (Ruby Smashers) -> ALMADEN SR, 40, TeamTypeAdult, 2.5+, GenderFemale, Ruby Smashers, Daytime
	// ALMADEN SR CM7.5A -> ALMADEN SR, 18, TeamTypeCombo, 7.5, GenderMale,, Evening

	var rawNameRegexp = regexp.MustCompile(`(.+)+\ (\d\d)?([A-Z])([A-Z])(\d\.\d\+?)([A-Z])(-DT)?\ ?(\((.+)\))?`)
	matches := rawNameRegexp.FindStringSubmatch(rawName)

	minAge := 18 // Default minimum age for USTA league
	if matches[2] != "" {
		var err error
		minAge, err = strconv.Atoi(matches[2])
		if err != nil {
			return fmt.Errorf("invalid age [%s] in raw name: %w", matches[2], err)
		}
	}

	var teamType TeamType
	switch matches[3] {
	case "M":
		teamType = TeamTypeMixed
	case "A":
		teamType = TeamTypeAdult
	case "C":
		teamType = TeamTypeCombo
	}

	var gender Gender
	switch matches[4] {
	case "M":
		gender = GenderMale
	case "W":
		gender = GenderFemale
	}

	scheduleGroup := TeamScheduleGroupEvening
	if matches[7] == "-DT" {
		scheduleGroup = TeamScheduleGroupDaytime
	}

	var name string
	if len(matches) > 8 {
		name = matches[9]
	}

	t.RawName = rawName
	t.OrgName = matches[1]
	t.MinAge = minAge
	t.Type = teamType
	t.Gender = gender
	t.Level = matches[5]
	t.ScheduleGroup = scheduleGroup
	t.Name = name

	return nil
}
