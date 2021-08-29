import React from "react";
import {useDrag} from "react-dnd";

import Hand from "../../../../components/Card/Hand";
import Face from "../../../../components/Card/Face";

function DraggableFace({card, ...props}) {
  const [, drag] = useDrag(() => ({ type: 'card', item: {card} }));
  return <Face ref={drag} card={card} {...props} />;
}

export default function MyHand({cards}) {
  return <Hand className="my">
    {cards.map(card =>
      <DraggableFace card={card} key={`${card.suite}${card.rank}`} />)}
  </Hand>;
}
