package durak

import (
	"cardgames/app/games/cards"
	"errors"
	"fmt"
)

type State struct {
	table    Table
	deck     cards.Deck
	trump    cards.Card
	started  bool
	players  PlayersCollection
	finished bool
}

type SerializableState struct {
	Table      *Table         `json:"table"`
	Me         *Player        `json:"me"`
	TrumpCard  *cards.Card    `json:"trump_card"`
	Players    []*OtherPlayer `json:"players"`
	CanConfirm bool           `json:"can_confirm"`
	DeckAmount int            `json:"deck_amount"`
	Started    bool           `json:"started"`
	Finished   bool           `json:"finished"`
}

func NewState(deckAmount int, players *PlayersCollection) *State {
	var deck = cards.NewDeck(deckAmount).Shuffle()
	var trump = deck.Last()

	players.Each(func(player *Player) {
		player.Hand = cards.Hand{Cards: []cards.Card{}}
		player.Winner = false
		player.Looser = false
		player.Ready = false
		player.Role = PLAYER_ROLE_IDLE
		player.Confirmed = false
	})

	return &State{
		table:   *NewTable(),
		deck:    *deck,
		trump:   *trump,
		started: false,
		players: *players,
	}
}

func (s *State) Start() {
	var i = 0
	var amount = s.players.GetAmount()

	s.players.Each(func(player *Player) {
		if 0 == i { //first
			player.Role = PLAYER_ROLE_ATTAKER
		} else if 1 == i { //second
			player.Role = PLAYER_ROLE_DEFENDER
		} else if (amount - 1) == i { //last
			player.Role = PLAYER_ROLE_SUB_ATTAKER
		} else {
			player.Role = PLAYER_ROLE_IDLE
		}

		player.Hand = cards.Hand{Cards: s.deck.Shift(6)}
		i++
	})

	s.started = true
}

func (s *State) isStarted() bool {
	return s.started
}

func (s *State) SetPlayerReady(playerId string) {
	var player, exist = s.players.GetPlayerById(playerId)
	if exist {
		player.Ready = true
	}
}

func (s *State) AreAllPlayersReady() bool {
	return s.players.Are(func(player *Player) bool {
		return player.Ready
	})
}

func (s *State) AddPlayer(playerId string) {
	var player = &Player{
		Id: playerId,
	}

	s.players.AddPlayer(player)
}

func (s *State) RemovePlayer(playerId string) {
	s.players.RemovePlayerById(playerId)
}

func (s *State) GetPlayer(playerId string) (*Player, bool) {
	return s.players.GetPlayerById(playerId)
}

func (s *State) GetAmountOfPlayers() int {
	return s.players.GetAmount()
}

func (s *State) Move(playerId string, cardId int, place *int) error {
	var player, hasPlayer = s.players.GetPlayerById(playerId)
	if !hasPlayer {
		return errors.New("Player not found")
	}

	if player.Confirmed {
		return errors.New("Player has confirmed, can not move")
	}

	card, hasCard := player.Hand.GetCardById(cardId)
	if !hasCard {
		return errors.New("Card not found")
	}

	if player.IsSubAttaker() {
		if false == s.table.HasCards() {
			return errors.New("Attaker should move first")
		}
	}

	//attak
	if player.IsAttaker() {

		var defender, found = s.players.Find(func(player *Player) bool {
			return player.IsDefender()
		})

		if false == found {
			return errors.New("Defender not found")
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

	//defend
	if player.IsDefender() {
		if len(s.table.Cards) == 0 {
			return errors.New("Table is empty")
		}
		var placeCard, placeIndex = s.table.GetCardFromPlace(place)

		if (placeCard.Suite == card.Suite && card.Rank > placeCard.Rank) ||
			(card.Suite == s.trump.Suite && placeCard.Suite != s.trump.Suite) {
			s.table.AddCard(card, &placeIndex)
			player.Hand.PopCardById(cardId)
			return nil
		}

		return errors.New("Defender can not make move with card")
	}

	return nil
}

func (s *State) Confirm(player *Player) error {

	var defenderTakesCards = s.table.HasUndefendedCards()

	player.Confirmed = true
	var areAllPlayersConfirmed = s.players.Are(func(player *Player) bool {
		if player.Winner {
			return true //autoconfirm for winners
		}
		if player.IsIdle() {
			return true //autoconfirm for idles
		}
		if false == defenderTakesCards && player.IsDefender() {
			return true //autoconfirm for defender, in case if defender defends all cards
		}

		return player.Confirmed
	})

	if false == areAllPlayersConfirmed {
		return nil
	}

	var attakerIndex, foundAttacker = s.players.FindIndex(func(player *Player) bool {
		return player.IsAttaker() && !player.IsSubAttaker()
	})
	if !foundAttacker {
		return errors.New("Attacker is not found")
	}

	s.players.EachStartingFromIndex(*attakerIndex, func(step int, player *Player) {
		if defenderTakesCards && player.IsDefender() {
			for _, cardsPlace := range s.table.Cards {
				for _, card := range cardsPlace {
					player.Hand.AddCard(card)
				}
			}
		} else {
			for len(player.Hand.Cards) < 6 && len(s.deck.Cards) > 0 {
				var cardFromDeck = s.deck.Shift(1)
				player.Hand.AddCard(cardFromDeck[0])
			}
		}
	})

	var iterateFromIndex, found = s.players.FindIndex(func(player *Player) bool {
		return player.IsDefender()
	})
	if !found {
		return errors.New("Defender is not found")
	}

	//if defenderTakesCards then next player after defender becomes attaker, next defender, next subattaker
	//else defender becomes attaker, next defender, next subattcker
	if defenderTakesCards {
		*iterateFromIndex++
	}

	var currentStep = 0
	var stateChangeCallback = func(step int, player *Player) {
		if player.Winner {
			player.Role = PLAYER_ROLE_IDLE
			return
		}

		player.Confirmed = false
		switch currentStep {
		case 0:
			player.Role = PLAYER_ROLE_ATTAKER
		case 1:
			player.Role = PLAYER_ROLE_DEFENDER
		case 2:
			player.Role = PLAYER_ROLE_SUB_ATTAKER
		default:
			player.Role = PLAYER_ROLE_IDLE
		}
		currentStep++
	}

	s.players.EachStartingFromIndex(*iterateFromIndex, stateChangeCallback)
	s.table = *NewTable()
	return nil
}

func (s *State) CanConfirm(player *Player) bool {

	if player.Confirmed {
		return false
	}

	if player.Winner {
		fmt.Println("player is winner")
		return false
	}

	if player.Looser {
		fmt.Println("player is looser")
		return false
	}

	if player.Role == PLAYER_ROLE_IDLE {
		fmt.Println("player is idle")
		return false
	}

	if player.IsAttaker() {
		if true == s.table.AllCardsDefended() {
			return true
		}

		var defender, _ = s.players.Find(func(player *Player) bool {
			return player.IsDefender()
		})

		fmt.Println("Defender confirmation")
		return defender.Confirmed
	}

	if player.IsDefender() {
		return s.table.HasUndefendedCards()
	}

	fmt.Println("Not any condition worked")
	return false
}

func (s *State) FinishGame(playerId string) (bool, error) {
	var player, found = s.GetPlayer(playerId)
	if !found {
		return false, errors.New("Can't find current player")
	}

	if len(player.Hand.Cards) > 0 {
		return false, nil
	}

	if len(s.deck.Cards) > 0 {
		return false, nil
	}

	player.Winner = true

	var winnersCount = s.players.Count(func(player *Player) bool {
		return player.Winner
	})

	if s.players.GetAmount()-winnersCount != 1 {
		return false, nil
	}

	var looser, looserFound = s.players.Find(func(player *Player) bool {
		return len(player.Hand.Cards) > 0
	})

	if !looserFound {
		return false, errors.New("Can't find a looser")
	}

	looser.Looser = true
	s.finished = true
	return true, nil
}

func (s *State) ToSerializable(currentPlayerId string) (*SerializableState, error) {
	var players = []*OtherPlayer{}
	var me *Player

	var index, found = s.players.FindIndex(func(player *Player) bool {
		return player.Id == currentPlayerId
	})
	if !found {
		return nil, errors.New("Can't find current player")
	}

	s.players.EachStartingFromIndex(*index, func(step int, player *Player) {
		if step == 0 {
			me = player
		} else {
			players = append(players, player.ToOtherPlayer())
		}
	})

	return &SerializableState{
		Table:      &s.table,
		Me:         me,
		TrumpCard:  &s.trump,
		Players:    players,
		CanConfirm: s.CanConfirm(me),
		DeckAmount: len(s.deck.Cards),
		Started:    s.started,
		Finished:   s.finished,
	}, nil
}
