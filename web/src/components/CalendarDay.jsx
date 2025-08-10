import { isSameDay, isSameMonth, isWeekendDay } from "../lib/date_utils"
import { CalendarEvent } from "./CalendarEvent"
import "./CalendarDay.css"
import { useState } from "react"
import { useDraggable } from '@dnd-kit/core';
import { Droppable } from "./Droppable";
import { Draggable } from "./Draggable";

export const CalendarDay = ({thisYear, thisMonth, year, month, day, events, setEvent, addEventLabel, allowAdds, allowEdits, allowDeletes, allowMoves, knownEvents, draggingID}) => {
    // console.log("calendar day: ", events)
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
            const id = `${year}_${month}_${day}_${slot}`
            // console.log({currentDay, slot, id, addEventLabel, title})
            setEvent({id: id, type: addEventLabel, slot: slot, date: currentDay, title: title});
        }
        setAddEventIdx(-1)
        setAddEventText("")
    }

    // TODO: don't show add event buttons in generated schedule?

    // Mix events + remaining slots and order results by slot
    let items = []
    allSlots.forEach((slot, i) => {
        const id = `${year}_${month}_${day}_${slot}`
        // console.log({id})
        const slotEvents = events.filter(e => e.slot == slot)
        if (slotEvents.length > 0) {
            // There is an event in this slot
            const slotEvent = slotEvents[0]

            items.push(
                <li className="calendar-day-event"  key={i}>
                    <Draggable id={id}>
                        <CalendarEvent year={year} month={month} day={day} event={slotEvent} setEvent={setEvent} allowEdit={allowEdits} allowDelete={allowDeletes} draggingID={draggingID} />
                    </Draggable>
                </li>
            )
        } else if (i == addEventIdx) {
            // An event is being added to this slot
            items.push(
                <li className="calendar-day-event" key={i}>
                    <form onBlur={() => {submitAddEvent(slot); return false;}} onSubmit={() => {submitAddEvent(slot); return false;}}>
                        <input type="text" autoFocus={true} value={addEventText} onChange={e => setAddEventText(e.target.value)} placeholder="enter event title"></input>
                    </form>
                </li>
            )
        } else if (allowAdds) {
            const className = `calendar-event new`
            items.push(
                <li className="calendar-day-event" key={i}>
                    <p className={className} onClick={() => {setAddEventIdx(i); return false;}}>add {slot} {addEventLabel}</p>
                </li>
            )
        } else if (allowMoves) {
            items.push(
                <Droppable id={id}>
                    {draggingID ? "Dropped!" : slot}
                </Droppable>
            )
        }
    })

    knownEvents = knownEvents.map(event => (
        <span className="known-event">{event.title}</span>
    ))

    return (
        <div className="calendar-day">
            <div className="header">
                <h4 className={dayClass}>{day}</h4>
                <span className="known-events">{knownEvents}</span>
            </div>
            <ol>{items}</ol>
        </div>
    )
}
