import './OrderedSelectionItem.css'

export const OrderedSelectionItem = ({isSelectable, label, value, order, onClick}) => {
    // console.log({label, isSelectable})
    if (!isSelectable) {
        return (
            <div className="ordered-selection-item-unselectable">
                <span className="order">{ order >= 0 ? order+1 : "" }</span>
                {label}
            </div>
        )
    }

    const isChecked = order >= 0
    const classes = isChecked ? "ordered-selection-item checked" : "ordered-selection-item unchecked"
    return (
        <div onClick={(e) => { e.stopPropagation(); onClick(value)}} className={classes}>
            <span className="order">{ order >= 0 ? order+1 : "" }</span>
            {label}
        </div>
    )
}