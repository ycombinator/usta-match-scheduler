import { getPaddedTime, isEventInDay } from "../lib/date_utils"
import "./CalendarEvent.css"

export const CalendarEvent = ({year, month, day, event}) => {
    const start = isEventInDay(year, month, day, event) 
        ? getPaddedTime(event.start)
        : ""
    const end = isEventInDay(year, month, day, event)
        ? getPaddedTime(event.end)
        : ""
    const title = event.title

    return (
        <p className="calendar-event">
            {start}-{end}: {title}
        </p>
    )
}

// Event
// {
//   "start": "2025-04-13T17:00:00Z",
//   "end": "2025-04-13T19:00:00Z",
//   "id": "a3e59ac3",
//   "title": "[M3.5] vs. Bramhall",
// }