import {useEffect} from "react";
import {useDrop} from "react-dnd";
import User from "../../service/User";
import {connect, send} from "../../websockets";

export function useDurakCardDrop(place) {
  return useDrop(() => ({
    accept: 'card',
    drop(item, monitor) {
      const {card} = item
      if (monitor.didDrop()) {
        return;
      }
      onCardDrop(card, place);
    },
    collect: (monitor) => ({
      isOver: monitor.isOver(),
    }),
  }), []);
}


export function useDurak(setState, setError) {
  useEffect(() => {
    connect(
      onMessage(setState, setError),
      onOpen,
      {playerId: User.getId()}
    )
  }, [])
}

const onOpen = () => {
  send({type: MESSAGE_CONNECTED, data: {playerId: User.getId()}})
}

export const onReady = (callback) => () => {
  send({type: MESSAGE_READY, data: {playerId: User.getId()}})
  callback()
}

export const onCardDrop = (card, place) => {
  send({type: MESSAGE_MOVE, data: {card, place: Number.isInteger(place) ? place : null}})
}


const MESSAGE_SELF_CONNECTED = 'self_connected';
const MESSAGE_CONNECTED = 'connected';
const MESSAGE_MOVE = 'move';
const MESSAGE_READY = 'ready';
const MESSAGE_STATE = 'state'
const MESSAGE_ERROR = 'error'

const onMessage = (setState, setError) => (message) => {
  console.log('MESSSAGE RECEIVED');
  console.log(message);

  switch(message.type) {
    
    case MESSAGE_SELF_CONNECTED:
      User.setId(message.data.playerId)
    break;

    case MESSAGE_STATE:
      setState(message.data)
    break;

    case MESSAGE_ERROR:
      setError(message.data)
    break;

    default:
      return
  }
}