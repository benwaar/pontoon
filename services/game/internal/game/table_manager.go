// Package game contains in-memory table management for a simplified blackjack-style game.
package game

import (
	"errors"
	"pontoon/game/internal/domain"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TableStatus represents current lifecycle of a table.
type TableStatus string

// Table status values.
const (
	TableWaiting  TableStatus = "waiting"  // waiting for first player to join
	TablePlaying  TableStatus = "playing"  // active hand in progress
	TableFinished TableStatus = "finished" // hand finished (player stuck/bust and dealer resolved)
)

// Table models a single-player vs dealer game for now.
type Table struct {
	ID          string
	CreatedAt   time.Time
	Status      TableStatus
	PlayerID    string       // simplistic single player; later multi-player slice
	PlayerHand  *domain.Hand // player's hand
	DealerCards []domain.Card
	Deck        []domain.Card
}

// State contains a serializable snapshot.
type State struct {
	ID           string            `json:"id"`
	Status       TableStatus       `json:"status"`
	PlayerID     string            `json:"playerId,omitempty"`
	PlayerCards  []domain.Card     `json:"playerCards,omitempty"`
	PlayerStatus domain.HandStatus `json:"playerStatus,omitempty"`
	DealerCards  []domain.Card     `json:"dealerCards,omitempty"`
	PlayerScore  int               `json:"playerScore"`
	DealerScore  int               `json:"dealerScore"`
}

// Manager holds tables in-memory.
type Manager struct {
	mu     sync.RWMutex
	tables map[string]*Table
}

// NewManager constructs a new in-memory table manager.
func NewManager() *Manager {
	return &Manager{tables: make(map[string]*Table)}
}

// Create creates a new table (no player yet).
func (m *Manager) Create() *Table {
	t := &Table{ID: uuid.NewString(), CreatedAt: time.Now(), Status: TableWaiting}
	m.mu.Lock()
	m.tables[t.ID] = t
	m.mu.Unlock()
	return t
}

// Join assigns the single player and deals initial cards.
func (m *Manager) Join(id, playerID string) (*Table, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tables[id]
	if !ok {
		return nil, errors.New("table not found")
	}
	if t.PlayerID != "" {
		return nil, errors.New("player already joined")
	}
	// Init deck & deal 2 each.
	deck := domain.NewDeck()
	domain.Shuffle(deck, nil)
	t.Deck = deck
	p1, p2 := t.Deck[0], t.Deck[1]
	d1, d2 := t.Deck[2], t.Deck[3]
	t.Deck = t.Deck[4:]
	t.PlayerID = playerID
	t.PlayerHand = domain.NewHand(p1, p2)
	t.DealerCards = []domain.Card{d1, d2}
	t.Status = TablePlaying
	// If player has blackjack, immediately finish dealer if needed.
	if t.PlayerHand.Status == domain.HandBlackjack {
		m.finishDealer(t)
	}
	return t, nil
}

// Action applies a move (hit|stick).
func (m *Manager) Action(id, playerID, move string) (*Table, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tables[id]
	if !ok {
		return nil, errors.New("table not found")
	}
	if t.Status != TablePlaying {
		return nil, errors.New("table not playing")
	}
	if t.PlayerID != playerID {
		return nil, errors.New("wrong player")
	}
	if t.PlayerHand == nil {
		return nil, errors.New("hand not initialized")
	}
	switch move {
	case "hit":
		if len(t.Deck) == 0 {
			return nil, errors.New("deck empty")
		}
		c := t.Deck[0]
		t.Deck = t.Deck[1:]
		t.PlayerHand.Add(c)
		if t.PlayerHand.Status == domain.HandBust || t.PlayerHand.Status == domain.HandStick {
			m.finishDealer(t)
		}
	case "stick":
		t.PlayerHand.Stick()
		m.finishDealer(t)
	default:
		return nil, errors.New("invalid move")
	}
	return t, nil
}

// finishDealer resolves dealer drawing rules then marks table finished.
func (m *Manager) finishDealer(t *Table) {
	// Dealer draws until score >=17 (typical blackjack rule) or deck empty.
	for domain.Score(t.DealerCards) < 17 && len(t.Deck) > 0 {
		t.DealerCards = append(t.DealerCards, t.Deck[0])
		t.Deck = t.Deck[1:]
	}
	t.Status = TableFinished
}

// Get returns current table state.
func (m *Manager) Get(id string) (*Table, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tables[id]
	if !ok {
		return nil, errors.New("table not found")
	}
	return t, nil
}

// Snapshot builds a serializable state.
func Snapshot(t *Table) State {
	playerScore := 0
	ps := domain.HandPlaying
	if t.PlayerHand != nil {
		playerScore = domain.Score(t.PlayerHand.Cards)
		ps = t.PlayerHand.Status
	}
	dealerScore := domain.Score(t.DealerCards)
	return State{
		ID:       t.ID,
		Status:   t.Status,
		PlayerID: t.PlayerID,
		PlayerCards: func() []domain.Card {
			if t.PlayerHand != nil {
				return t.PlayerHand.Cards
			}
			return nil
		}(),
		PlayerStatus: ps,
		DealerCards:  t.DealerCards,
		PlayerScore:  playerScore,
		DealerScore:  dealerScore,
	}
}
