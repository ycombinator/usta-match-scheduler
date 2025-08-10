import React, {useState} from 'react'
import { DndContext } from '@dnd-kit/core';
import { isEventInMonth, nextMonth, previousMonth } from "../lib/date_utils"
import { CalendarMonth } from "./CalendarMonth"
import { Droppable } from "./Droppable"
import "./CalendarMonthGroup.css"

export class CalendarMonthGroup extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            draggingMatch: null,
            droppedID: null,
        }
    }

    setDraggingMatch = (draggingMatch) => this.setState({draggingMatch})
    setDroppedID = (droppedID) => this.setState({droppedID})

    handleDragEnd = (event) => {
        if (event.over) {
            // console.log(`Dropped event with ID = ${event.active.id} on ID = ${event.over.id}`)
            this.props.moveEvent(event.active.id, event.over.id)
            this.setDroppedID(event.over.id)
            this.setDraggingMatch(null);
        }
    }

    handleDragStart = (event) => {
        const fromID = event.active.id
        this.props.events.forEach(event => {
            if (event.id == fromID) {
                this.setDraggingMatch(event)
            }
        })
    }

    render() {
        const {startYear, startMonth, setStartYearMonth, numMonths, events, setEvent, addEventLabel, allowAdds, allowEdits, allowDeletes, allowMoves, header, knownEvents} = this.props
        // console.log("calendar month group: ", events)
        const months = []
        let year = startYear
        let month = startMonth
        for (let i = 0; i < numMonths; i++) {
            // Include events from previous month, current month, and next month so display
            // works correctly
            const monthEvents = events.filter(monthEventFilter(year, month))
            // console.log({knownEvents})
            const monthKnownEvents = knownEvents.filter(monthEventFilter(year, month))
            // console.log({monthKnownEvents})
            months.push(
                <div key={i} className="calendar-month-container">
                    <CalendarMonth
                        year={year} month={month}
                        setStartYearMonth={setStartYearMonth}
                        events={monthEvents} setEvent={setEvent} addEventLabel={addEventLabel}
                        allowAdds={allowAdds} allowEdits={allowEdits} allowDeletes={allowDeletes} allowMoves={allowMoves}
                        knownEvents={monthKnownEvents}
                        draggingMatch={this.state.draggingMatch}
                    />
                </div>
            )

            month++
            // Check if we should start the new year
            if (month == 12) {
                year++
                month = 0
            }
        }

        return (
            <div>
                <DndContext onDragEnd={this.handleDragEnd} onDragStart={this.handleDragStart}>
                    { header }
                    <div className="calendar-month-group">
                        { months }
                    </div>
                </DndContext>
            </div>
        )
    }
}

function monthEventFilter(year, month) {
    return (event) => {
        return isEventInMonth(year, previousMonth(month), event) ||
            isEventInMonth(year, month, event) ||
            isEventInMonth(year, nextMonth(month), event)
    }
}