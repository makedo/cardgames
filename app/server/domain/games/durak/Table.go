package durak

import (
	"cardgames/domain/cards"
)


type Table struct {
	Cards [][]cards.Card `json:"cards"`
}

func NewTable() *Table {
	return &Table{Cards: [][]cards.Card{}}
}

func (t *Table) HasCards() bool {
	return len(t.Cards) > 0
}


func (t *Table) AddCard(card *cards.Card, index *int) {
	if nil == index {
		var cards = []cards.Card{*card}
		t.Cards = append(t.Cards, cards)
	} else {
		//TODO check if card[index] was allocated
		t.Cards[*index] = append(t.Cards[*index], *card)
	}
}

func(t *Table) HasCardOfSameRank(card *cards.Card) bool {
	for _, cardPlaces := range t.Cards {
		for _, tableCard := range cardPlaces {
			if (tableCard.Rank == card.Rank) {
				return true
			}
		}
	}

	return false
}

func (t *Table) GetCardFromPlace(place int ) (*cards.Card) {
	
	var card = t.Cards[place][0]
	
	//@TODO check if card and place exists

	return &card
}