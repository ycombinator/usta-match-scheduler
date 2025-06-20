import { useState } from 'react'
import { OrderedSelectionGroup } from "./OrderedSelectionGroup"
import "./TeamPreferences.css"

const dayOfWeekMap = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]

export const TeamPreferences = ({teams, changeDayPreference}) => {
    teams = teams.map((team, teamIdx)=> {
        const allDays = []
        for (let i = 0; i < 7; i++) {
            allDays.push(dayOfWeekMap[i])
        }

        const [selectedDays, setSelectedDays] = useState([])

        return (
            <tr key={team.id}>
                <td>{team.name}</td>
                <td>{team.captain}</td>
                <td className="days"><OrderedSelectionGroup allItems={allDays} selectedItems={selectedDays} setSelectedItems={setSelectedDays} /></td>
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
