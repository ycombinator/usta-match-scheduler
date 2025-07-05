import { useState } from "react"
import { doesEventEndInDay, doesEventStartInDay, getPaddedTime } from "../lib/date_utils"
import "./CalendarEvent.css"
import { faTrash } from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

export const CalendarEvent = ({year, month, day, event, setEvent, allowDelete}) => {
    const start = doesEventStartInDay(year, month, day, event) 
        ? getPaddedTime(event.start)
        : "..."
    const end = doesEventEndInDay(year, month, day, event)
        ? getPaddedTime(event.end)
        : "..."
    const className = `calendar-event ${event.type}`

    const [editEvent, setEditEvent] = useState(false)
    const [editEventText, setEditEventText] = useState(event.title)

    const submitEditEvent = () => {
        setEditEvent(false)
        const title = editEventText.trim()
        setEvent({id: event.id, type: event.type, slot: event.slot, start: event.start, end: event.end, title: title}); ;
    }

    let deleteButton, deleteEvent
    if (allowDelete) {
        deleteButton = <FontAwesomeIcon icon={faTrash} onClick={() => {deleteEvent(); return false;}}/>
        deleteEvent = () => {
            setEvent({id: event.id, type: event.type, slot: event.slot, start: event.start, end: event.end, title: ""}); ;
        }
    }

    if (editEvent) {
        return (
            <form onBlur={() => {submitEditEvent(); return false;}} onSubmit={() => {submitEditEvent(); return false;}}>
                <input type="text" autoFocus={true} value={editEventText} onChange={e => setEditEventText(e.target.value)}></input>
            </form>
        )
    }

    // TODO: show trash icon for deleting

    return (
        <span className={className}>
            <span onClick={() => {setEditEvent(true); return false;}}>
                {getSlotLabel(event.slot)}: {event.title}
                {/* {start}-{end}: {title} */}
            </span>
            {deleteButton}
        </span>
    )
}

function getSlotLabel(slot) {
    switch (slot) {
        case "morning":   return "Morn."
        case "afternoon": return "Aft."
        case "evening":   return "Eve."
        default:          return slot.substring(0,3)
    }
}

// Event
// {
//   "type": "match", // or "blackout"
//   "slot": "evening", // or "morning" or "afternoon"
//   "start": "2025-04-13T17:00:00Z",
//   "end": "2025-04-13T19:00:00Z",
//   "id": "a3e59ac3",
//   "title": "[M3.5] vs. Bramhall",
// }