package durak

import (
	"cardgames/domain/cards"
	"errors"
)

type State struct {
	table   Table
	deck    cards.Deck
	trump   cards.Card
	started bool
	players []*Player
}

type SerializableState struct {
	Table      *Table       `json:"table"`
	Hand       []cards.Card `json:"hand"`
	TrumpCard  *cards.Card  `json:"trump_card"`
	TrumpSuite *cards.Suite `json:"trump_suite"`
	Hands      []uint       `json:"hands"`
}

func NewState(deckAmount int, players []*Player) *State {
	var deck = cards.NewDeck(deckAmount).Shuffle()
	var trump = deck.Last()

	return &State{
		table:   *NewTable(),
		deck:    *deck,
		trump:   *trump,
		started: false,
		players: players,
	}
}

func (s *State) Start() {

	for i, player := range s.players {
		if 0 == i {
			player.State = PLAYER_STATE_ATTAKER
		} else {
			player.State = PLAYER_STATE_DEFENDER
		}
		player.Hand = &cards.Hand{Cards: s.deck.Pop(6)}
	}

	s.started = true
}

func (s *State) isStarted() bool {
	return s.started
}

func (s *State) SetPlayerReady(playerId string) {
	var player, exist = s.GetPlayer(playerId)
	if exist {
		player.Ready = true
	}
}

func (s *State) AreAllPlayersReady() bool {
	for _, player := range s.players {
		if false == player.Ready {
			return false
		}
	}

	return true
}

func (s *State) AddPlayer(playerId string) {
	if _, exists := s.GetPlayer(playerId); exists {
		return
	}

	var player = &Player{
		Id: playerId,
	}
	s.players = append(s.players, player)
}

func (s *State) RemovePlayer(playerId string) {
	var newPlayers []*Player

	for _, player := range s.players {
		if player.Id != playerId {
			newPlayers = append(newPlayers, player)
		}
	}

	s.players = newPlayers
}

func (s *State) GetPlayer(playerId string) (*Player, bool) {
	for _, player := range s.players {
		if player.Id == playerId {
			return player, true
		}
	}

	return nil, false
}

func (s *State) GetAmountOfPlayers() int {
	return len(s.players)
}

func (s *State) GetPlayers() []*Player {
	return s.players
}

func (s *State) Move(playerId string, cardId int, place *int) error {
	var player, hasPlayer = s.GetPlayer(playerId)
	if !hasPlayer {
		return errors.New("Player not found")
	}

	card, hasCard := player.Hand.GetCardById(cardId)
	if !hasCard {
		return errors.New("Card id not found")
	}

	if player.IsAttaker() {
		if false == s.table.HasCards() || s.table.HasCardOfSameRank(card) {
			s.table.AddCard(card, nil)
			player.Hand.PopCardById(cardId)
			return nil
		}

		return errors.New("Attaker can not make move with card")
	}

	if player.IsDefender() {
		if nil == place {
			return errors.New("Place is nil")
		}

		var placeCard = s.table.GetCardFromPlace(*place)

		if ((placeCard.Suite == card.Suite && card.Rank > placeCard.Rank) ||
			(card.Suite == s.trump.Suite && placeCard.Suite != s.trump.Suite)) {
			s.table.AddCard(card, place)
			player.Hand.PopCardById(cardId)
			return nil
		}

		return errors.New("Defender can not make move with card")
	}

	return nil
}

func (s *State) ToSerializable(currentPlayerId string) (*SerializableState, error) {
	var myHand *cards.Hand
	var hands = make([]uint, len(s.players)-1)

	var i = 0

	for _, player := range s.players {
		if player.Id != currentPlayerId {
			hands[i] = uint(len(player.Hand.Cards))
			i++
		} else {
			myHand = player.Hand
		}
	}

	if myHand == nil {
		return nil, errors.New("Can't find my hand")
	}

	return &SerializableState{
		Table:      &s.table,
		Hand:       myHand.Cards,
		TrumpCard:  &s.trump,
		TrumpSuite: &s.trump.Suite,
		Hands:      hands,
	}, nil
}
