import './TextCheckBox.css'

export const TextCheckBox = ({children, isChecked, onClick}) => {
    const classes = isChecked ? "textcheckbox checked" : "textcheckbox unchecked"
    return (
        <span onClick={onClick} className={classes}>{children}</span>
    )
}