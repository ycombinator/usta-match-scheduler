import { useEffect, useState } from 'react'
import './App.css'
import { TeamPreferences } from './components/TeamPreferences'
import { CalendarMonthGroup } from './components/CalendarMonthGroup'

const asrcOrganizationID = 225
const initialDayPreferences = {
    0: false, // Monday
    1: false, // Tuesday
    2: false, // Wednesday
    3: false, // Thursday
    4: false, // Friday
    5: false, // Saturday
    6: false, // Sunday
}

function App() {
    const now = new Date()
    const [startYearMonth, setStartYearMonth] = useState({year: now.getFullYear(), month: now.getMonth()})

    const startYear = startYearMonth.year
    const startMonth = startYearMonth.month

    const events = [
        { start: new Date("2025-04-08T22:00:00Z"), end: new Date("2025-04-09T12:00:00Z"), title: "Test multiday"},
        { start: new Date("2025-04-12T19:00:00Z"), end: new Date("2025-04-12T22:00:00Z"), title: "[W3.5] vs. Morgan Hill Tennis Club"},
        { start: new Date("2025-04-13T16:00:00Z"), end: new Date("2025-04-13T19:00:00Z"), title: "[W3.5DT] vs. Bay Club Courtside"},
        { start: new Date("2025-04-13T19:30:00Z"), end: new Date("2025-04-13T22:30:00Z"), title: "[M4.5] vs. Los Gatos"},
        { start: new Date("2025-04-13T23:00:00Z"), end: new Date("2025-04-14T02:00:00Z"), title: "[M3.5] vs. Bramhall"},
    ]

    const matches = [
        { date: new Date("2025-05-03T07:00:00Z"), home_team: { id: 105417 } },
        { date: new Date("2025-05-15T07:00:00Z"), home_team: { id: 105417 } },
        { date: new Date("2025-05-21T07:00:00Z"), home_team: { id: 105587 } },
        { date: new Date("2025-05-22T07:00:00Z"), home_team: { id: 105772 } },
        { date: new Date("2025-05-28T07:00:00Z"), home_team: { id: 105587 } },
        { date: new Date("2025-06-08T07:00:00Z"), home_team: { id: 105772 } },
        { date: new Date("2025-06-13T07:00:00Z"), home_team: { id: 105587 } },
    ]

    const [teams, setTeams] = useState([]);
    useEffect(() => {
        fetchUpcomingTeams(asrcOrganizationID)
        .then(teams => { teams.forEach(team => team.day_preferences = structuredClone(initialDayPreferences)); return teams })
        .then(setTeams)
    }, [])

    const changeDayPreference = function(teamIdx, dayIdx) {
        const preference = teams[teamIdx].day_preferences[dayIdx]
        const newTeams = structuredClone(teams)
        if (preference == true) {
            newTeams[teamIdx].day_preferences[dayIdx] = false
        } else {
            newTeams[teamIdx].day_preferences[dayIdx] = true
        }
        setTeams(newTeams)
    }

    // States:
    // - Blackout dates set
    // - Team preferences set
    // - Schedule generated
    const appState = "set_team_preferences"
    let component = <TeamPreferences teams={teams} changeDayPreference={changeDayPreference} />
    switch (appState) {
        case "set_blackout_dates":
            component = <BlackoutDates events={events} />
        case "edit_schedule":
            component = <EditSchedule events={events} />
    }

    return (
        <div className="App">
            <h1>USTA match scheduler</h1>
            { component }
        </div>
    )
}

async function fetchUpcomingTeams(organizationID) {
    return fetch("http://localhost:3000/api/usta/organization/"+organizationID+"/teams?upcoming=true")
    .then(r => r.json())
    .then(j => j.teams)
}

export default App