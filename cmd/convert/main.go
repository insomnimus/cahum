package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/insomnimus/cahum/cah"
)

var (
	in  = flag.String("i", "cards.json", "input file")
	out = flag.String("o", "decks.json", "output file")
)

func main() {
	flag.Parse()

	f, err := os.Open(*in)
	if err != nil {
		panic(err)
	}

	decks := make([]cah.Deck, 0)
	json.NewDecoder(f).Decode(&decks)

	var id uint32
	deckMap := make(cah.Decks)
	for _, deck := range decks {
		for i := range deck.Black {
			deck.Black[i].ID = id
			id++
		}
		for i := range deck.White {
			deck.White[i].ID = id
			id++
		}
		deckMap[deck.Name] = deck
	}

	buf, err := json.MarshalIndent(deckMap, "", "\t")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(*out, buf, 0o644); err != nil {
		panic(err)
	}
}
