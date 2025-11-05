# Next steps (with proto / contracts)

## Phase 1 – Local Setup

Established a fully scripted local stack: custom Postgres (role `game`, DB `pontoon`) + custom Keycloak (realm auto-import with seeded dev user & client), port conflict resolution (Keycloak 8081, Postgres 55432), health verification (containers + DB SELECT 1), and helper scripts for build/up/down/reset/realm configuration, psql access, cleanup. This gives a reproducible baseline for gameplay and auth work.

See detailed recap: [Phase 1 Expanded](./phase1-expanded.md).

---

## Phase 2 – Game logic

Next [Phase 2 Plan](./phase2-expanded.md).

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

---

## Future Enhancements (Backlog)
Short-term backlog items not yet scheduled to a phase:

### Auth / Realm
* Add roles & role mappings to realm JSON
* Token fetch helper script (`tools/get-token.sh`)

### Data / Persistence
* Initial schema migration scripts for game service
* Automated migrations runner in build pipeline

### Testing / Quality
* Automated tests for health endpoints and auth flow
* Add integration test hitting Keycloak token + protected endpoint

### Tooling / Dev UX
* Improve realm configure script to remove grep exit 141 (parse JSON)
* Make build script optionally run `docker compose pull`

### Observability (Later)
* Add basic structured logging & request tracing (OpenTelemetry)
* Add /metrics endpoints (Prometheus) for game & ai services

### Stretch
* WebSocket realtime table updates
* Horizontal scaling considerations (stateless services + shared DB)
* Feature flags for experimental AI assistance
