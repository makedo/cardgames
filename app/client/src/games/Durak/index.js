import React from "react";

import Face from "../../components/Card/Face";
import Back from "../../components/Card/Back";

import MyHand from "./Hand/My";
import OtherHand from "./Hand/Other";
import Table from "./Table";

import { useDurakCardDrop, useDurakState, useDurakWebsockets } from "./hooks";

import "./style.css";

export default function Durak() {

  const sendMessage = useDurakWebsockets();

  const ready = function() {
      sendMessage('{"type": "ready", "data": {"playerId": "abc"}}');
  }

  const createOnCardDrop = (val) => (card) => {
    console.log(val);
    console.log(card);
  }
  const [{ isOver }, drop] = useDurakCardDrop(createOnCardDrop('ON TABLE'));
  const [state] = useDurakState();

  if (!state) {
    return "Loading...";
  }

  return <div className="durak">
    <div className="header">
      <OtherHand count={state.hands[0]} />
    </div>

    <div ref={drop} className="game-table" style={isOver ? { "border": "2px solid black" } : {}}>
      <div className="deck">
        <Back className="top-card" />
        <Face className="main-suite" card={state.trump_card} />
      </div>

      <Table cards={state.table.cards || []} onCardDrop={createOnCardDrop('ON CARD')} />
    </div>

    <div className="footer">
      <button onClick={ready}>Ready</button>
      <MyHand cards={state.hand} />
    </div>
  </div>;
}

