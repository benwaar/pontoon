# Next steps (with proto / contracts)

## Phase 1 – Local skeleton (now)
1. `cd infra && docker compose up --build`
   - keycloak (http://localhost:8080, admin/admin)
   - postgres (game/game, db=pontoon)
   - game-service (Go) → http://localhost:9000/healthz
   - ai-service (Python) → http://localhost:9001/healthz
2. Create Keycloak realm `pontoon`, client `pontoon-api`.
3. ✅ Goal: everything runs locally with hardcoded/dev creds.

---

## Phase 2 – Game logic
1. In `services/game/internal/domain/` implement basic Pontoon:
   - deck / shuffle
   - player hand, dealer hand
   - score calculation
2. Add endpoints:
   - `POST /api/table` → create table
   - `POST /api/table/{id}/join`
   - `POST /api/table/{id}/action` (hit/stick)
   - `GET /api/table/{id}` → current state
   - (later) `GET /ws/table/{id}` for realtime
3. Persist finished hands in Postgres.
4. ✅ Goal: you can play a round locally via HTTP.

---

## Phase 3 – Local auth (Keycloak)
1. In Keycloak, add a user → log in → get token.
2. In Go, add OIDC middleware using `OIDC_ISSUER` + `OIDC_AUDIENCE`.
3. On first request, upsert player in `players` using token `sub`.
4. ✅ Goal: only authenticated users can create/join tables.

---

## Phase 4 – Scores / Leaderboard
1. Add migrations (golang-migrate or similar) to create:
   - `players`
   - `player_stats`
   - `game_results`
2. On game end, write to `game_results` and update `player_stats`.
3. Expose `GET /api/leaderboard`.
4. ✅ Goal: scoreboard works end-to-end.

---

## Phase 4b – API contracts (proto / OpenAPI)
**Why now?** By this point you know what the real endpoints are, so you can freeze them and generate code.

1. Create folders:
   ```text
   contracts/
   ├── openapi/
   │   └── pontoon.yaml
   └── proto/
       └── pontoon.proto
   ```
2. In **OpenAPI** (`contracts/openapi/pontoon.yaml`), define:
   - `POST /api/table`
   - `POST /api/table/{id}/join`
   - `POST /api/table/{id}/action`
   - `GET /api/table/{id}`
   - `GET /api/leaderboard`
   And schemas: `Player`, `Table`, `Card`, `LeaderboardEntry`, `Error`.
3. In **proto** (`contracts/proto/pontoon.proto`), define minimal service for internal calls:
   - `GetState`
   - (later) `PlayAction`
4. Add `tools/generate.sh` (or `make generate`) to:
   - run `protoc` → Go + Python stubs into each service
   - run an OpenAPI generator → Go HTTP types/server
   - (optional) generate Dart models
5. Commit generated files so CI doesn’t need protoc yet.
6. ✅ Goal: Go, Python, Dart all agree on message shapes.

---

## Phase 5 – CI (safe/default)
1. Keep existing `.github/workflows/ci.yml`:
   - checkout
   - `go test ./...` (services/game)
   - install Python deps + `pytest` (services/ai)
   - build Docker images
2. **No secrets yet, no deploy yet.**
3. Later: add a second workflow that runs codegen if `contracts/**` changed.
4. ✅ Goal: every push is at least buildable/testable.

---

## Phase 6 – Prod-ish auth swap
1. Frontend (Flutter) → switch to Firebase Auth (Google/Facebook/Apple).
2. Go service → change `OIDC_ISSUER` to Firebase/Google Identity Platform.
3. Redeploy containers to Cloud Run.
4. ✅ Goal: prod is “just auth and play,” no Keycloak.

---

## Phase 7 – Infra/Terraform (later)
1. Add Terraform modules in `infra/terraform/modules/`:
   - `cloud-run-service`
   - `postgres` / `cloud-sql`
2. Add `infra/terraform/envs/dev/` for project-specific vars.
3. CI: only `terraform fmt` and `terraform validate` (no apply).
4. Later: add GitHub secrets / Workload Identity to allow deploy.
5. ✅ Goal: infra codified, deployments still under your control.
