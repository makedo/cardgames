package cards

type Hand struct {
	Cards []Card `json:"cards"`
}

func (h *Hand) PopCardById(id int) (Card, bool) {
	var newCards []Card
	var popCard Card
	var ok = false
	for _, card := range h.Cards {
		if id != card.Id {
			newCards = append(newCards, card)
		} else {
			ok = true
			popCard = card
		}
	}

	h.Cards = newCards

	return popCard, ok
}

func (h *Hand) GetCardById(id int) (*Card, bool) {
	for _, card := range h.Cards {
		if id == card.Id {
			return &card, true
		}
	}

	return nil, false
}

func (h *Hand) AddCard(card Card) {
	h.Cards = append(h.Cards, card)
}
