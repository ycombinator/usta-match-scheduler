import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBackward, faForward } from '@fortawesome/free-solid-svg-icons'
import { getMonthName, getPreviousYearMonth, getNextYearMonth, weeksInMonth } from "../lib/date_utils"
import { CalendarWeek } from "./CalendarWeek"
import "./CalendarMonth.css"

export const CalendarMonth = ({year, month, setStartYearMonth, events}) => {
    const numWeeks = weeksInMonth(year, month)
    const monthName = getMonthName(year, month)

    const calendarWeeks = []
    for(let i = 0; i < numWeeks; i++) {
        const key = year+"_"+month+"_"+i
        calendarWeeks.push(<div key={key}><CalendarWeek year={year} month={month} week={i} events={events} /></div>)
    }

    function goBack(e) {
        e.preventDefault()
        const { prevYear, prevMonth } = getPreviousYearMonth(year, month)
        setStartYearMonth({ year: prevYear, month: prevMonth })
    }

    function goForward(e) {
        e.preventDefault()
        const { nextYear, nextMonth } = getNextYearMonth(year, month)
        setStartYearMonth({ year: nextYear, month: nextMonth })
    }

    return (
        <div className="calendar-month">
            <div className="header">
                <a href="" onClick={goBack}><FontAwesomeIcon icon={faBackward} /></a>
                <h3>{monthName} {year}</h3>
                <a href="" onClick={goForward}><FontAwesomeIcon icon={faForward} /></a>
            </div>
            { calendarWeeks }
        </div>
    )
}