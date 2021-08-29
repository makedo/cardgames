import React from "react"
import Face from "../../../components/Card/Face";
import { useDurakCardDrop } from "../hooks";

  
function DroppableFace({ onCardDrop, ...props }) {
    const [{ isOver }, drop] = useDurakCardDrop(onCardDrop);
    return <Face
      ref={drop}
      style={isOver ? { "border": "2px solid black" } : {}}
      {...props}
    />
}

export default function Table({ cards: cardsByPlace, onCardDrop }) {
    const renderCards = function (cards) {
      const CardComponent = cards.length === 2 ? Face : DroppableFace;
      const classes = ['bottom', 'top'];
      
      return <div className="card-placeholder">
        {cards.map((card, key) =>
          <CardComponent
            key={card}
            onCardDrop={onCardDrop}
            card={card}
            className={classes[key]}
          />)
        }</div>;
    }
  
    return <div className="table">{cardsByPlace.map(renderCards)}</div>
  }