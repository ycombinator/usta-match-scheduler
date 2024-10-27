package models

type Team struct {
	Title          string   `yaml:"title"`
	Captain        string   `yaml:"captain"`
	SchedulingType string   `yaml:"scheduling_type"`
	DayPreferences []string `yaml:"day_preferences"`
	Weeks          []string `yaml:"weeks"`
}
