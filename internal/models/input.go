package models

import (
	"sort"
	"time"
)

type DaySlot string

const (
	SlotMorning   DaySlot = "morning"
	SlotAfternoon DaySlot = "afternoon"
	SlotEvening   DaySlot = "evening"
)

type Input struct {
	Teams  []SchedulingTeam `json:"teams"`
	Events []Event          `json:"events"`
}

type schedulingType struct {
	StartTimes startTimes `yaml:"start_times"`
}

type startTimes struct {
	Monday    string `yaml:"monday"`
	Tuesday   string `yaml:"tuesday"`
	Wednesday string `yaml:"wednesday"`
	Thursday  string `yaml:"thursday"`
	Friday    string `yaml:"friday"`
	Saturday  string `yaml:"saturday"`
	Sunday    string `yaml:"sunday"`
}

func (i *Input) FirstDayOfMatches() *time.Time {
	startDates := sortStartDates(i.Teams)
	startDate := startDates[0]
	return &startDate
}

func (i *Input) LastDayOfMatches() *time.Time {
	startDates := sortStartDates(i.Teams)
	finalStartDate := startDates[len(startDates)-1]
	finalStartDate = finalStartDate.Add(6 * 24 * time.Hour)

	return &finalStartDate
}

func sortStartDates(teams []SchedulingTeam) []time.Time {
	var startDates []time.Time
	for _, team := range teams {
		if len(team.Weeks) == 0 {
			continue
		}

		startDates = append(startDates, team.Weeks...)
	}

	if len(startDates) == 0 {
		return nil
	}

	sort.Slice(startDates, func(i, j int) bool {
		return startDates[i].Before(startDates[j])
	})

	return startDates
}
