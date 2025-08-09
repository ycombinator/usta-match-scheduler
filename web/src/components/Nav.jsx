import "./Nav.css"
import { faBackwardStep, faForwardStep, faHourglassHalf, faSpinner } from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

export function Nav({previous, previousLabel, isPreviousProcessing, next, nextLabel, isNextProcessing}) {
    let previousNav, nextNav
    if (previousLabel) {
        if (isPreviousProcessing) {
            previousNav = <span>{previousLabel} <FontAwesomeIcon icon={faHourglassHalf} /></span>
        } else {
            previousNav = (
                <a href="#" className="nav-item" onClick={() => {previous(); return false}}>
                    <FontAwesomeIcon icon={faBackwardStep} />
                    <span>{previousLabel}</span>
                </a>
            )
        }
    }
    if (nextLabel) {
        if (isNextProcessing) {
            nextNav =<span>{nextLabel} <FontAwesomeIcon icon={faHourglassHalf} /></span>
        } else {
            nextNav = (
                <a href="#" className="nav-item" onClick={() => {next(); return false}}>
                    <span>{nextLabel}</span>
                    <FontAwesomeIcon icon={faForwardStep} />
                </a>
            )
        }
    }

    return (
        <span className="nav">
            {previousNav}
            {nextNav}
        </span>
    )

}