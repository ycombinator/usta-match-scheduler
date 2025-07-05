import "./Nav.css"
import { faBackwardStep, faForwardStep } from "@fortawesome/free-solid-svg-icons"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

export function Nav({previous, previousLabel, next, nextLabel}) {
    let previousNav, nextNav
    if (previousLabel) {
        previousNav = (
            <a href="#" className="nav-item" onClick={() => {previous(); return false}}>
                <FontAwesomeIcon icon={faBackwardStep} />
                <span>{previousLabel}</span>
            </a>
        )
    }
    if (nextLabel) {
        nextNav = (
            <a href="#" className="nav-item" onClick={() => {next(); return false}}>
                <span>{nextLabel}</span>
                <FontAwesomeIcon icon={faForwardStep} />
            </a>
        )
    }

    return (
        <span className="nav">
            {previousNav}
            {nextNav}
        </span>
    )

}