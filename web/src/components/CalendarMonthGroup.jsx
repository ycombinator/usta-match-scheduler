import { isEventInMonth } from "../lib/date_utils"
import { CalendarMonth } from "./CalendarMonth"
import "./CalendarMonthGroup.css"

export const CalendarMonthGroup = ({startYear, startMonth, setStartYearMonth, numMonths, events, addEventLabel}) => {
    const months = []
    let year = startYear
    let month = startMonth
    for (let i = 0; i < numMonths; i++) {
        const monthEvents = events.filter(event => isEventInMonth(year, month, event))
        months.push(
            <div key={i} className="calendar-month-container">
                <CalendarMonth 
                    year={year} month={month} 
                    setStartYearMonth={setStartYearMonth}
                    events={monthEvents} addEventLabel={addEventLabel}
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
        <div className="calendar-month-group">
            { months }
        </div>
    )
}