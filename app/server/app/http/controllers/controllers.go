package controllers

import (
	"cardgames/domain/cards"
	"encoding/json"
	"net/http"
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
