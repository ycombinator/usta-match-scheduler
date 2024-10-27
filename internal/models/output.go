package models

import "time"

type Output []booking

type booking struct {
	StartTime time.Time
	EndTime   time.Time
}
