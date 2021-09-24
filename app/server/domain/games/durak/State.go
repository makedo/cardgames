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
	Table      *Table         `json:"table"`
	Me         *Player        `json:"me"`
	TrumpCard  *cards.Card    `json:"trump_card"`
	Players    []*OtherPlayer `json:"players"`
	CanConfirm bool           `json:"can_confirm"`
	DeckAmount int            `json:"deck_amount"`
	Started    bool           `json:"started"`
}

func NewState(deckAmount int, players []*Player) *State {
	var deck = cards.NewDeck(deckAmount).Shuffle()
	var trump = deck.Last()

	for _, player := range players {
		player.Hand = cards.Hand{Cards: []cards.Card{}}
		player.Winner = false
		player.Looser = false
		player.Ready = false
		player.State = PLAYER_STATE_IDLE
	}

	return &State{
		table:   *NewTable(),
		deck:    *deck,
		trump:   *trump,
		started: false,
		players: players,
	}
}

func (s *State) Start() {

	//@TODO make it work for more than 2 players
	for i, player := range s.players {
		if 0 == i {
			player.State = PLAYER_STATE_ATTAKER
		} else {
			player.State = PLAYER_STATE_DEFENDER
		}
		var handCards = s.deck.Shift(6)
		player.Hand = cards.Hand{Cards: handCards}
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
		var defender *Player
		for _, otherPlayer := range s.players {
			if otherPlayer.IsDefender() {
				defender = otherPlayer
			}
		}

		//@TODO Check first attak with 5 cards
		//@TODO Check no more than 6 cards on table
		if s.table.CountUndefendedCards()+1 > len(defender.Hand.Cards) {
			return errors.New("Defender has too many cards to defend")
		}

		if false == s.table.HasCards() || s.table.HasCardOfSameRank(card) {
			s.table.AddCard(card, nil)
			player.Hand.PopCardById(cardId)
			return nil
		}

		return errors.New("Attaker can not make move with card")
	}

	if player.IsDefender() {
		if len(s.table.Cards) == 0 {
			return errors.New("Table is empty")
		}
		var placeCard, validPlace = s.table.GetCardFromPlace(place)

		if (placeCard.Suite == card.Suite && card.Rank > placeCard.Rank) ||
			(card.Suite == s.trump.Suite && placeCard.Suite != s.trump.Suite) {
			s.table.AddCard(card, &validPlace)
			player.Hand.PopCardById(cardId)
			return nil
		}

		return errors.New("Defender can not make move with card")
	}

	return nil
}

func (s *State) Confirm(player *Player) error {
	if false == s.CanConfirm(player) {
		return errors.New("Player can not confirm")
	}

	if player.IsDefender() {
		for _, cardsPlace := range s.table.Cards {
			for _, card := range cardsPlace {
				player.Hand.AddCard(card)
			}
		}
	}
	s.table = *NewTable()

	for _, otherPlayer := range s.players {
		for len(otherPlayer.Hand.Cards) < 6 && len(s.deck.Cards) > 0 {
			var cardFromDeck = s.deck.Shift(1)
			otherPlayer.Hand.AddCard(cardFromDeck[0])
		}
	}

	if false == player.IsDefender() {
		//@TODO mke work for more than 2 players
		for _, otherPlayer := range s.players {
			if otherPlayer.IsAttaker() {
				otherPlayer.State = PLAYER_STATE_DEFENDER
			} else {
				if otherPlayer.IsDefender() {
					otherPlayer.State = PLAYER_STATE_ATTAKER
				}
			}
		}
	}

	return nil
}

func (s *State) CanConfirm(player *Player) bool {

	if player.IsAttaker() {
		return s.table.AllCardsDefended()
	}

	if player.IsDefender() {
		return s.table.HasUndefendedCards()
	}

	return false
}

func (s *State) ToSerializable(currentPlayerId string) (*SerializableState, error) {
	var players = []*OtherPlayer{}
	var me *Player

	for _, player := range s.players {
		if player.Id != currentPlayerId {
			players = append(players, player.ToOtherPlayer())
		} else {
			me = player
		}
	}

	if me == nil {
		return nil, errors.New("Can't find my hand")
	}

	return &SerializableState{
		Table:      &s.table,
		Me:         me,
		TrumpCard:  &s.trump,
		Players:    players,
		CanConfirm: s.CanConfirm(me),
		DeckAmount: len(s.deck.Cards),
		Started:    s.started,
	}, nil
}
