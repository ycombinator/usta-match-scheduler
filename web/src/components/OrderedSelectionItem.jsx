import './OrderedSelectionItem.css'

export const OrderedSelectionItem = ({label, order, onClick}) => {
    console.log({label, order})
    const isChecked = order >= 0
    const classes = isChecked ? "ordered-selection-item checked" : "ordered-selection-item unchecked"
    return (
        <div onClick={(e) => { e.stopPropagation(); onClick(label)}} className={classes}>
            <span className="order">{ order >= 0 ? order+1 : "" }</span>
            {label}
        </div>
    )
}