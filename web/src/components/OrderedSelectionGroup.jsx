import "./OrderedSelectionGroup.css"
import { OrderedSelectionItem } from "./OrderedSelectionItem"

// allItems = [ "Sunday", "Monday": "Tuesday" ]
// selectedItems = [ "Monday", "Sunday" ]
export const OrderedSelectionGroup = ({allItems, selectedItems, setSelectedItems}) => {
    let  sItems = [...selectedItems]
    const selectItem = item => sItems.push(item)
    const deselectItem = item => sItems = sItems.filter(i => i != item)
    const toggleItem = item => {
        sItems.includes(item) ? deselectItem(item) : selectItem(item)
        setSelectedItems(sItems)
        return false
    }

    const items = allItems.map(item => (<OrderedSelectionItem label={item} order={selectedItems.indexOf(item)} onClick={toggleItem} />))
    // const items = selectedItems.map(item => (<OrderedSelectionItem label={item} order={selectedItems.indexOf(item)} onClick={toggleItem} />))
    //     .concat(allItems.filter(i => !selectedItems.includes(i)).map(item => (<OrderedSelectionItem label={item} onClick={toggleItem} />)))
    return (<div className="ordered-selection-group">{items}</div>)
}