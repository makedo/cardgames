package websocket

import (
	"fmt"
)

type Pool struct {
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan Message
	BroadcastTo chan MessagePool

	Clients map[*Client]bool
}

type MessagePool struct {
	Message Message
	Clients []*Client
}

func NewMessagePool(message Message, clients ...*Client) *MessagePool {
	return &MessagePool{
		Message: message,
		Clients: clients,
	}
}

func NewPool() *Pool {
	return &Pool{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan Message),
		BroadcastTo: make(chan MessagePool),

		Clients: make(map[*Client]bool),
	}
}

func (pool *Pool) Listen() {
	for {
		select {
		case client := <-pool.Register:
			var message = &Message{
				Type: MESSAGE_TYPE_CONNECTED,
				Data: map[string]interface{}{"playerId": client.Id},
			}
			for client, _ := range pool.Clients {
				client.Write(message)
			}
			pool.Clients[client] = true
			fmt.Println(pool.Clients)
			break

		case client := <-pool.Unregister:
			client.Conn.Close()
			delete(pool.Clients, client)

			var message = &Message{
				Type: MESSAGE_TYPE_DISCONNECTED,
				Data: map[string]interface{}{"playerId": client.Id},
			}
			for client, _ := range pool.Clients {
				client.Write(message)
			}
			break

		case message := <-pool.Broadcast:
			for client, _ := range pool.Clients {
				if err := client.Write(&message); err != nil {
					fmt.Println(err)
					return
				}
			}
			break

		case messagePool := <-pool.BroadcastTo:
			for _, client := range messagePool.Clients {
				if true == pool.Clients[client] {
					if err := client.Write(&messagePool.Message); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
			break
		}
	}
}
