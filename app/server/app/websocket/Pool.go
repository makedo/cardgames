package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message

	Clients map[*Client]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Listen() {
	for {
		select {
		case client := <-pool.Register:

			client.Write(&Message{
				Type: SELF_CONNECTED,
				Data: struct {
					Id string `json:"id"`
				}{
					Id: client.Id,
				},
			})

			for client, _ := range pool.Clients {
				client.Write(&Message{Type: CONNECTED, Data: "New User Joined..."})
			}
			pool.Clients[client] = true
			break

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			var message = &Message{Type: DISCONNECTED, Data: "User has disconnected Joined..."}
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
		}
	}
}
