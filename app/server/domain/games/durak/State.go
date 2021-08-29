package durak

import (
	"cardgames/domain/cards"
	"errors"
)

type Hand struct {
	Id    string
	Cards []cards.Card
}

type Table struct {
	Cards [][2]cards.Card `json:"cards"`
}

type State struct {
	table Table
	deck  cards.Deck
	hands []Hand
	trump cards.Card
}

type SerializableState struct {
	Table      *Table       `json:"table"`
	Hand       []cards.Card `json:"hand"`
	TrumpCard  *cards.Card  `json:"trump_card"`
	TrumpSuite *cards.Suite `json:"trump_suite"`
	Hands      []uint       `json:"hands"`
}

func NewState(deckAmount int, playerIds []string) *State {
	var deck = cards.NewDeck(deckAmount).Shuffle()
	var trump = deck.Last()
	var hands = make([]Hand, len(playerIds))

	var i = 0
	for _, playerId := range playerIds {
		hands[i] = Hand{Id: playerId, Cards: deck.Pop(6)}
		i++
	}

	return &State{
		table: Table{},
		deck:  *deck,
		hands: hands,
		trump: *trump,
	}
}

func (s *State) ToSerializable(currentPlayerId string) (*SerializableState, error) {
	var myHand *Hand = nil
	var hands = make([]uint, len(s.hands)-1)

	var i = 0
	for _, hand := range s.hands {
		if hand.Id != currentPlayerId {
			hands[i] = uint(len(hand.Cards))
			i++
		} else {
			myHand = &hand
		}
	}

	if nil == myHand {
		return nil, errors.New("can't find my hand")
	}

	return &SerializableState{
		Table:      &s.table,
		Hand:       myHand.Cards,
		TrumpCard:  &s.trump,
		TrumpSuite: &s.trump.Suite,
		Hands:      hands,
	}, nil
}
