import React, { Component } from "react";
import { connect, sendMessage } from "./../../websockets";
import ChatHistory from "./ChatHistory";
import ChatInput from "./ChatInput";

class Chat extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: []
    }
  }

  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log(this.state);
    });
  }

  send(event) {
    if(event.keyCode === 13) {
      sendMessage(event.target.value);
      event.target.value = "";
    }
  }

  render() {
    return (
      <div className="chat">
        <ChatHistory chatHistory={this.state.chatHistory} />
        <ChatInput send={this.send} />
      </div>
    );
  }
}

export default Chat;