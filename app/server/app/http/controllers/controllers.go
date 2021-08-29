package controllers

import (
	"api/domain/cards"
	"api/domain/games/durak"
	
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func GetDeck(w http.ResponseWriter, r *http.Request) {
	deck := cards.NewDeck(36)
	deck.Shuffle()

	out, err := json.Marshal(deck.Cards)
	if nil != err {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func GetState(w http.ResponseWriter, r *http.Request) {
	palyers := []string{uuid.New().String(), uuid.New().String()}

	state := durak.NewState(52, palyers)
	serState, err := state.ToSerializable(palyers[0])
	if nil != err {
		panic(err)
	}

	out, err := json.Marshal(serState)
	if nil != err {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
