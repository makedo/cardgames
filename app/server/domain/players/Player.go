package players

type Player struct {
	Id string
	IsReady bool
}

func NewPlayer(id string) *Player {
	return &Player{
		Id:id,
		IsReady: false,
	}
}

func (p *Player) Ready() {
	p.IsReady = true
}