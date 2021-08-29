import React from 'react';
import Suite from "../Suite";
import Rank from "../Rank";

import './style.css';

function parseCard(card) {
  return [card.suite, card.rank];
}

const Face = React.forwardRef(({card, className, style}, ref) => {
  const [suite, rank] = parseCard(card);

  return <div ref={ref} style={style} className={["card", className].filter((c) => !!c).join(' ')} >
    <div className="content">
      <div className="top">
        <Rank rank={rank} />
        <Suite suite={suite} />
      </div>
      <div className="middle">
        <h1><Rank rank={rank} /><Suite suite={suite} /></h1>
      </div>
      <div className="bottom"><Rank rank={rank} /><Suite suite={suite} /></div>
    </div>
  </div>
});

export default Face;
