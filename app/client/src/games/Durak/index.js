import React, { useState } from "react";

import Table from "./Table";

import { useDurakCardDrop, useDurak, onReady, onConfirm, onRestart } from "./hooks";

import "./style.css";
import Me from "./Player/Me";
import OtherPlayer from "./Player/Other";
import Deck from "./Deck";
import Endgame from "./Endgame";

export default function Durak() {

  const [state, setState] = useState(null);
  const [error, setError] = useState(null);
  

  const [{ isOver }, drop] = useDurakCardDrop();
  useDurak(setState, setError);
  
  if (error) {
    return <div style={{color: "red"}}>{error.message}</div>
  }

  if (!state || (!state.me.ready && !state.started)) {
    return <button onClick={onReady}>Ready</button>;
  }

  if (state.me.ready && !state.started) {
    return "Waiting for other players to ready...";
  }

  return <div className="durak" >
    <div className="header">
      <OtherPlayer player={state.players[0]} />
    </div>

    <div
      className="game-table" 
      ref={drop}
      style={isOver ? { "border": "2px solid black" } : {}} 
    >
      <Deck trump_card={state.trump_card} amount={state.deck_amount} />
      <Table table={state.table || {}} />
      {state.me.looser && "LOOOSER"}
      {state.me.winner && "WINNER"}
    </div>

    <div className="footer">
      <Me player={state.me || {}} can_confirm={state.can_confirm || false} onConfirm={onConfirm} />
    </div>

    <Endgame onClose={onRestart} me={state.me} />
  </div>;
}

