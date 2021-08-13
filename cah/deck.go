package cah

// DrawWhite draws a single white card from the deck.
//
// Calling `DrawWhite()` will panic if the deck has no white cards left.
func (d *Deck) DrawWhite() White {
	return d.DrawWhites(1)[0]
}

// DrawWhites draws `n` white cards from the deck.
//
// Calling `DrawWhites(n)` with `n > len(deck.White)` will panic.
func (d *Deck) DrawWhites(n int) []White {
	cards := d.White[:n]
	d.White = d.White[n:]
	return cards
}

// DrawBlack draws a single black card from the deck.
//
// Calling `DrawBlack()` when the deck has no black cards left will panic.
func (d *Deck) DrawBlack() Black {
	return d.DrawBlacks(1)[0]
}

// DrawBlacks will draw `n` black cards from the deck.
//
// Calling `DrawBlacks(n)` with `n > len(deck.Black)` will panic.
func (d *Deck) DrawBlacks(n int) []Black {
	cards := d.Black[:n]
	d.Black = d.Black[n:]
	return cards
}
