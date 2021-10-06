package controllers

import (
	"cardgames/app/websocket"
	"cardgames/app/games/durak"

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

var handler = durak.NewHandler() //@todo choose handler according to a game

func Websocket() func(http.ResponseWriter, *http.Request) {

	pool := websocket.NewPool()
	go pool.Listen()

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("WebSocket Endpoint Hit")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(w, "%+v\n", err)
		}

		var client = createClient(r, conn, pool)

		read(client)
	}
}

func read(client *websocket.Client) {

	handler.Connect(client)
	defer func() {
		handler.Disconnect(client)
	}()

	for {
		var message, err = client.Read()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Message Received: %+v\n", message)
		handler.Handle(client, message)
	}
}

func createClient(r *http.Request, conn *gorillaWebsocket.Conn, pool *websocket.Pool) *websocket.Client {
	var query = r.URL.Query()["playerId"]
	var playerId string

	if len(query) > 0 && len(query[0]) > 0 {
		playerId = query[0]
	} else {
		playerId = uuid.New().String()
	}

	return &websocket.Client{
		Id:   playerId,
		Conn: conn,
		Pool: pool,
	}
}
