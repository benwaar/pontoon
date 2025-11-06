package domain

// Score computes the best blackjack-style score for a set of cards (Aces flexible 11->1).
func Score(cards []Card) int {
	total := 0
	aces := 0
	for _, c := range cards {
		v := c.Value()
		if c.Rank == Ace {
			aces++
		}
		total += v
	}
	// Downgrade aces from 11 to 1 while busting.
	for total > 21 && aces > 0 {
		total -= 10 // convert one Ace (11) to 1
		aces--
	}
	return total
}
