import { useState } from 'react'
import { OrderedSelectionGroup } from "./OrderedSelectionGroup"
import "./TeamPreferences.css"

const dayOfWeekMap = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]

export const TeamPreferences = ({teams, changePreferredMatchDays}) => {
    teams = teams.map((team, teamIdx)=> {
        const allDays = []
        for (let i = 0; i < 7; i++) {
            allDays.push(dayOfWeekMap[i])
        }

        const setSelectedMatchDays = days => changePreferredMatchDays(teamIdx, days)
        const teamType = team.type == "Adult"?  "" : team.type

        return (
            <tr key={team.id}>
                <td>{teamType} {team.gender} {team.min_age}+ {team.level} <span class="team-name">{team.name}</span></td>
                <td>{team.captain}</td>
                <td className="days"><OrderedSelectionGroup allItems={allDays} selectedItems={team.preferred_match_days} setSelectedItems={setSelectedMatchDays} /></td>
            </tr>
        )
    })

    return (
        <table>
            <thead>
                <th>Team</th>
                <th>Captain</th>
                <th>Preferred Match Days</th>
            </thead>
            <tbody>
                {teams}
            </tbody>
        </table>
    )
}
