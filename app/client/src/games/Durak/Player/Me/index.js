import React from "react";
import {useDrag} from "react-dnd";

import Hand from "../../../../components/Card/Hand";
import Face from "../../../../components/Card/Face";

function DraggableFace({card, ...props}) {
  const [, drag] = useDrag(() => ({ type: 'card', item: {card} }));
  return <Face ref={drag} card={card} {...props} />;
}

export default function Me({me, can_confirm, finished, onConfirm, onRestart}) {
  return <div>
    {me.role}

    {can_confirm && <button onClick={onConfirm}>
      {me.role === 'defender' ? 'Take' : 'Confirm'}
    </button>
    }

    {finished && <button onClick={onRestart}>Restart</button>}

    <p>{me.winner && 'You win!'}</p>
    <p>{me.looser && 'You loose!'}</p>

    <Hand className="my">
      {me.hand.cards && me.hand.cards.length > 0 && me.hand.cards.map(card =>
        <DraggableFace card={card} key={`${card.suite}${card.rank}`} />)}
    </Hand>
  </div>;
}
