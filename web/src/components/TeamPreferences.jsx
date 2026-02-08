import { useState } from 'react'
import { OrderedSelectionGroup } from "./OrderedSelectionGroup"
import "./TeamPreferences.css"

const dayOfWeekMap = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]

export const TeamPreferences = ({teams, changePreferredMatchDays, changeWeight}) => {
    teams = teams.map((team, teamIdx)=> {
        const allDays = []
        for (let i = 0; i < 7; i++) {
            allDays.push(dayOfWeekMap[i])
        }

        const unselectableDays = (team.schedule_group == "Daytime") ? ["Saturday", "Sunday"] : []

        const setSelectedMatchDays = days => {
            // console.log({days})
            changePreferredMatchDays(teamIdx, days)
        }
        const teamType = team.type == "Adult"?  "" : team.type
        const daytime = team.schedule_group == "Daytime" ? "Daytime" : ""

        return (
            <tr key={team.id}>
                <td>{teamType} {team.gender} {team.min_age}+ {team.level} {daytime} <span class="team-name">{team.name}</span></td>
                <td>{team.captain}</td>
                <td className="days">
                    <OrderedSelectionGroup
                        allItems={allDays}
                        unselectableItems={unselectableDays}
                        selectedItems={team.day_preferences}
                        setSelectedItems={setSelectedMatchDays}
                    />
                </td>
                <td>
                    <input
                        type="number"
                        min="0"
                        value={team.scheduling_weight || ''}
                        placeholder="0"
                        onChange={e => changeWeight(teamIdx, parseInt(e.target.value) || 0)}
                        style={{width: '60px'}}
                    />
                </td>
            </tr>
        )
    })

    return (
        <table>
            <thead>
                <th>Team</th>
                <th>Captain</th>
                <th>Preferred Match Days</th>
                <th>Scheduling Weight</th>
            </thead>
            <tbody>
                {teams}
            </tbody>
        </table>
    )
}
