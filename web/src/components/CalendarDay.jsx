import { isSameDay, isSameMonth, isWeekendDay } from "../lib/date_utils"
import { CalendarEvent } from "./CalendarEvent"
import "./CalendarDay.css"

export const CalendarDay = ({thisYear, thisMonth, year, month, day, events}) => {
    const currentDay = new Date(year, month, day)
    const today = new Date()
    const isToday = isSameDay(today, currentDay)
    const isInThisMonth = isSameMonth(new Date(thisYear, thisMonth, 1), currentDay)
    const isWeekend = isWeekendDay(currentDay)

    let dayClass = ""
    if (isInThisMonth) {
        if (isToday) {
            dayClass = "today"
        } else if (isWeekend) {
            dayClass = "weekend"
        }
    } else {
        dayClass = "not-this-month"
    }

    // TODO: add ability to add events
    // - two slots on weekdays: morning, evening
    // - three slots on weekends: morning, afternoon, evening
    return (
        <div className="calendar-day">
            <h4 className={dayClass}>{day}</h4>
            <ol>
            { 
                events.map((event, i) => {
                    return (
                        <li className="calendar-day-event" key={i}><CalendarEvent year={year} month={month} day={day} event={event} /></li>
                    )
                })
            }
            </ol>
        </div>
    )
}
    