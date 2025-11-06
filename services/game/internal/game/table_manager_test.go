package game

import (
    "testing"
)

// Minimal integration-style tests using the public HTTP handlers in main would
// require refactoring main. Instead we directly test Manager behavior here.

func TestTableLifecycle(t *testing.T) {
    mgr := NewManager()
    table := mgr.Create()
    if table.Status != TableWaiting { t.Fatalf("expected waiting, got %s", table.Status) }

    joined, err := mgr.Join(table.ID, "player-1")
    if err != nil { t.Fatalf("join failed: %v", err) }
    if joined.Status != TablePlaying { t.Fatalf("expected playing after join, got %s", joined.Status) }
    snap := Snapshot(joined)
    if len(snap.PlayerCards) != 2 { t.Fatalf("expected 2 player cards, got %d", len(snap.PlayerCards)) }

    // Perform a hit action.
    _, err = mgr.Action(table.ID, "player-1", "hit")
    if err != nil { t.Fatalf("hit failed: %v", err) }

    // Stick should finish table.
    _, err = mgr.Action(table.ID, "player-1", "stick")
    if err != nil { t.Fatalf("stick failed: %v", err) }
    final, _ := mgr.Get(table.ID)
    if final.Status != TableFinished { t.Fatalf("expected finished, got %s", final.Status) }
}

// Example of potential future HTTP test placeholder (shows how we'd construct handlers).
// Removed placeholder HTTP handler test (will add real handler tests after refactor).
