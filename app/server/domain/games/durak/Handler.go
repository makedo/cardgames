package durak

import (
	"cardgames/app/websocket"
	"fmt"
	"log"
	"time"
)

const TIME_TO_WAIT_FOR_RECONNECT_SEC = 5
const CARDS_IN_DECK = 36

type Handler struct {
	state      *State
	maxPlayers int
	timers     map[string]*time.Timer
}

func NewHandler() *Handler {
	return &Handler{
		state:      NewState(CARDS_IN_DECK, []*Player{}),
		maxPlayers: 2,
		timers:     make(map[string]*time.Timer),
	}
}

func (h *Handler) Handle(client *websocket.Client, message *websocket.Message) {
	switch message.Type {
	case websocket.MESSAGE_TYPE_CONNECTED:
		fmt.Println("CONNECTED")
		var playerId = client.Id
		var _, hasPlayer = h.state.GetPlayer(playerId)

		var timer, hasTimer = h.timers[playerId]
		fmt.Println("TIMERS")
		fmt.Println(h.timers)

		if hasTimer {
			fmt.Println("TIMER")
			fmt.Println(h.timers)
			timer.Stop()
			delete(h.timers, playerId)
		}

		if h.state.isStarted() && hasPlayer {
			h.broadcastState(client)
			return
		}

		if h.state.isStarted() && false == hasPlayer {
			var message = &websocket.Message{
				Type: websocket.MESSAGE_TYPE_ERROR,
				Data: map[string]interface{}{"message": "Game has been alreay started"},
			}
			client.Pool.BroadcastTo <- websocket.NewMessagePool(message, client)
			return
		}

		if h.state.GetAmountOfPlayers() > h.maxPlayers {
			var message = &websocket.Message{
				Type: websocket.MESSAGE_TYPE_ERROR,
				Data: map[string]interface{}{"message": "Too many players"},
			}
			client.Pool.BroadcastTo <- websocket.NewMessagePool(message, client)
			return
		}

		h.state.AddPlayer(client.Id)
		var message = &websocket.Message{
			Type: websocket.MESSAGE_TYPE_SELF_CONNECTED,
			Data: map[string]interface{}{"playerId": client.Id},
		}
		client.Pool.BroadcastTo <- websocket.NewMessagePool(message, client)
		break

	case websocket.MESSAGE_TYPE_READY:
		var playerId = client.Id
		var _, hasPlayer = h.state.GetPlayer(playerId)

		if false == hasPlayer {
			return
		}

		if true == h.state.isStarted() {
			return
		}

		h.state.SetPlayerReady(client.Id)
		if false == h.state.AreAllPlayersReady() || h.state.GetAmountOfPlayers() < h.maxPlayers {
			return
		}

		h.state.Start()
		for client, _ := range client.Pool.Clients {
			if err := h.broadcastState(client); nil != err {
				log.Fatal(err)
			}
		}

		return
	}

	// client.Pool.Broadcast <- *message
}

func (h *Handler) Disconnect(client *websocket.Client) {
	var playerId = client.Id

	client.Pool.Unregister <- client

	if false == h.state.isStarted() {
		h.state.RemovePlayer(playerId)
	} else {
		var timer = time.NewTimer(TIME_TO_WAIT_FOR_RECONNECT_SEC * time.Second)
		h.timers[playerId] = timer
		go func() {
			<-timer.C
			fmt.Println("TIMER FIRED")
			h.state.RemovePlayer(playerId)
			h.state = NewState(CARDS_IN_DECK, h.state.GetPlayers())
		}()
	}
}

func (h *Handler) Connect(client *websocket.Client) {
	client.Pool.Register <- client
}

func (h *Handler) broadcastState(client *websocket.Client) error {
	var _, hasPlayer = h.state.GetPlayer(client.Id)

	if false == hasPlayer {
		return nil
	}

	serializableState, err := h.state.ToSerializable(client.Id)
	if nil != err {
		log.Fatal(err)
		return err
	}

	var outBoundMessage = &websocket.Message{
		Type: websocket.MESSAGE_TYPE_STATE,
		Data: serializableState,
	}

	client.Pool.BroadcastTo <- websocket.NewMessagePool(outBoundMessage, client)

	return nil
}
