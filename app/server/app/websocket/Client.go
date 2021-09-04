package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Write(message *Message) error {
	var err = c.Conn.WriteJSON(message)
	return err
}

func (c *Client) Read() (*Message, error) {
	var message = &Message{}
	err := c.Conn.ReadJSON(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}