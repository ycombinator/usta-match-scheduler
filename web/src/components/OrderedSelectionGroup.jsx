import "./OrderedSelectionGroup.css"
import { OrderedSelectionItem } from "./OrderedSelectionItem"

// allItems = [ "Sunday", "Monday": "Tuesday" ]
// selectedItems = [ 1, 0 ] // Monday, Sunday
export const OrderedSelectionGroup = ({allItems, unselectableItems, selectedItems, setSelectedItems}) => {
    let  sItems = [...selectedItems]
    const selectItem = item => sItems.push(item)
    const deselectItem = item => sItems = sItems.filter(i => i != item)
    const toggleItem = item => {
        sItems.includes(item) ? deselectItem(item) : selectItem(item)
        setSelectedItems(sItems)
        return false
    }

    // console.log({unselectableItems})

    const items = allItems.map((item, idx) => (
        <OrderedSelectionItem isSelectable={!unselectableItems.includes(item)} label={item} value={idx} order={selectedItems.indexOf(idx)} onClick={toggleItem} />
    ))
    // const items = selectedItems.map(item => (<OrderedSelectionItem label={item} order={selectedItems.indexOf(item)} onClick={toggleItem} />))
    //     .concat(allItems.filter(i => !selectedItems.includes(i)).map(item => (<OrderedSelectionItem label={item} onClick={toggleItem} />)))
    return (<div className="ordered-selection-group">{items}</div>)
}