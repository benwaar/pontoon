// Package domain defines fundamental card and hand primitives used by the game logic.
package domain

import (
	"math/rand"
	"time"
)

// Suit represents a card suit.
type Suit string

// Suit constants.
const (
	Clubs    Suit = "C"
	Diamonds Suit = "D"
	Hearts   Suit = "H"
	Spades   Suit = "S"
)

// Rank represents card rank.
type Rank string

// Rank constants.
const (
	Ace   Rank = "A"
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "10"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
)

// Card represents a single playing card.
type Card struct {
	Rank Rank
	Suit Suit
}

// Value returns the base value (A handled specially in scoring).
func (c Card) Value() int {
	switch c.Rank {
	case Ace:
		return 11 // treat as 11 initially; scoring will downgrade if needed
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten, Jack, Queen, King:
		return 10
	default:
		return 0
	}
}

// NewDeck returns a new 52 card deck (no jokers).
func NewDeck() []Card {
	suits := []Suit{Clubs, Diamonds, Hearts, Spades}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
	deck := make([]Card, 0, 52)
	for _, s := range suits {
		for _, r := range ranks {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}
	return deck
}

// Shuffle shuffles the deck in-place using provided rand.Source (if nil, uses time seed).
func Shuffle(deck []Card, src rand.Source) {
	if src == nil {
		src = rand.NewSource(time.Now().UnixNano())
	}
	r := rand.New(src)
	for i := len(deck) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}
