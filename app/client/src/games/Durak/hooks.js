import {useState, useEffect} from "react";
import {useDrop} from "react-dnd";

import api from "./api";
import {useWebsockets} from "../../websockets/hooks";

export function useDurakCardDrop(onDrop) {
  return useDrop(() => ({
    accept: 'card',
    drop({card}, monitor) {
      const didDrop = monitor.didDrop();
      if (didDrop) {
        return;
      }
      onDrop(card);
    },
    collect: (monitor) => ({
      isOver: monitor.isOver(),
    }),
  }), []);
}

export function useDurakState() {
  const [state, setState] = useState(null);

  useEffect(() => {
    const fetchState = async () => {
      const state = await api('/durak');
      setState(state);
    }

    fetchState();
  }, []);
  
  return [state, setState];
}

const MESSAGE_SELF_CONNECTED = 'self_connected';
const MESSAGE_READY = 'reay';

const onMessageReceive = function(message) {
  console.log('MESSSAGE RECEIVED');
  console.log(message);

  switch(message.type) {
    case MESSAGE_SELF_CONNECTED:
      localStorage.setItem("id", message.data.id);
    break;
    case MESSAGE_READY:
      
    break;
    default:
      console.log(message)
  }
}

export function useDurakWebsockets() {
  return useWebsockets(onMessageReceive);
}