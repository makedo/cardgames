package controllers

import (
	"cardgames/app/websocket"

	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	gorillaWebsocket "github.com/gorilla/websocket"
)

var upgrader = gorillaWebsocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func read(client *websocket.Client) {
	client.Pool.Register <- client

	defer func() {
		client.Pool.Unregister <- client
		client.Conn.Close()
	}()

	for {
		var message, err = client.Read()
		if err != nil {
			log.Println(err)
			return
		}

		client.Pool.Broadcast <- *message
		fmt.Printf("Message Received: %+v\n", message)
	}
}

func Websocket() func(http.ResponseWriter, *http.Request) {

	pool := websocket.NewPool()
	go pool.Listen()

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("WebSocket Endpoint Hit")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(w, "%+v\n", err)
		}

		client := &websocket.Client{
			Id:   uuid.New().String(),
			Conn: conn,
			Pool: pool,
		}

		go read(client)
	}
}
