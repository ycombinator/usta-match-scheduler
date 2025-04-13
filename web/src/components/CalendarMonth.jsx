import { getMonthName, weeksInMonth } from "../lib/date_utils"
import { CalendarWeek } from "./CalendarWeek"
import "./CalendarMonth.css"

export const CalendarMonth = ({year, month, events}) => {
    const numWeeks = weeksInMonth(year, month)
    const monthName = getMonthName(year, month)

    const calendarWeeks = []
    for(let i = 0; i < numWeeks; i++) {
        const key = year+"_"+month+"_"+i
        calendarWeeks.push(<div key={key}><CalendarWeek year={year} month={month} week={i} events={events} /></div>)
    }

    return (
        <div className="calendar-month">
            <h3>{monthName} {year}</h3>
            { calendarWeeks }
        </div>
    )
}