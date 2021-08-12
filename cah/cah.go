package cah

type Deck struct {
	Name     string  `json:"name"`
	White    []White `json:"white"`
	Black    []Black `json:"black"`
	Official bool    `json:"official"`
}

type White struct {
	Text string `json:"text"`
	ID   uint32 `json:"id"`
	// Pack int    `json:"pack"`
}

type Black struct {
	Text string `json:"text"`
	ID   uint32 `json:"id"`

	// Pick int    `json:"pick"`
	// Pack int    `json:"pack"`
}

type Player struct {
	Name  string  `json:"name"`
	Cards []White `json:"cards,omitempty"`
	Score uint32  `json:"score"`
	ID    uint32  `json:"id"`
}

func (p *Player) RemoveCard(id uint32) {
	for i, c := range p.Cards {
		if c.ID == id {
			p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
			return
		}
	}
}
