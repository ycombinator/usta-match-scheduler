import React from 'react'
import { toPng } from 'html-to-image';

import './App.css'
import { TeamPreferences } from './components/TeamPreferences'
import { CalendarMonthGroup } from './components/CalendarMonthGroup'
import { Step } from './components/Step'
import { Nav } from './components/Nav'

const asrcOrganizationID = 225

export default class App extends React.Component {
    constructor(props) {
        super(props)
        this.componentRef = React.createRef();

        const now = new Date()
        this.state = {
            appState: "set_team_preferences",
            year: now.getFullYear(),
            month: now.getMonth(),
            numMonths: 1,
            teams: [],
            events: [],
            blackoutEvents: [],
            knownEvents: [],
            isGeneratingSchedule: false,
        }
    }

    async componentDidMount() {
            const knownEvents = await fetchKnownEvents(asrcOrganizationID)
            knownEvents.forEach(event => {
                event.date = new Date(event.date)
            })
            this.setKnownEvents(knownEvents)

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
        const setEvents = events => this.setState({events})
        const setBlackoutEvents = events => this.setState({blackoutEvents: events})
        const setIsGeneratingSchedule = isGeneratingSchedule => this.setState({isGeneratingSchedule})

        // const now = new Date()
        // const [startYearMonth, setStartYearMonth] = useState({year: now.getFullYear(), month: now.getMonth()})

        // const [events, setEvents] = useState([])
        // const [blackoutEvents, setBlackoutEvents] = useState([])
        // const [knownEvents, setKnownEvents] = useState([])

        // const [isGeneratingSchedule, setIsGeneratingSchedule] = useState(false)

        // useEffect(async () => {
        //     const knownEvents = await fetchKnownEvents(asrcOrganizationID)
        //     knownEvents.forEach(event => {
        //         event.date = new Date(event.date)
        //     })
        //     setKnownEvents(knownEvents)
        // }, [])

        // const [teams, setTeams] = useState([]);
        // useEffect(async () => {
        //     const teams = await fetchUpcomingTeams(asrcOrganizationID)
        //     teams.forEach(team => team.day_preferences = [])
        //     setTeams(teams)
        // }, [])

        const changePreferredMatchDays = function(teamIdx, days) {
            const newTeams = structuredClone(self.state.teams)
            newTeams[teamIdx].day_preferences = days
            setTeams(newTeams)
        }

        const setEvent = function(e) {
            const newEvents = []
            let found = false
            for (let i = 0; i < self.state.events.length; i++) {
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
                    allowDeletes={true}
                    header={header}
                    knownEvents={self.state.knownEvents}
                />
                step = 2
                stepLabel = "Set blackout slots"

                navPreviousLabel = "Set team preferences"
                navPrevious = () => setAppState("set_team_preferences")
                navNextLabel = this.state.isGeneratingSchedule ? "Generating..." : "Generate schedule"
                isNextProcessing = this.state.isGeneratingSchedule
                navNext = async () => {
                    setBlackoutEvents(self.state.events)
                    setIsGeneratingSchedule(true)
                    const schedule = await generateSchedule(self.state.teams, self.state.events)
                    setEvents(schedule.scheduled_events)

                    const { eventsStartYear, eventsStartMonth, eventsNumMonths } = findEventsBounds(schedule.scheduled_events)
                    console.log({ eventsStartYear, eventsStartMonth, eventsNumMonths })
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
                            addEventLabel="match"
                            allowAdds={false}
                            allowDeletes={false}
                            knownEvents={self.state.knownEvents}
                        />
                    </div>
                step = 3
                stepLabel = "Review schedule"

                navPreviousLabel = "Set blackout slots"
                navPrevious = () => {
                    setEvents(self.state.blackoutEvents)
                    setAppState("set_blackout_slots")
                }

                navNextLabel = "Print"
                navNext = () => {
                    if (self.componentRef === null) {
                        return
                    }

                    toPng(self.componentRef.current, { cacheBust: true })
                        .then((dataUrl) => {
                            const link = document.createElement('a');
                            link.download = 'schedule.png';
                            link.href = dataUrl;
                            link.click();
                        })
                        .catch((err) => {
                            console.error('Oops, something went wrong!', err);
                        });
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
    console.log({events})

    const firstEvent = events[0]
    const firstEventDate = new Date(firstEvent.date)
    eventsStartYear = firstEventDate.getFullYear()
    eventsStartMonth = firstEventDate.getMonth()

    if (events.length > 1) {
        const lastEvent = events[events.length - 1]
        const lastEventDate = new Date(lastEvent.date)
        const eventsEndYear = lastEventDate.getFullYear()
        const eventsEndMonth = lastEventDate.getMonth()

        console.log({eventsEndYear, eventsStartYear, eventsEndMonth, eventsStartMonth})

        eventsNumMonths = eventsEndMonth - eventsStartMonth
            + 12 * (eventsEndYear - eventsStartYear)        
    }

    eventsNumMonths += 1

    return { eventsStartYear, eventsStartMonth, eventsNumMonths }
}