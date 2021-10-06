package durak

import (
	"cardgames/app/websocket"
	"fmt"
	"log"
	"time"
)

const TIME_TO_WAIT_FOR_RECONNECT_SEC = 5
const CARDS_IN_DECK = 36
const MAX_PLAYERS = 4

type Handler struct {
	state      *State
	maxPlayers int
	timers     map[string]*time.Timer
}

func NewHandler() *Handler {
	return &Handler{
		state:      NewState(CARDS_IN_DECK, NewPlayersCollection([]*Player{})),
		maxPlayers: MAX_PLAYERS, //@TODO should be in state and set from client
		timers:     make(map[string]*time.Timer),
	}
}

func (h *Handler) Handle(client *websocket.Client, message *websocket.Message) {

	var playerId = client.Id
	var player, hasPlayer = h.state.GetPlayer(playerId)
	if false == hasPlayer && message.Type != websocket.MESSAGE_TYPE_CONNECTED {
		log.Fatalln("Player not found!")
		return
	}

	switch message.Type {
	case websocket.MESSAGE_TYPE_CONNECTED:
		var timer, hasTimer = h.timers[playerId]

		if hasTimer {
			timer.Stop()
			delete(h.timers, playerId)
		}

		if h.state.isStarted() && hasPlayer {
			h.broadcastState(client)
			return
		}

		if h.state.isStarted() && false == hasPlayer {
			var message = websocket.Message{
				Type: websocket.MESSAGE_TYPE_ERROR,
				Data: map[string]interface{}{"message": "Game has been alreay started"},
			}
			client.Pool.BroadcastTo <- *websocket.NewMessagePool(message, client)
			return
		}

		if h.state.GetAmountOfPlayers() >= h.maxPlayers {
			var message = websocket.Message{
				Type: websocket.MESSAGE_TYPE_ERROR,
				Data: map[string]interface{}{"message": "Too many players"},
			}
			client.Pool.BroadcastTo <- *websocket.NewMessagePool(message, client)
			return
		}

		h.state.AddPlayer(client.Id)
		var message = websocket.Message{
			Type: websocket.MESSAGE_TYPE_SELF_CONNECTED,
			Data: map[string]interface{}{"playerId": client.Id},
		}
		client.Pool.BroadcastTo <- *websocket.NewMessagePool(message, client)
		break

	case websocket.MESSAGE_TYPE_READY:
		if true == h.state.isStarted() {
			log.Println("Game already started")
			return
		}

		if false == player.Ready {
			h.state.SetPlayerReady(client.Id)
		}

		if h.state.AreAllPlayersReady() && h.state.GetAmountOfPlayers() == h.maxPlayers {
			h.state.Start()
		}

		h.broadcastStateForAll(client.Pool)

		return

	case websocket.MESSAGE_TYPE_RESTART:

		if true == h.state.isStarted() {
			log.Println("Game did not finished")
			return
		}

		if false == player.Ready {
			h.state.SetPlayerReady(client.Id)
		}

		if h.state.AreAllPlayersReady() && h.state.GetAmountOfPlayers() == h.maxPlayers {
			h.state.Start()
			h.broadcastStateForAll(client.Pool)
		} else {
			h.broadcastState(client)
		}

		return

	case websocket.MESSAGE_TYPE_MOVE:
		var data = message.GetData()

		//@TODO separate message type and cast map to type
		var cardData = data["card"].(map[string]interface{})
		cardId, ok := cardData["id"]
		if !ok {
			log.Println("Wrong structure of move message")
			return
		}
		cardIdInt, ok := cardId.(float64)
		if !ok {
			log.Printf("Got data of type %T but wanted int for card.id", cardIdInt)
			return
		}

		var place *int
		var placeData = data["place"]
		if nil != placeData {
			placeDataFloat, ok := placeData.(float64)
			if ok {
				var placeInt = int(placeDataFloat)
				place = &placeInt
			}
		}

		var error = h.state.Move(playerId, int(cardIdInt), place)
		if nil != error {
			log.Println(error)
			return
		}

		//@TODO calculate for more than 2 players
		var IsGameFinished, errorFinished = h.state.FinishGame(playerId)
		if nil != errorFinished {
			log.Println(errorFinished)
			return
		}

		h.broadcastStateForAll(client.Pool)
		if IsGameFinished {
			h.state = NewState(CARDS_IN_DECK, &h.state.players)
			fmt.Println("NEW STATE")
		}

		return

	case websocket.MESSAGE_TYPE_CONFIRM:

		if false == h.state.CanConfirm(player) {
			log.Println("Player can not confirm")
			return
		}

		err := h.state.Confirm(player)
		if nil != err {
			log.Println(err)
			return
		}

		h.broadcastStateForAll(client.Pool)

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
			h.state = NewState(CARDS_IN_DECK, &h.state.players)
			//@TODO send RESTART message to other players
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

	var outBoundMessage = websocket.Message{
		Type: websocket.MESSAGE_TYPE_STATE,
		Data: serializableState,
	}

	client.Pool.BroadcastTo <- *websocket.NewMessagePool(outBoundMessage, client)

	return nil
}

func (h *Handler) broadcastStateForAll(pool *websocket.Pool) {
	for client := range pool.Clients {
		if err := h.broadcastState(client); nil != err {
			log.Fatal(err)
		}
	}
}
