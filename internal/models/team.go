package models

import "time"

type Team struct {
	ID        int       `json:"id"`
	Name      string    `yaml:"name" json:"name"`
	Captain   string    `yaml:"captain" json:"captain"`
	StartDate time.Time `yaml:"start_date" json:"start_date"`
}

type TeamMatch struct {
	Match
	Location MatchLocation `json:"location"`
}

type SchedulingTeam struct {
	Team
	SchedulingType string   `yaml:"scheduling_type"`
	DayPreferences []string `yaml:"day_preferences"`
	Weeks          []string `yaml:"weeks"`
}
