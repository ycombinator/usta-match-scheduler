package models

import (
	"slices"
	"time"
)

type Input struct {
	Teams           map[string]Team `yaml:"teams"`
	SchedulingTypes struct {
		Daytime schedulingType
		Evening schedulingType
	} `yaml:"scheduling_types"`
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

func (i *Input) FirstDayOfMatches() (*time.Time, error) {
	var startDates []string
	for _, team := range i.Teams {
		if len(team.Weeks) == 0 {
			continue
		}

		startDates = append(startDates, team.Weeks...)
	}

	if len(startDates) == 0 {
		return nil, nil
	}

	slices.Sort(startDates)

	startDate, err := time.Parse("20060102", startDates[0])
	if err != nil {
		return nil, err
	}

	return &startDate, nil
}

func (i *Input) LastDayOfMatches() (*time.Time, error) {
	var startDates []string
	for _, team := range i.Teams {
		if len(team.Weeks) == 0 {
			continue
		}

		startDates = append(startDates, team.Weeks...)
	}

	if len(startDates) == 0 {
		return nil, nil
	}

	slices.Sort(startDates)
	slices.Reverse(startDates)

	finalStartDate, err := time.Parse("20060102", startDates[0])
	if err != nil {
		return nil, err
	}
	finalStartDate = finalStartDate.Add(6 * 24 * time.Hour)

	return &finalStartDate, nil
}
