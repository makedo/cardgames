import React from "react"
import Face from "../../../components/Card/Face";
import { useDurakCardDrop } from "../hooks";

  
function DroppableFace({place, ...props }) {
    const [{ isOver }, drop] = useDurakCardDrop(place);
    return <Face
      ref={drop}
      style={isOver ? { "border": "2px solid black" } : {}}
      {...props}
    />
}

const Table = ({ cards: cardsByPlace }) => {
    const renderCards = function (cards, place) {
      const CardComponent = cards.length === 2 ? Face : DroppableFace;
      const classes = ['bottom', 'top'];
      
      return <div key={place} className="card-placeholder">
        {cards.map((card, key) =>
          <CardComponent
            key={`${card.suite}${card.rank}`} 
            className={classes[key]}
            card={card}
            place={place}
          />)
        }</div>;
    }
  
    return <div className="table">{cardsByPlace.map(renderCards)}</div>
  }

  export default Table