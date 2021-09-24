import React from "react"
import ReactModal from "react-modal"

export default function Endgame({onClose, me}) {
    return <ReactModal
        isOpen={me.winner || me.looser}
        onRequestClose={onClose}
        shouldCloseOnEsc={false}
        shouldCloseOnOverlayClick={true}
        style={{
        content: {
            inset: '30%',
        }
        }}
    
    >
        <p>{me.winner && 'You win!'}</p>
        <p>{me.looser && 'You loose!'}</p>
        <button onClick={onClose}>Restart</button>
    </ReactModal>;
}