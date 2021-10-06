package cards

import (
	"math/rand"
	"time"
)

type Deck struct {
	Cards []Card
}

func NewDeck(amount int) *Deck {
	//@TODO check amount%4
	var cards = make([]Card, amount)
	var i = 0
	var minRank = Rank(int(MAX_RANK) - (amount / 4) + 1)

	for _, suite := range SuiteList {
		for rank := minRank; rank <= MAX_RANK; rank++ {
			cards[i] = *NewCard(i, suite, rank)
			i++
		}
	}

	return &Deck{cards}
}

func (d *Deck) Shuffle() *Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})

	return d
}

func (d *Deck) Shift(amount int) []Card {
	//TODO check if deck empty
	var cards = d.Cards[:amount]
	d.Cards = d.Cards[amount:len(d.Cards)]

	var newCards = make([]Card, len(cards))
	copy(newCards, cards)
	return newCards
}

func (d *Deck) Last() *Card {
	//TODO check if deck empty
	return &d.Cards[len(d.Cards)-1]
}
