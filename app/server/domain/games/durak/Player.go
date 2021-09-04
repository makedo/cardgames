package durak

import "cardgames/domain/cards"

type Player struct {
	Id   string
	Hand *cards.Hand
	Ready bool
}

func NewPlayer(id string) *Player {
	return &Player{
		Id: id,
		Ready: false,
	}
}

type PlayerCollection struct {
	players []*Player
}

