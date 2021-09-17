import React, { useState } from "react";

import Face from "../../components/Card/Face";
import Back from "../../components/Card/Back";

import MyHand from "./Hand/My";
import OtherHand from "./Hand/Other";
import Table from "./Table";

import { useDurakCardDrop, useDurak, onReady } from "./hooks";

import "./style.css";

export default function Durak() {

  const [state, setState] = useState(null);
  const [error, setError] = useState(null);
  const [isReady, setIsReady] = useState(false);
  

  const [{ isOver }, drop] = useDurakCardDrop();
  useDurak(setState, setError);
  
  if (error) {
    return <div style={{color: "red"}}>{error.message}</div>
  }

  if (!state && !isReady) {
    const onReadyClick = onReady(() => setIsReady(true))
    return <button onClick={onReadyClick}>Ready</button>;
  }

  if (!state && isReady) {
    return "Waiting for other players to ready...";
  }

  return <div className="durak" >
    <div className="header">
      <OtherHand count={state.hands[0] || 0} />
    </div>

    <div
      className="game-table" 
      ref={drop}
      style={isOver ? { "border": "2px solid black" } : {}} 
    >
      <div className="deck">
        <Back className="top-card" />
        <Face className="main-suite" card={state.trump_card} />
      </div>
      <Table cards={state.table.cards || []} />
    </div>

    <div className="footer">
      <MyHand cards={state.hand || []} />
    </div>
  </div>;
}

