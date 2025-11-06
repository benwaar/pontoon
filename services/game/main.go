package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"pontoon/game/internal/game"
)

func main() {
	mgr := game.NewManager()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Create table
	http.HandleFunc("/api/table", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		t := mgr.Create()
		writeJSON(w, game.Snapshot(t))
	})

	// Join table
	http.HandleFunc("/api/table/join", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		tableID := r.URL.Query().Get("id")
		playerID := r.URL.Query().Get("player")
		if tableID == "" || playerID == "" { http.Error(w, "missing id or player", http.StatusBadRequest); return }
		t, err := mgr.Join(tableID, playerID)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		writeJSON(w, game.Snapshot(t))
	})

	// Action endpoint: /api/table/{id}/action
	http.HandleFunc("/api/table/", func(w http.ResponseWriter, r *http.Request) {
		// crude path parsing for /api/table/{id} and /api/table/{id}/action
		if !strings.HasPrefix(r.URL.Path, "/api/table/") { http.NotFound(w, r); return }
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/table/"), "/")
		if len(parts) == 0 { http.NotFound(w, r); return }
		tableID := parts[0]
		if len(parts) == 1 && r.Method == http.MethodGet { // GET /api/table/{id}
			t, err := mgr.Get(tableID)
			if err != nil { http.Error(w, err.Error(), http.StatusNotFound); return }
			writeJSON(w, game.Snapshot(t))
			return
		}
		if len(parts) == 2 && parts[1] == "action" && r.Method == http.MethodPost {
			// read move from JSON body
			var body struct { Move string `json:"move"`; Player string `json:"player"` }
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil { http.Error(w, "invalid json", http.StatusBadRequest); return }
			if body.Move == "" || body.Player == "" { http.Error(w, "missing move or player", http.StatusBadRequest); return }
			t, err := mgr.Action(tableID, body.Player, body.Move)
			if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
			writeJSON(w, game.Snapshot(t))
			return
		}
		http.Error(w, "unsupported endpoint", http.StatusNotFound)
	})

	log.Println("game service listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
