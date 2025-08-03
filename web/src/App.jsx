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
    const [blackoutEvents, setBlackoutEvents] = useState([])
    // const [events, setEvents] = useState([
    //     { start: new Date("2025-07-08T16:00:00Z"), end: new Date("2025-07-08T20:00:00Z"), title: "Club social", type:"blackout", slot:"morning"},
    //     { start: new Date("2025-07-12T19:00:00Z"), end: new Date("2025-07-12T22:00:00Z"), title: "[W3.5] vs. Morgan Hill Tennis Club", type:"match", slot:"afternoon"},
    //     { start: new Date("2025-07-13T16:00:00Z"), end: new Date("2025-07-13T19:00:00Z"), title: "[W3.5DT] vs. Bay Club Courtside", type:"match", slot:"morning"},
    //     { start: new Date("2025-07-13T19:30:00Z"), end: new Date("2025-07-13T22:30:00Z"), title: "[M4.5] vs. Los Gatos", type:"match", slot:"afternoon"},
    //     { start: new Date("2025-07-13T23:00:00Z"), end: new Date("2025-07-14T02:00:00Z"), title: "[M3.5] vs. Bramhall", type:"match", slot:"evening"},
    // ])

    const [teams, setTeams] = useState([]);
    useEffect(async () => {
        const teams = await fetchUpcomingTeams(asrcOrganizationID)
        teams.forEach(team => team.day_preferences = [])
        setTeams(teams)
    }, [])

    const changePreferredMatchDays = function(teamIdx, days) {
        const newTeams = structuredClone(teams)
        newTeams[teamIdx].day_preferences = days
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
            console.log("adding new event: ", e)
            // Add new event
            newEvents.push(e)
        }

        setEvents(newEvents)
    }

    console.log({events})

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
            header = <h3>Add club events, USTA Sectionals events, etc.</h3>
            component = <CalendarMonthGroup
                startYear={startYear} 
                startMonth={startMonth} 
                numMonths={1} 
                setStartYearMonth={setStartYearMonth} 
                events={events} 
                setEvent={setEvent}
                addEventLabel="blackout"
                allowAdds={true}
                allowDeletes={true}
                header={header}
            />
            step = 2
            stepLabel = "Set blackout slots"

            navPreviousLabel = "Set team preferences"
            navPrevious = () => setAppState("set_team_preferences")
            navNextLabel = "Generate schedule"
            navNext = async () => {
                setBlackoutEvents(events)
                const schedule = await generateSchedule(teams, events)
                setEvents(schedule.scheduled_events)
                setAppState("edit_schedule")
            }
            break
        case "edit_schedule":
            component = <CalendarMonthGroup
                startYear={startYear} 
                startMonth={startMonth} 
                numMonths={1} 
                setStartYearMonth={setStartYearMonth} 
                events={events} 
                setEvent={setEvent}
                addEventLabel="match"
                allowAdds={false}
                allowDeletes={false}
            />
            step = 3
            stepLabel = "Review schedule"

            navPreviousLabel = "Set blackout slots"
            navPrevious = () => {
                setEvents(blackoutEvents)
                setAppState("set_blackout_slots")
            }
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
    const response = await fetch("/api/usta/organization/"+organizationID+"/teams?upcoming=true")
    const json = await response.json()
    return json.te
    ams
    // return [json.teams[0], json.teams[1], json.teams[2]]
}

async function mockEvents(events) {
    const scheduleEvents = structuredClone(events)
    scheduleEvents.push(
        { start: new Date("2025-07-11T19:00:00Z"), end: new Date("2025-07-11T22:00:00Z"), title: "[CW3.5] vs. Morgan Hill Tennis Club", type:"match", slot:"evening"},
        { start: new Date("2025-07-13T16:00:00Z"), end: new Date("2025-07-13T19:00:00Z"), title: "[CW3.5DT] vs. Bay Club Courtside", type:"match", slot:"morning"},
        { start: new Date("2025-07-13T19:30:00Z"), end: new Date("2025-07-13T22:30:00Z"), title: "[CM4.5] vs. Los Gatos", type:"match", slot:"afternoon"},
        { start: new Date("2025-07-13T23:00:00Z"), end: new Date("2025-07-14T02:00:00Z"), title: "[CM3.5] vs. Bramhall", type:"match", slot:"evening"},        
        { start: new Date("2025-07-16T19:00:00Z"), end: new Date("2025-07-16T22:00:00Z"), title: "[CW2.5+DT] vs. Brookside", type:"match", slot:"morning"},
    )
    return new Promise((resolve, reject) => resolve(scheduleEvents))
}

async function generateSchedule(teams, events) {
    // return mockEvents(events)
    const response = await fetch("/api/schedule", {
        method: "POST",
        body: JSON.stringify({teams, events})
    })
    const json = await response.json()
    console.log({json})
    return json
}

export default App