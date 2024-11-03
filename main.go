package main

import (
	"fmt"
	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"github.com/ycombinator/usta-match-scheduler/internal/scheduler"
	"gopkg.in/yaml.v3"
	"os"
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
		fmt.Println(day.Date.Format("Mon, 01/02/2006"))
		if day.DaytimeTeam != nil {
			fmt.Println("  Daytime:", day.DaytimeTeam.Title)
		}
		if day.EveningTeam != nil {
			fmt.Println("  Evening:", day.EveningTeam.Title)
		}
	}
}
