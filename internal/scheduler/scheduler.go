package scheduler

import "github.com/ycombinator/usta-match-scheduler/internal/models"

type Scheduler interface {
	Run() (*models.Schedule, error)
}
