import { isSameDay, isSameMonth, isWeekendDay } from "../lib/date_utils"
import { CalendarEvent } from "./CalendarEvent"
import "./CalendarDay.css"
import { useState } from "react"

export const CalendarDay = ({thisYear, thisMonth, year, month, day, events, setEvent, addEventLabel}) => {
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

    // FIXME? extract out of component since it's not generic enough for component?
    const allSlots = isWeekend ? ["morning", "afternoon", "evening"] : ["morning", "evening"]
    addEventLabel = addEventLabel || "event"

    const [addEventIdx, setAddEventIdx] = useState(-1)
    const [addEventText, setAddEventText] = useState("")

    const submitAddEvent = (slot) => {
        const title = addEventText.trim()
        if (title != "") {
            const id = `${year}${month}${day}_${slot}`
            setEvent({id: id, type: addEventLabel, slot: slot, start: currentDay, end: currentDay, title: title});
        }
        setAddEventIdx(-1)
        setAddEventText("")
    }

    // TODO: don't show add event buttons in generated schedule?

    // Mix events + remaining slots and order results by slot
    let items = []
    allSlots.forEach((slot, i) => {
        const slotEvents = events.filter(e => e.slot == slot)
        if (slotEvents.length > 0) {
            // There is an event in this slot
            const slotEvent = slotEvents[0]
            items.push(
                <li className="calendar-day-event" key={i}>
                    <CalendarEvent year={year} month={month} day={day} event={slotEvent} setEvent={setEvent}/>
                </li>
            )
        } else if (i == addEventIdx) {
            // An event is being added to this slot
            items.push(
                <li className="calendar-day-event" key={i}>
                    <form onBlur={() => {submitAddEvent(slot); return false;}} onSubmit={() => {submitAddEvent(slot); return false;}}>
                        <input type="text" autoFocus={true} value={addEventText} onChange={e => setAddEventText(e.target.value)}></input>
                    </form>
                </li>
            )
        } else {
            const className = `calendar-event new`
            items.push(
                <li className="calendar-day-event" key={i}>
                    <p className={className}>
                        <a href="#" onClick={() => {setAddEventIdx(i); return false;}}>add {slot} {addEventLabel}</a>
                    </p>
                </li>
            )
        }
    })
    return (
        <div className="calendar-day">
            <h4 className={dayClass}>{day}</h4>
            <ol>{items}</ol>
        </div>
    )
}
    