package websocket

import (
	"fmt"
)

type Pool struct {
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan Message
	BroadcastTo chan ClientsPool

	Clients map[*Client]bool
}

type ClientsPool struct {
	Message Message
	Clients []*Client
}

func NewClientsPool(message Message, clients ...*Client) *ClientsPool {
	return &ClientsPool{
		Message: message,
		Clients: clients,
	}
}

func NewPool() *Pool {
	return &Pool{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan Message),
		BroadcastTo: make(chan ClientsPool),

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
			for client := range pool.Clients {
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
			for client := range pool.Clients {
				client.Write(message)
			}
			break

		case message := <-pool.Broadcast:
			for client := range pool.Clients {
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
