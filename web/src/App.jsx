import { useEffect, useState } from 'react'
import './App.css'
import { TeamPreferences } from './components/TeamPreferences'
import { CalendarMonthGroup } from './components/CalendarMonthGroup'
import { Step } from './components/Step'
import { Nav } from './components/Nav'

const asrcOrganizationID = 225

function App() {
    const now = new Date()
    const [startYearMonth, setStartYearMonth] = useState({year: now.getFullYear(), month: now.getMonth()})
    const startYear = startYearMonth.year
    const startMonth = startYearMonth.month

    const [events, setEvents] = useState([])
    // const [events, setEvents] = useState([
    //     { start: new Date("2025-07-08T16:00:00Z"), end: new Date("2025-07-08T20:00:00Z"), title: "Club social", type:"blackout", slot:"morning"},
    //     { start: new Date("2025-07-12T19:00:00Z"), end: new Date("2025-07-12T22:00:00Z"), title: "[W3.5] vs. Morgan Hill Tennis Club", type:"match", slot:"afternoon"},
    //     { start: new Date("2025-07-13T16:00:00Z"), end: new Date("2025-07-13T19:00:00Z"), title: "[W3.5DT] vs. Bay Club Courtside", type:"match", slot:"morning"},
    //     { start: new Date("2025-07-13T19:30:00Z"), end: new Date("2025-07-13T22:30:00Z"), title: "[M4.5] vs. Los Gatos", type:"match", slot:"afternoon"},
    //     { start: new Date("2025-07-13T23:00:00Z"), end: new Date("2025-07-14T02:00:00Z"), title: "[M3.5] vs. Bramhall", type:"match", slot:"evening"},
    // ])

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

    const setEvent = function(e) {
        const newEvents = []
        let found = false
        for (let i = 0; i < events.length; i++) {
            if (e.id != events[i].id) {
                newEvents.push(events[i])
                continue
            }

            found = true 
            if (e.title == "") {
                // Delete event by not adding it to newEvents
                continue
            }

            // Update event
            newEvents.push(e)
        }

        if (!found) {
            // Add new event
            newEvents.push(e)
        }

        setEvents(newEvents)
    }

    // States:
    // - Team preferences set
    // - Blackout slots set
    // - Schedule generated
    const [ appState, setAppState ] = useState("set_team_preferences")
    let component, step, stepLabel
    let navPrevious, navNext, navPreviousLabel, navNextLabel
    const totalSteps = 3
    switch (appState) {
        case "set_team_preferences":
            component = <TeamPreferences teams={teams} changePreferredMatchDays={changePreferredMatchDays} />
            step = 1
            stepLabel = "Set team preferences"

            navNextLabel = "Set blackout slots"
            navNext = () => setAppState("set_blackout_slots")
            break
        case "set_blackout_slots":
            component = <CalendarMonthGroup
                startYear={startYear} 
                startMonth={startMonth} 
                numMonths={1} 
                setStartYearMonth={setStartYearMonth} 
                events={events} 
                setEvent={setEvent}
                addEventLabel="blackout"
            />
            step = 2
            stepLabel = "Set blackout slots"

            navPreviousLabel = "Set team preferences"
            navPrevious = () => setAppState("set_team_preferences")
            navNextLabel = "Generate schedule"
            navNext = () => setAppState("edit_schedule")
            break
        case "edit_schedule":
            component = <CalendarMonthGroup
                startYear={startYear} 
                startMonth={startMonth} 
                numMonths={1} 
                setStartYearMonth={setStartYearMonth} 
                events={events} 
                setEvent={setEvent}
            />
            step = 3
            stepLabel = "Review schedule"

            navPreviousLabel = "Set blackout slots"
            navPrevious = () => setAppState("set_blackout_slots")
            break
    }

    return (
        <div className="App">
            <header>
                <h1>USTA match scheduler</h1>
                <h2>
                    <Step current={step} total={totalSteps} label={stepLabel} />
                    <Nav
                        previous={navPrevious}
                        previousLabel={navPreviousLabel}
                        next={navNext}
                        nextLabel={navNextLabel}
                    />
                </h2>
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