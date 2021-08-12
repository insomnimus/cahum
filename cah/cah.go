package cah

type Decks map[string]Deck

type Deck struct {
	Name     string  `json:"name"`
	White    []White `json:"white"`
	Black    []Black `json:"black"`
	Official bool    `json:"official"`
}

type Card struct {
	Text string `json:"text"`
	ID   uint32 `json:"id"`
	Pack int    `json:"pack"`
}

type White Card

type Black struct {
	Card
	Pick int `json:"pick"`
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
