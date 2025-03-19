package main

import (
	"fmt"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"github.com/ycombinator/usta-match-scheduler/internal/scheduler"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

// TODO make CLI arg
const inputFilePath = "./input.yml"

func main() {
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var input models.Input
	if err := yaml.Unmarshal(data, &input); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	//fmt.Printf("%#+v\n", input)

	//s, err := scheduler.NewEager(input)
	s, err := scheduler.NewPreferring(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	schedule, err := s.Run()
	//s.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}

	// Print schedule
	schedule.ResetIterator()
	for day := schedule.Next(); day != nil; day = schedule.Next() {
		fmt.Printf("%12s: ", day.Date.Format("Mon, 01/02"))
		teamTitle := ""
		if day.DaytimeTeam != nil {
			teamTitle = day.DaytimeTeam.Title
		}
		fmt.Printf("%25s\t", teamTitle)

		teamTitle = ""
		if day.EveningTeam != nil {
			teamTitle = day.EveningTeam.Title
		}
		fmt.Printf("%15s\n", teamTitle)

		if day.Date.Weekday() == time.Sunday {
			fmt.Println()
		}
	}
}
