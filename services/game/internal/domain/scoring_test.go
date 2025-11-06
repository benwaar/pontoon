package domain

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreBasicValues(t *testing.T) {
	// Ace + King = 21
	ace := Card{Rank: Ace, Suit: Hearts}
	king := Card{Rank: King, Suit: Spades}
	s := Score([]Card{ace, king})
	assert.Equal(t, 21, s)
}

func TestScoreMultipleAces(t *testing.T) {
	cards := []Card{{Rank: Ace, Suit: Hearts}, {Rank: Ace, Suit: Spades}, {Rank: Nine, Suit: Clubs}}
	s := Score(cards)
	// A(11)+A(1)+9=21
	assert.Equal(t, 21, s)
}

func TestDeckShuffleDeterministic(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()
	seed := int64(42)
	Shuffle(deck1, rand.NewSource(seed))
	Shuffle(deck2, rand.NewSource(seed))
	assert.Equal(t, deck1, deck2, "Decks should match with same seed")
}
