package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTeamParseRawName(t *testing.T) {
	cases := map[string]struct {
		rawName  string
		expected Team
	}{
		// ALMADEN SR 18MX6.0A (Skirts & Balls) -> ALMADEN SR, 18, TeamTypeMixed, 6.0,, Skirts & Balls, TeamScheduleGroupEvening
		// ALMADEN SR 40AW2.5+A-DT (Ruby Smashers) -> ALMADEN SR, 40, TeamTypeAdult, 2.5+, GenderFemale, Ruby Smashers, TeamScheduleGroupDaytime
		// ALMADEN SR CM7.5A -> ALMADEN SR, 18, TeamTypeCombo, 7.5, GenderMale,, Evening
		"18_mixed_6.0": {
			rawName: "ALMADEN SR 18MX6.0A (Skirts & Balls)",
			expected: Team{
				RawName:       "ALMADEN SR 18MX6.0A (Skirts & Balls)",
				OrgName:       "ALMADEN SR",
				MinAge:        18,
				Type:          TeamTypeMixed,
				Level:         "6.0",
				Gender:        "",
				Name:          "Skirts & Balls",
				ScheduleGroup: TeamScheduleGroupEvening,
			},
		},
		"40_womens_2.5+_daytime": {
			rawName: "ALMADEN SR 40AW2.5+A-DT (Ruby Smashers)",
			expected: Team{
				RawName:       "ALMADEN SR 40AW2.5+A-DT (Ruby Smashers)",
				OrgName:       "ALMADEN SR",
				MinAge:        40,
				Type:          TeamTypeAdult,
				Level:         "2.5+",
				Gender:        GenderFemale,
				Name:          "Ruby Smashers",
				ScheduleGroup: TeamScheduleGroupDaytime,
			},
		},
		"combo_mens_7.5_evening": {
			rawName: "ALMADEN SR CM7.5A",
			expected: Team{
				RawName:       "ALMADEN SR CM7.5A",
				OrgName:       "ALMADEN SR",
				MinAge:        18, // Default for combo teams
				Type:          TeamTypeCombo,
				Level:         "7.5",
				Gender:        GenderMale,
				Name:          "",
				ScheduleGroup: TeamScheduleGroupEvening,
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			var team Team
			team.SetRawName(test.rawName)
			require.Equal(t, test.expected, team)
		})
	}
}
