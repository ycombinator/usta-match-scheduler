import './TextCheckBox.css'

export const TextCheckBox = ({children, isChecked, onClick}) => {
    const classes = isChecked ? "checked" : ""
    return (
        <span onClick={onClick} className={classes}>{children}</span>
    )
}