package cards

type Suite uint8
const (
    SPADES Suite = iota + 1
    CLUBS
    HEARTS                
    DIAMONDS     
)
var SuiteList = []Suite{SPADES, CLUBS, HEARTS, DIAMONDS}

type Rank uint8
const MIN_RANK = Rank(2)
const MAX_RANK = Rank(14)

type Card struct {
	Suite Suite `json:"suite"`
	Rank Rank `json:"rank"`
}

func NewCard(suite Suite, rank Rank) *Card {
	return &Card{
		Suite:suite,
		Rank: rank,
	}
}