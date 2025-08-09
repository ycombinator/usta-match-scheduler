import React from 'react'
import { isEventInMonth, nextMonth, previousMonth } from "../lib/date_utils"
import { CalendarMonth } from "./CalendarMonth"
import "./CalendarMonthGroup.css"

export class CalendarMonthGroup extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        const {startYear, startMonth, setStartYearMonth, numMonths, events, setEvent, addEventLabel, allowAdds, allowDeletes, header, knownEvents} = this.props
        console.log("calendar month group: ", events)
        const months = []
        let year = startYear
        let month = startMonth
        for (let i = 0; i < numMonths; i++) {
            // Include events from previous month, current month, and next month so display
            // works correctly
            const monthEvents = events.filter(monthEventFilter(year, month))
            console.log({knownEvents})
            const monthKnownEvents = knownEvents.filter(monthEventFilter(year, month))
            console.log({monthKnownEvents})
            months.push(
                <div key={i} className="calendar-month-container">
                    <CalendarMonth 
                        year={year} month={month} 
                        setStartYearMonth={setStartYearMonth}
                        events={monthEvents} setEvent={setEvent} addEventLabel={addEventLabel}
                        allowAdds={allowAdds} allowDeletes={allowDeletes}
                        knownEvents={monthKnownEvents}
                    />
                </div>
            )

            month++
            // Check if we should start the new year
            if (month == 12) {
                year++
                month = 0
            }
        }

        return (
            <div>
            { header }
                <div className="calendar-month-group">
                    { months }
                </div>
            </div>
        )
    }
}

function monthEventFilter(year, month) {
    return (event) => {
        return isEventInMonth(year, previousMonth(month), event) ||
            isEventInMonth(year, month, event) ||
            isEventInMonth(year, nextMonth(month), event)
    }
}