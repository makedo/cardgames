package durak

import (
	"cardgames/domain/cards"
)

const (
	PLAYER_STATE_ATTAKER     string = "attaker"
	PLAYER_STATE_SUB_ATTAKER string = "sub_attaker"
	PLAYER_STATE_DEFENDER    string = "defender"
	PLAYER_STATE_IDLE        string = "idle"
)

type Player struct {
	Id        string     `json:"id"`
	Hand      cards.Hand `json:"hand"`
	State     string     `json:"state"`
	Ready     bool       `json:"ready"`
	Winner    bool       `json:"winner"`
	Looser    bool       `json:"looser"`
	Confirmed bool       `json:"confirmed"`
}

type OtherPlayer struct {
	Id        string `json:"id"`
	Hand      int    `json:"hand"`
	State     string `json:"state"`
	Confirmed bool   `json:"confirmed"`
	Winner    bool   `json:"winner"`
	Looser    bool   `json:"looser"`
}

func NewPlayer(id string) *Player {
	return &Player{
		Id:        id,
		Ready:     false,
		State:     PLAYER_STATE_IDLE,
		Winner:    false,
		Looser:    false,
		Confirmed: false,
	}
}

func (p *Player) IsAttaker() bool {
	return p.State == PLAYER_STATE_ATTAKER || p.State == PLAYER_STATE_SUB_ATTAKER
}

func (p *Player) IsSubAttaker() bool {
	return p.State == PLAYER_STATE_SUB_ATTAKER
}

func (p *Player) IsDefender() bool {
	return p.State == PLAYER_STATE_DEFENDER
}

func (p *Player) IsIdle() bool {
	return p.State == PLAYER_STATE_IDLE
}

func (p *Player) ToOtherPlayer() *OtherPlayer {
	return &OtherPlayer{
		Id:        p.Id,
		State:     p.State,
		Hand:      len(p.Hand.Cards),
		Confirmed: p.Confirmed,
		Winner:    p.Winner,
		Looser:    p.Looser,
	}
}

type PlayersCollection struct {
	players []*Player
}

func NewPlayersCollection(players []*Player) *PlayersCollection {
	return &PlayersCollection{players}
}

func (pc *PlayersCollection) AddPlayer(player *Player) bool {
	if _, exists := pc.GetPlayerById(player.Id); exists {
		return false
	}

	pc.players = append(pc.players, player)
	return true
}

func (pc *PlayersCollection) GetPlayerById(playerId string) (*Player, bool) {
	for _, player := range pc.players {
		if player.Id == playerId {
			return player, true
		}
	}

	return nil, false
}

func (pc *PlayersCollection) RemovePlayerById(playerId string) {
	var newPlayers []*Player

	for _, player := range pc.players {
		if player.Id != playerId {
			newPlayers = append(newPlayers, player)
		}
	}

	pc.players = newPlayers
}

func (pc *PlayersCollection) GetAmount() int {
	return len(pc.players)
}

func (pc *PlayersCollection) Each(callback func(player *Player)) {
	for _, player := range pc.players {
		callback(player)
	}
}

func (pc *PlayersCollection) Are(callback func(player *Player) bool) bool {
	for _, player := range pc.players {
		if false == callback(player) {
			return false
		}
	}

	return true
}

func (pc *PlayersCollection) Find(callback func(player *Player) bool) (*Player, bool) {
	for _, player := range pc.players {
		if true == callback(player) {
			return player, true
		}
	}

	return nil, false
}

func (pc *PlayersCollection) FindIndex(callback func(player *Player) bool) (*int, bool) {
	for index, player := range pc.players {
		if true == callback(player) {
			return &index, true
		}
	}

	return nil, false
}

func (pc *PlayersCollection) EachStartingFromIndex(index int, callback func(step int, player *Player)) {
	for i := 0; i < pc.GetAmount(); i++ {
		if index >= pc.GetAmount() {
			index = 0
		}
		callback(i, pc.players[index])
		index++
	}
}

func (pc *PlayersCollection) Count(callback func(player *Player) bool) int {
	var count = 0

	for _, player := range pc.players {
		if true == callback(player) {
			count++
		}
	}

	return count
}
