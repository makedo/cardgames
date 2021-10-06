package websocket

type Handler interface {
	Handle(client *Client, message *Message)
	Disconnect(client *Client)
	Connect(client *Client)
}
