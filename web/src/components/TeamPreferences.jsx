import "./TeamPreferences.css"
import { TextCheckBox } from "./TextCheckBox"

const dayOfWeekMap = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]

export const TeamPreferences = ({teams, changeDayPreference}) => {
    teams = teams.map((team, teamIdx)=> {
        const days = []
        for (let i = 0; i < 7; i++) {
            days.push(<TextCheckBox isChecked={team.day_preferences[i]} onClick={() => changeDayPreference(teamIdx, i)}>{dayOfWeekMap[i]}</TextCheckBox>)
        }

        return (
            <tr key={team.id}>
                <td>{team.name}</td>
                <td>{team.captain}</td>
                <td className="days">{days}</td>
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
