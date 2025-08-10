import { useDraggable } from '@dnd-kit/core';
import "./Draggable.css"

export const Draggable = (props) => {
  const {attributes, listeners, setNodeRef, transform} = useDraggable({
    id: props.id,
  });
  const style = transform ? {
    transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
    zIndex: 10,
    position: `absolute`,
  } : undefined;

  
  return (
    <button ref={setNodeRef} style={style} {...listeners} {...attributes} className="draggable">
      {props.children}
    </button>
  );
}