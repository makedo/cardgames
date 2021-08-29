import React from 'react';

import Back from "../../../../components/Card/Back";
import Hand from "../../../../components/Card/Hand";

export default function OtherHand({count}) {
  const cards = [];
  for (let i = 0; i < count; i++) {
    cards.push(<Back key={i} />);
  }
  return <Hand>{cards.map(c => c)}</Hand>
}
