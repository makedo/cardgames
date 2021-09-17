package durak

import "cardgames/domain/cards"

const (
	PLAYER_STATE_ATTAKER  string = "attaker"
	PLAYER_STATE_DEFENDER string = "defender"
	PLAYER_STATE_IDLE     string = "idle"
)

type Player struct {
	Id    string
	Hand  *cards.Hand
	Ready bool
	State string
}

func NewPlayer(id string) *Player {
	return &Player{
		Id:    id,
		Ready: false,
		State: PLAYER_STATE_IDLE,
	}
}

func (p *Player) IsAttaker() bool {
	return p.State == PLAYER_STATE_ATTAKER
}

func (p *Player) IsDefender() bool {
	return p.State == PLAYER_STATE_DEFENDER
}
