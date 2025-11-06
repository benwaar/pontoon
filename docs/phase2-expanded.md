## Phase 2: Gameplay & Platform Foundations

Phase 1 delivered a fully automated local environment (Docker Compose stack, custom Postgres + Keycloak images, deterministic realm import with seeded user/client, health checks, scripts). Phase 2 shifts focus from "environment ready" to a functioning minimal game loop with authenticated players, persistent results, and scaffolding for future AI and realtime features.

### High-Level Goals
1. You can play a full round of Pontoon (Blackjack variant) locally via authenticated HTTP requests.
2. Game state transitions are validated server-side; finished hands persist to Postgres.
3. Auth (Keycloak) tokens protect game endpoints; JWT validation middleware in place.
4. AI service contract stub exists (suggested move endpoint + timeout fallback) for future logic.
5. Foundational CI, migrations, logging, and metrics are established.

### Scope (Included)
- Core domain: deck, shuffle, dealing, scoring, hit/stick flow, hand termination rules.
- REST endpoints for table lifecycle and actions.
- Basic persistence schema (tables: players, tables, hands, hand_cards, moves).
- DB migrations tooling & script.
- JWT verification against Keycloak (no advanced roles yet; just authenticated user).
- AI stub endpoint and internal interface definition.
- CI pipeline (build + lint + tests + health check smoke).
- Observability (structured logs + Prometheus metrics endpoint).

### Out of Scope (Deferred to Phase 3+)
- Realtime WebSocket updates.
- Advanced roles/permissions mapping & admin flows.
- Sophisticated AI heuristics/model integration.
- Leaderboards, matchmaking, in-depth analytics.
- Secrets rotation & production hardening.

### Workstreams

#### 1. Database Migrations & Schema
Contract:
- Tool: `golang-migrate` (CLI) or alternative lightweight approach.
- Script: `tools/migrate.sh` (up by default; flag for down).
- Initial schema objects:
  - players (id, external_auth_id, display_name, created_at)
  - tables (id, status, created_at, owner_player_id)
  - hands (id, table_id, player_id, dealer_score, player_score, outcome, finished_at)
  - hand_cards (id, hand_id, card_rank, card_suit, is_dealer BOOLEAN, dealt_order)
  - moves (id, hand_id, player_id, move_type, created_at)
Edge Cases: multiple rapid joins, aborted hands, invalid foreign keys.

#### 2. Game Domain Implementation (Go)
Existing Draft (from initial notes):
```
1. In `services/game/internal/domain/` implement basic Pontoon:
   - deck / shuffle
   - player hand, dealer hand
   - score calculation
```
Expanded:
- Card model (Rank, Suit) + derived value (Aces flexible: 1 or 11).
- Deck service (seedable for deterministic tests).
- Hand state: list of cards, status (playing|bust|stick|blackjack|finished).
- Scoring rules: treat Aces adaptively; blackjack check on initial deal.
- Validation for moves (cannot hit after bust/stick; enforce turn ordering if multi-player later).
Edge Cases: multiple Aces; dealer hitting logic; immediate blackjack.

#### 3. HTTP API Layer
Endpoints (initial):
- POST `/api/table` → create table
- POST `/api/table/{id}/join` → join as player
- POST `/api/table/{id}/action` (hit|stick)
- GET `/api/table/{id}` → current state snapshot
Deferred: GET `/ws/table/{id}` for realtime
Consider pagination for historical hands later.

#### 4. Auth Integration
- Middleware: verify JWT (Keycloak realm `pontoon`, audience/client `pontoon-api`).
- Extract user sub → map/create player record (lazy provisioning on first hit).
- Reject unauthenticated requests; return 401 with JSON error contract.
Edge Cases: expired token, missing kid, JWKS rotation.

#### 5. AI Service Contract (Stub)
- Define request JSON: `{ "tableId": "uuid", "hand": { ... }, "legalMoves": ["hit","stick"] }`.
- Response JSON: `{ "recommended": "hit", "alternatives": ["hit","stick"], "confidence": 0.62 }`.
- Add internal client with timeout (e.g. 300ms) and fallback (no recommendation → default to manual only).
Testing: simulate mock server in unit tests.

#### 6. Shared Contracts Directory
- `contracts/` with JSON Schemas or Go types exported for AI stub & game API.
- Optional code generation later (proto or OpenAPI).

#### 7. Dev Experience
- Optional: integrate `air` for live reload on Go changes.
- Extend `build.sh` to run migrations (or separate explicit `tools/migrate.sh`).
- Add `tools/jwt-decode.sh` (decode token segments for debugging).
 
##### 🛠 Dev Container Adoption
Adopt `.devcontainer/devcontainer.json` early in Phase 2 to lock consistent toolchains (Go 1.22, Python 3.11) and speed onboarding.

Planned adjustments:
- Add `remoteUser` (non-root) for better file permissions.
- Add `postCreateCommand` to run `./tools/check-health.sh || true` and optionally `./tools/migrate.sh` once created.
- (Optional) Pre-install `air` for live reload and `golang-migrate` CLI inside container.

Benefits:
- Deterministic versions across machines.
- Faster pairing/codespaces trial later.
- Eliminates local Go/Python installation requirement.

Exit criteria: opening repo in VS Code → “Reopen in Container” yields a ready environment where `go test ./...` and `./tools/check-health.sh` succeed.

#### 8. CI Bootstrap
Artifacts: test reports (later).

Revision: Phase 2 expansion draft (devcontainer note added).
- Basic metrics: total tables created, active hands, moves per minute.
- Export `/metrics` (Prometheus client for Go).

#### 10. Security & Hardening (Early Hooks)
- Central config for allowed origins (future CORS).
- Note: leave client secret static for dev; add doc note for rotation strategy.
- Input validation: ensure UUID formats; reject oversized payloads.

#### 11. Documentation Enhancements
- `docs/api.md` describing endpoints, auth requirements, sample requests.
- Update `docs/architecture.md` with domain + flow diagram (deck → hand → scoring → persistence).
- Add quickstart snippet: "Play a round".

### Implementation Order (Lean Path)
1. Migrations + schema
2. Domain (deck, hand, scoring)
3. API endpoints + basic tests
4. JWT auth middleware
5. Persistence integration (hands & moves recording)
6. AI stub + contract
7. CI workflow
8. Metrics + logging polish
9. Docs expansion

### Edge Cases & Considerations
- Ace handling with multiple Aces.
- Simultaneous actions (future concurrency locking strategy — optimistic for now).
- Token expiration mid-session (re-auth vs soft fail).
- Partial failures: DB write failure after move accepted → rollback or mark inconsistent state.
- Deterministic test deck vs production randomization.

### Quick Wins (Do Early for Velocity)
- Replace grep in `tools/configure-realm.sh` with `jq` to eliminate exit code 141 noise.
- Add `/api/health` endpoint returning versions & DB connectivity.
- Seed deterministic deck when `DEBUG_DECK_SEED` env present.

### Assumptions
- Single-player vs dealer initial; multi-player extension later.
- One Postgres instance sufficient for dev & tests.
- AI recommendations advisory only (never auto-applied).

### Success Criteria
- curl sequence can create table, join, deal, perform hit/stick, and observe final persisted result.
- Auth enforced (401 when missing token).
- CI green on pushes.
- Logs & metrics exposed.

### Example Flow (Happy Path)
1. Client obtains Keycloak token.
2. POST /api/table (Authorization: Bearer ...) → returns tableId.
3. POST /api/table/{id}/join → player added (lazy provision).
4. POST /api/table/{id}/action {"move":"hit"} → returns updated state.
5. Repeat until finished → final state includes outcome + persisted hand entry.

### Next Step (Actionable)
Start with migrations + domain scaffolding in `services/game/internal/domain/` and add minimal tests.

---
Revision: Initial Phase 2 expansion draft.