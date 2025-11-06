package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewHandBlackjack ensures a natural 21 is marked as blackjack.
func TestNewHandBlackjack(t *testing.T) {
	h := NewHand(Card{Rank: Ace, Suit: Hearts}, Card{Rank: King, Suit: Spades})
	assert.Equal(t, HandBlackjack, h.Status)
	assert.Equal(t, 21, Score(h.Cards))
	// Adding a card should be ignored due to status not playing.
	before := len(h.Cards)
	h.Add(Card{Rank: Two, Suit: Clubs})
	assert.Equal(t, before, len(h.Cards))
	assert.Equal(t, HandBlackjack, h.Status)
}

// TestHandAddBust covers transition to bust when exceeding 21.
func TestHandAddBust(t *testing.T) {
	h := NewHand(Card{Rank: Ten, Suit: Hearts}, Card{Rank: Six, Suit: Clubs}) // 16
	assert.Equal(t, HandPlaying, h.Status)
	h.Add(Card{Rank: King, Suit: Diamonds}) // 26 -> bust
	assert.Equal(t, HandBust, h.Status)
	assert.Equal(t, 26, Score(h.Cards))
	// Further adds ignored.
	h.Add(Card{Rank: Two, Suit: Spades})
	assert.Equal(t, 3, len(h.Cards))
}

// TestHandAutoStick verifies auto-stick at exactly 21.
func TestHandAutoStick(t *testing.T) {
	h := NewHand(Card{Rank: Ten, Suit: Hearts}, Card{Rank: Five, Suit: Clubs}) // 15
	h.Add(Card{Rank: Six, Suit: Spades})                                       // 21 -> stick
	assert.Equal(t, HandStick, h.Status)
	assert.Equal(t, 21, Score(h.Cards))
}

// TestManualStick ensures manual stick prevents further adds.
func TestManualStick(t *testing.T) {
	h := NewHand(Card{Rank: Ten, Suit: Hearts}, Card{Rank: Seven, Suit: Clubs}) // 17
	h.Stick()
	assert.Equal(t, HandStick, h.Status)
	before := len(h.Cards)
	h.Add(Card{Rank: Two, Suit: Diamonds}) // ignored
	assert.Equal(t, before, len(h.Cards))
	assert.Equal(t, 17, Score(h.Cards))
}

// TestAceAdjustmentAfterAdd tests dynamic Ace downgrade still allows auto-stick.
func TestAceAdjustmentAfterAdd(t *testing.T) {
	h := NewHand(Card{Rank: Ace, Suit: Hearts}, Card{Rank: Nine, Suit: Clubs}) // 20
	assert.Equal(t, HandPlaying, h.Status)
	h.Add(Card{Rank: Ace, Suit: Spades}) // would be 31 -> adjust to 21 -> stick
	assert.Equal(t, HandStick, h.Status)
	assert.Equal(t, 21, Score(h.Cards))
}
