import React from 'react';

import Back from "../../../../components/Card/Back";
import Hand from "../../../../components/Card/Hand";

export default function OtherPlayer({player}) {
  const cards = [];
  const cardsCount = player.hand || 0;

  for (let i = 0; i < player.hand; i++) {
    cards.push(<Back key={i} />);
  }
  return <div>
    <Hand>{cards.map(c => c)}</Hand>
    {player.state}
  </div>
}
