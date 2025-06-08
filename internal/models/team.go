package models

import "time"

type TeamScheduleGroup string

const (
	TeamScheduleGroupDaytime TeamScheduleGroup = "daytime"
	TeamScheduleGroupEvening TeamScheduleGroup = "evening"
)

type Team struct {
	ID            int               `json:"id"`
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
