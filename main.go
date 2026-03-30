package main

import (
	"fmt"
	"os"

	"github.com/ycombinator/usta-match-scheduler/internal/models"
	"github.com/ycombinator/usta-match-scheduler/internal/scheduler"
	"gopkg.in/yaml.v3"
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
	s := scheduler.NewPreferring(&input)

	schedule, err := s.Run()
	//s.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}

	// Print schedule
	for _, event := range schedule.ScheduledEvents {
		fmt.Printf("%12s: %s [%s]\n", event.Date.Format("Mon, 01/02"), event.Title, event.Slot)
	}

	if len(schedule.UnscheduledEvents) > 0 {
		fmt.Printf("\nUnscheduled: %d events\n", len(schedule.UnscheduledEvents))
	}
}
