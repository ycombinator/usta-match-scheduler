import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBackward, faForward } from '@fortawesome/free-solid-svg-icons'
import { getMonthName, getPreviousYearMonth, getNextYearMonth, weeksInMonth } from "../lib/date_utils"
import { CalendarWeek } from "./CalendarWeek"
import "./CalendarMonth.css"
import "./CalendarWeek.css"

export const CalendarMonth = ({year, month, setStartYearMonth, events, setEvent, addEventLabel, allowAdds, allowDeletes, knownEvents}) => {
    console.log("calendar month: ", events)
    const numWeeks = weeksInMonth(year, month)
    const monthName = getMonthName(year, month)

    const weekdayNames = ["Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"]
        .map(name => (
            (name === "Sunday" || name === "Saturday")
            ? <div className="calendar-weekday weekend">{name}</div>
            : <div className="calendar-weekday">{name}</div>
        ))

    const calendarWeeks = []
    for(let i = 0; i < numWeeks; i++) {
        const key = year+"_"+month+"_"+i
        calendarWeeks.push(
            <div key={key}>
                <CalendarWeek
                    year={year} month={month} week={i}
                    events={events} setEvent={setEvent} addEventLabel={addEventLabel}
                    allowAdds={allowAdds} allowDeletes={allowDeletes}
                    knownEvents={knownEvents}
                />
            </div>
        )
    }

    function goBack(e) {
        e.preventDefault()
        const { prevYear, prevMonth } = getPreviousYearMonth(year, month)
        setStartYearMonth(prevYear, prevMonth)
    }

    function goForward(e) {
        e.preventDefault()
        const { nextYear, nextMonth } = getNextYearMonth(year, month)
        setStartYearMonth(nextYear, nextMonth)
    }

    return (
        <div className="calendar-month">
            <div className="header">
                <a href="#" onClick={goBack}><FontAwesomeIcon icon={faBackward} /></a>
                <h3>{monthName} {year}</h3>
                <a href="#" onClick={goForward}><FontAwesomeIcon icon={faForward} /></a>
            </div>
            <div className="calendar-week">
                {weekdayNames}
            </div>
            { calendarWeeks }
        </div>
    )
}