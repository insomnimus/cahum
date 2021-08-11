package cah

type Deck struct {
	Name     string  `json:"name"`
	White    []White `json:"white"`
	Black    []Black `json:"black"`
	Official bool    `json:"official"`
}

type White struct {
	Text string `json:"text"`
	Pack int    `json:"pack"`
}

type Black struct {
	Text string `json:"text"`
	Pick int    `json:"pick"`
	Pack int    `json:"pack"`
}

type Player struct {
	Name  string  `json:"name"`
	Cards []White `json:"cards,omitempty"`
	Score uint32  `json:"score"`
}
