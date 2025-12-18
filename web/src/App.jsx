import React from 'react'
import { toPng } from 'html-to-image';

import './App.css'
import { TeamPreferences } from './components/TeamPreferences'
import { CalendarMonthGroup } from './components/CalendarMonthGroup'
import { Step } from './components/Step'
import { Nav } from './components/Nav'

const asrcOrganizationID = 225
const EVENTS_STORAGE_KEY = 'events';

export default class App extends React.Component {
    constructor(props) {
        super(props)
        this.componentRef = React.createRef();

        let appState = "set_team_preferences"

        const now = new Date()
        let year = now.getFullYear()
        let month = now.getMonth()
        let numMonths = 1

        const events = loadEventsFromStorage()
        if (events.length > 0) {
            appState = "edit_schedule"
            const { eventsStartYear, eventsStartMonth, eventsNumMonths } = findEventsBounds(events)
            year = eventsStartYear
            month = eventsStartMonth
            numMonths = eventsNumMonths
        }

        this.state = {
            appState: appState,
            year: year,
            month: month,
            numMonths: numMonths,
            teams: [],
            events: events,
            blackoutEvents: [],
            knownEvents: [],
            isGeneratingSchedule: false,
            isPrinting: false,
        }
    }

    async componentDidMount() {
            const knownEvents = await fetchKnownEvents(asrcOrganizationID)
            knownEvents.forEach(event => {
                event.date = new Date(event.date)
            })
            this.setKnownEvents(knownEvents)
            // this.setKnownEvents([])

            const teams = await fetchUpcomingTeams(asrcOrganizationID)
            teams.forEach(team => team.day_preferences = [])
            this.setTeams(teams)
    }

    setKnownEvents = events => this.setState({knownEvents: events})

    setTeams = (teams) => this.setState({teams})

    render() {
        const self = this

        const setAppState = appState => this.setState({appState})
        const setCalendarBounds = (year, month, numMonths) => {
            this.setState({year, month})
            if (numMonths) {
                this.setState({numMonths})
            }
        }
        const setTeams = this.setTeams
        const setEvents = events => this.setState({events}, () => {
            saveEventsToStorage(events);
        })
        const moveEvent = (fromID, toID) => {
            // console.log({fromID, toID})
            const newEvents = structuredClone(self.state.events)
            newEvents.forEach(event => {
                // console.log(event)
                if (event.id == fromID) {
                    event.id = toID
                    const [ year, month, day, slot] = toID.split("_")
                    let currentDate = new Date(event.date)
                    currentDate.setFullYear(year)
                    currentDate.setMonth(month)
                    currentDate.setDate(day)
                    event.date = currentDate.toISOString()
                    event.slot = slot
                    // console.log({event})
                }
            })
            setEvents(newEvents)
        }
        const setBlackoutEvents = events => this.setState({blackoutEvents: events})
        const setIsGeneratingSchedule = isGeneratingSchedule => this.setState({isGeneratingSchedule})

        const changePreferredMatchDays = function(teamIdx, days) {
            const newTeams = structuredClone(self.state.teams)
            newTeams[teamIdx].day_preferences = days
            setTeams(newTeams)
        }

        const setEvent = function(e) {
            const newEvents = []
            let found = false
            for (let i = 0; i < self.state.events.length; i++) {
                if (e.id != self.state.events[i].id) {
                    newEvents.push(self.state.events[i])
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

        console.log(this.state)

        // States:
        // - Team preferences set
        // - Blackout slots set
        // - Schedule generated
        // const [ appState, setAppState ] = useState("set_team_preferences")
        let component, step, stepLabel, header
        let navPrevious, navNext, navPreviousLabel, navNextLabel
        let isPreviousProcessing, isNextProcessing
        const totalSteps = 3
        switch (this.state.appState) {
            case "set_team_preferences":
                component = <TeamPreferences 
                    teams={self.state.teams} 
                    changePreferredMatchDays={changePreferredMatchDays} 
                />
                step = 1
                stepLabel = "Set team preferences"

                navNextLabel = "Set blackout slots"
                navNext = () => setAppState("set_blackout_slots")
                break
            case "set_blackout_slots":
                header = <h5>Blackout any slots where you don't want matches to be scheduled, e.g. for club events.</h5>
                component = <CalendarMonthGroup
                    startYear={this.state.year}
                    startMonth={this.state.month}
                    numMonths={this.state.numMonths}
                    setStartYearMonth={setCalendarBounds}
                    events={self.state.events}
                    setEvent={setEvent}
                    addEventLabel="blackout"
                    allowAdds={true}
                    allowEdits={true}
                    allowDeletes={true}
                    allowMoves={false}
                    header={header}
                    knownEvents={self.state.knownEvents}
                />
                step = 2
                stepLabel = "Set blackout slots"

                navPreviousLabel = "Set team preferences"
                navPrevious = () => {
                    setEvents([])
                    setAppState("set_team_preferences")
                }
                navNextLabel = this.state.isGeneratingSchedule ? "Generating..." : "Generate schedule"
                isNextProcessing = this.state.isGeneratingSchedule
                navNext = async () => {
                    setBlackoutEvents(self.state.events)
                    setIsGeneratingSchedule(true)
                    const schedule = await generateSchedule(self.state.teams, self.state.events)
                    setEvents(schedule.scheduled_events)

                    const { eventsStartYear, eventsStartMonth, eventsNumMonths } = findEventsBounds(schedule.scheduled_events)
                    // console.log({ eventsStartYear, eventsStartMonth, eventsNumMonths })
                    setCalendarBounds(eventsStartYear, eventsStartMonth, eventsNumMonths)
                    setIsGeneratingSchedule(false)
                    setAppState("edit_schedule")
                }
                break
            case "edit_schedule":
                component = <div ref={this.componentRef}>
                        <CalendarMonthGroup
                            startYear={this.state.year}
                            startMonth={this.state.month}
                            numMonths={this.state.numMonths}
                            setStartYearMonth={setCalendarBounds}
                            events={self.state.events}
                            setEvent={setEvent}
                            moveEvent={moveEvent}
                            addEventLabel="match"
                            allowAdds={false}
                            allowEdits={false}
                            allowDeletes={true}
                            allowMoves={!this.state.isPrinting}
                            knownEvents={self.state.knownEvents}
                        />
                    </div>
                step = 3
                stepLabel = "Review schedule"

                navPreviousLabel = "Set blackout slots"
                navPrevious = () => {
                    setEvents(self.state.blackoutEvents)
                    const now = new Date()
                    setCalendarBounds(now.getFullYear(), now.getMonth(), 1)
                    setAppState("set_blackout_slots")
                }

                navNextLabel = "Print"
                navNext = () => {
                    if (self.componentRef === null) {
                        return
                    }

                    self.setState({isPrinting: true})

                    toPng(self.componentRef.current, { cacheBust: true })
                        .then((dataUrl) => {
                            const link = document.createElement('a');
                            link.download = 'schedule.png';
                            link.href = dataUrl;
                            link.click();
                        })
                        .catch((err) => {
                            console.error('Oops, something went wrong!', err);
                        })
                        .finally(() => self.setState({isPrinting: false}))
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
                            isPreviousProcessing={isPreviousProcessing}
                            next={navNext}
                            nextLabel={navNextLabel}
                            isNextProcessing={isNextProcessing}
                        />
                    </h2>
                </header>
                { component }
            </div>
        )
    }
}

async function fetchUpcomingTeams(organizationID) {
    const response = await fetch("/api/usta/organization/"+organizationID+"/teams?upcoming=true")
    const json = await response.json()
    return json.teams
    // return [json.teams[0], json.teams[1], json.teams[2]]
}

async function fetchKnownEvents(organizationID) {
    const response = await fetch(`https://raw.githubusercontent.com/ycombinator/usta-match-scheduler/refs/heads/main/data/known-events-${organizationID}.json`)
    const json = await response.json()
    return json.events
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
    json.scheduled_events.forEach(event => {
        const eventDate = new Date(event.date)
        const year = eventDate.getFullYear()
        const month = eventDate.getMonth()
        const day = eventDate.getDate()
        const slot = event.slot
        event.id = `${year}_${month}_${day}_${slot}`
    })
    return json
}

function findEventsBounds(events) {
    const now = new Date()
    let eventsStartYear = now.getFullYear()
    let eventsStartMonth = now.getMonth()
    let eventsNumMonths = 0

    if (events.length == 0) {
        return { eventsStartYear, eventsStartMonth, eventsNumMonths }
    }

    events.sort((e1, e2) => Date.parse(e1.date) - Date.parse(e2.date))
    // console.log({events})

    const firstEvent = events[0]
    const firstEventDate = new Date(firstEvent.date)
    eventsStartYear = firstEventDate.getFullYear()
    eventsStartMonth = firstEventDate.getMonth()

    if (events.length > 1) {
        const lastEvent = events[events.length - 1]
        const lastEventDate = new Date(lastEvent.date)
        const eventsEndYear = lastEventDate.getFullYear()
        const eventsEndMonth = lastEventDate.getMonth()

        // console.log({eventsEndYear, eventsStartYear, eventsEndMonth, eventsStartMonth})

        eventsNumMonths = eventsEndMonth - eventsStartMonth
            + 12 * (eventsEndYear - eventsStartYear)        
    }

    eventsNumMonths += 1

    return { eventsStartYear, eventsStartMonth, eventsNumMonths }
}

  function loadEventsFromStorage() {
    try {
      const stored = localStorage.getItem(EVENTS_STORAGE_KEY);
      return stored ? JSON.parse(stored) : [];
    } catch (err) {
      console.error('Failed to load events from localStorage', err);
      return [];
    }
  }


  function saveEventsToStorage(events) {
    try {
      localStorage.setItem(EVENTS_STORAGE_KEY, JSON.stringify(events));
    } catch (err) {
      console.error('Failed to save events to localStorage', err);
    }
  }