import { daysInMonth, isEventInDay, monthDaysInFirstWeek, getPreviousYearMonth, getNextYearMonth } from "../lib/date_utils"
import { CalendarDay } from "./CalendarDay"
import "./CalendarWeek.css"

export const CalendarWeek = ({year, month, week, events, setEvent, addEventLabel, allowAdds, allowDeletes, knownEvents}) => {
    console.log("calendar week: ", events)
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
                const { prevYear, prevMonth } = getPreviousYearMonth(thisYear, thisMonth)
                year = prevYear
                month = prevMonth

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

                const { nextYear, nextMonth } = getNextYearMonth(thisYear, thisMonth)
                year = nextYear
                month = nextMonth
            }
        }

        const dayEvents = events.filter(dayEventFilter(year, month, day))
        const dayKnownEvents = knownEvents.filter(dayEventFilter(year, month, day))
        const key = year+"_"+month+"_"+day
        days.push(
            <div key={key}>
                <CalendarDay
                    thisYear={thisYear} thisMonth={thisMonth}
                    year={year} month={month} day={day}
                    events={dayEvents} setEvent={setEvent} addEventLabel={addEventLabel}
                    allowAdds={allowAdds} allowDeletes={allowDeletes}
                    knownEvents={dayKnownEvents}
                />
            </div>
        )
    }

    return (
        <div className="calendar-week">
            { days }
        </div>

    )
}

function dayEventFilter(year, month, day) {
    return (event) => {
        return isEventInDay(year, month, day, event)
    }
}