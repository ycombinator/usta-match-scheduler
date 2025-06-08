import './Step.css'

export const Step = ({current, total, label}) => {
    label = label ? ": " + label : ""
    return (
        <span>Step {current} of {total}{label}</span>
    )
}