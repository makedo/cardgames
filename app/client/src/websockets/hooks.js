import { useEffect } from "react";
import {connect, sendMessage} from '.'

export function useWebsockets(onMessageReceive) {
    useEffect(() => {
        connect(onMessageReceive);
    }, [onMessageReceive]);
      
    return sendMessage;
}