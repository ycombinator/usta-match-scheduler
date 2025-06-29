import { useEffect, useState } from 'react'
import './App.css'
import { TeamPreferences } from './components/TeamPreferences'
import { CalendarMonthGroup } from './components/CalendarMonthGroup'
import { Step } from './components/Step'

const asrcOrganizationID = 225

function App() {
    const now = new Date()
    const [startYearMonth, setStartYearMonth] = useState({year: now.getFullYear(), month: now.getMonth()})

    const startYear = startYearMonth.year
    const startMonth = startYearMonth.month

    const events = [
        { start: new Date("2025-04-08T16:00:00Z"), end: new Date("2025-04-08T20:00:00Z"), title: "Club social", type:"blackout", slot:"morning"},
        { start: new Date("2025-04-12T19:00:00Z"), end: new Date("2025-04-12T22:00:00Z"), title: "[W3.5] vs. Morgan Hill Tennis Club", type:"match", slot:"afternoon"},
        { start: new Date("2025-04-13T16:00:00Z"), end: new Date("2025-04-13T19:00:00Z"), title: "[W3.5DT] vs. Bay Club Courtside", type:"match", slot:"morning"},
        { start: new Date("2025-04-13T19:30:00Z"), end: new Date("2025-04-13T22:30:00Z"), title: "[M4.5] vs. Los Gatos", type:"match", slot:"afternoon"},
        { start: new Date("2025-04-13T23:00:00Z"), end: new Date("2025-04-14T02:00:00Z"), title: "[M3.5] vs. Bramhall", type:"match", slot:"evening"},
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
        .then(teams => { teams.forEach(team => team.preferred_match_days = []); return teams })
        .then(setTeams)
    }, [])

    const changePreferredMatchDays = function(teamIdx, days) {
        const newTeams = structuredClone(teams)
        newTeams[teamIdx].preferred_match_days = days
        setTeams(newTeams)
    }

    const addEvent = function(e) {
        // TODO: 
    }

    // States:
    // - Team preferences set
    // - Blackout slots set
    // - Schedule generated
    // const appState = "set_team_preferences"
    const appState = "set_blackout_slots"
    let component, step, stepLabel
    const totalSteps = 3
    switch (appState) {
        case "set_team_preferences":
            component = <TeamPreferences teams={teams} changePreferredMatchDays={changePreferredMatchDays} />
            step = 1
            stepLabel = "Set team preferences"
            break
        case "set_blackout_slots":
            component = <CalendarMonthGroup
                startYear={startYear} 
                startMonth={startMonth} 
                numMonths={1} 
                setStartYearMonth={setStartYearMonth} 
                events={events} 
                addEvent={addEvent}
                addEventLabel="blackout"
            />
            step = 2
            stepLabel = "Set blackout slots"
            break
        case "edit_schedule":
            component = <CalendarMonthGroup
                startYear={startYear} 
                startMonth={startMonth} 
                numMonths={1} 
                setStartYearMonth={setStartYearMonth} 
                events={events} 
                addEvent={addEvent}
            />
            step = 3
            stepLabel = "Review schedule"
            break
    }

    return (
        <div className="App">
            <header>
                <h1>USTA match scheduler</h1>
                <h2><Step current={step} total={totalSteps} label={stepLabel} /></h2>
            </header>
            { component }
        </div>
    )
}

async function fetchUpcomingTeams(organizationID) {
    return fetch("/api/usta/organization/"+organizationID+"/teams?upcoming=true")
    .then(r => r.json())
    .then(j => j.teams)
}

export default App