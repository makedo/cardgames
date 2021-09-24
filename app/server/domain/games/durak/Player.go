package durak

import "cardgames/domain/cards"

const (
	PLAYER_STATE_ATTAKER  string = "attaker"
	PLAYER_STATE_DEFENDER string = "defender"
	PLAYER_STATE_IDLE     string = "idle"
)

type Player struct {
	Id     string     `json:"id"`
	Hand   cards.Hand `json:"hand"`
	State  string     `json:"state"`
	Ready  bool       `json:"ready"`
	Winner bool       `json:"winner"`
	Looser bool       `json:"looser"`
}

type OtherPlayer struct {
	Id    string `json:"id"`
	Hand  int    `json:"hand"`
	State string `json:"state"`
}

func NewPlayer(id string) *Player {
	return &Player{
		Id:     id,
		Ready:  false,
		State:  PLAYER_STATE_IDLE,
		Winner: false,
		Looser: false,
	}
}

func (p *Player) IsAttaker() bool {
	return p.State == PLAYER_STATE_ATTAKER
}

func (p *Player) IsDefender() bool {
	return p.State == PLAYER_STATE_DEFENDER
}

func (p *Player) ToOtherPlayer() *OtherPlayer {
	return &OtherPlayer{
		Id:    p.Id,
		State: p.State,
		Hand:  len(p.Hand.Cards),
	}
}
