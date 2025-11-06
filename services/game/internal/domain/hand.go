// Package domain defines fundamental card and hand primitives used by the game logic.
package domain

// HandStatus indicates current status of a hand.
type HandStatus string

// Hand status constants.
const (
	HandPlaying   HandStatus = "playing"
	HandBust      HandStatus = "bust"
	HandStick     HandStatus = "stick"
	HandBlackjack HandStatus = "blackjack"
	HandFinished  HandStatus = "finished"
)

// Hand represents a player's hand.
type Hand struct {
	Cards  []Card
	Status HandStatus
}

// NewHand creates a new hand with two initial cards.
func NewHand(c1, c2 Card) *Hand {
	h := &Hand{Cards: []Card{c1, c2}, Status: HandPlaying}
	if Score(h.Cards) == 21 {
		h.Status = HandBlackjack
	}
	return h
}

// Add adds a card if hand is still playable.
func (h *Hand) Add(c Card) {
	if h.Status != HandPlaying {
		return
	}
	h.Cards = append(h.Cards, c)
	score := Score(h.Cards)
	if score > 21 {
		h.Status = HandBust
	} else if score == 21 {
		h.Status = HandStick // auto-stick on 21
	}
}

// Stick marks the hand as stick if playing.
func (h *Hand) Stick() {
	if h.Status == HandPlaying {
		h.Status = HandStick
	}
}
