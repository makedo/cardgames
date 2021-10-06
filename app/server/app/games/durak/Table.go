package durak

import (
	"cardgames/app/games/cards"
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

func (t *Table) HasCardOfSameRank(card *cards.Card) bool {
	for _, cardPlaces := range t.Cards {
		for _, tableCard := range cardPlaces {
			if tableCard.Rank == card.Rank {
				return true
			}
		}
	}

	return false
}

func (t *Table) GetCardFromPlace(place *int) (*cards.Card, int) {
	//@TODO check if card and place exists
	if place == nil {
		return &t.Cards[len(t.Cards)-1][0], len(t.Cards) - 1
	} else {
		return &t.Cards[*place][0], *place
	}

}

func (t *Table) AllCardsDefended() bool {

	if false == t.HasCards() {
		return false
	}

	for _, cardsPlace := range t.Cards {
		if len(cardsPlace) < 2 {
			return false
		}
	}

	return true
}

func (t *Table) HasUndefendedCards() bool {
	if false == t.HasCards() {
		return false
	}

	for _, cardsPlace := range t.Cards {
		if len(cardsPlace) == 1 {
			return true
		}
	}

	return false
}

func (t *Table) CountUndefendedCards() int {
	var count = 0

	for _, cardsPlace := range t.Cards {
		if len(cardsPlace) == 1 {
			count++
		}
	}

	return count
}

func (t *Table) GetPlainCars() []cards.Card {
	var cards = []cards.Card{}
	for _, cardsPlace := range t.Cards {
		for _, card := range cardsPlace {
			cards = append(cards, card)
		}
	}

	return cards
}
