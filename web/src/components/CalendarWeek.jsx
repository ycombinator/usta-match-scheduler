import { daysInMonth, getPreviousMonth, getPreviousMonthYear, getNextMonth, getNextMonthYear, isEventInDay, monthDaysInFirstWeek } from "../lib/date_utils"
import { CalendarDay } from "./CalendarDay"
import "./CalendarWeek.css"

export const CalendarWeek = ({year, month, week, events}) => {
    const thisYear = year
    const thisMonth = month

    const firstDayOfMonth = new Date(thisYear, thisMonth, 1)
    const firstDayOfMonthWeekday = firstDayOfMonth.getDay()

    let nextMonthDays = 0
    const days = []
    for (let i = 0; i < 7; i++) {
        let year = thisYear
        let month = thisMonth
        let day = 0

        if (week == 0) {
            if (i < firstDayOfMonthWeekday) {
                month = getPreviousMonth(thisMonth)
                year = getPreviousMonthYear(thisYear, thisMonth)

                const dayDiff = firstDayOfMonthWeekday - i - 1
                day = daysInMonth(year, month) - dayDiff
            } else {
                day = i - firstDayOfMonthWeekday + 1
            }
        } else {
            const startDayOffset = ((week - 1) * 7) + monthDaysInFirstWeek(firstDayOfMonthWeekday)
            day = startDayOffset + i + 1

            if (day > daysInMonth(year, month)) {
                day = 1 + nextMonthDays
                nextMonthDays++ 

                month = getNextMonth(thisMonth)
                year = getNextMonthYear(thisYear, thisMonth)
            }
        }

        const dayEvents = events.filter(event => isEventInDay(year, month, day, event))
        const key = year+"_"+month+"_"+day
        days.push(<div key={key}><CalendarDay thisYear={thisYear} thisMonth={thisMonth} year={year} month={month} day={day} events={dayEvents} /></div>)
    }

    return (
        <div className="calendar-week">
            { days }
        </div>
        
    )
}