import { useState } from 'react'
import './App.css'
import { CalendarMonthGroup } from './components/CalendarMonthGroup'

function App() {
    const now = new Date()
    const [startYearMonth, setStartYearMonth] = useState({year: now.getFullYear(), month: now.getMonth()})

    const startYear = startYearMonth.year
    const startMonth = startYearMonth.month

    const events = [
        { start: new Date("2025-04-12T19:00:00Z"), end: new Date("2025-04-12T22:00:00Z"), title: "[W3.5] vs. Morgan Hill Tennis Club"},
        { start: new Date("2025-04-13T16:00:00Z"), end: new Date("2025-04-13T19:00:00Z"), title: "[W3.5DT] vs. Bay Club Courtside"},
        { start: new Date("2025-04-13T19:30:00Z"), end: new Date("2025-04-13T22:30:00Z"), title: "[M4.5] vs. Los Gatos"},
        { start: new Date("2025-04-13T23:00:00Z"), end: new Date("2025-04-14T02:00:00Z"), title: "[M3.5] vs. Bramhall"},
    ]

    return (
        <div className="App">
            <div className="calendar-container">
                <CalendarMonthGroup 
                    startYear={startYear} 
                    startMonth={startMonth} 
                    setStartYearMonth={setStartYearMonth}
                    numMonths={3} 
                    events={events}
                />
            </div>
        </div>
    )
}

export default App