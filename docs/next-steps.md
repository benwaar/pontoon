# Next steps (with proto / contracts)

## Phase 1 – Local Setup

Established a fully scripted local stack: custom Postgres (role `game`, DB `pontoon`) + custom Keycloak (realm auto-import with seeded dev user & client), port conflict resolution (Keycloak 8081, Postgres 55432), health verification (containers + DB SELECT 1), and helper scripts for build/up/down/reset/realm configuration, psql access, cleanup. This gives a reproducible baseline for gameplay and auth work.

See detailed recap: [Phase 1 Expanded](./phase1-expanded.md).

---

## Phase 2 – Game logic

Focus: implement core Pontoon round (deck, shuffle, scoring, hit/stick flow), persist finished hands, secure endpoints with JWT, and stub AI move suggestion. Introduce migrations, logging, metrics, and initial CI. This sets the foundation for auth enforcement & leaderboard work.

Tooling add-on: adopt Dev Container (`.devcontainer/devcontainer.json`) to lock Go/Python versions, add `postCreateCommand` for health check, and set a non-root `remoteUser` for consistent collaborative environment.

See detailed scope: [Phase 2 Expanded](./phase2-expanded.md).

### 🛠 Tooling Tips (Phase 2)
- 🧪 Testing: introduce `testify` for richer assertions and table-driven tests.
- 🛠 Linting: add `golangci-lint` (local + CI) for aggregated static analysis.
- 🗄 Migrations: install `golang-migrate` CLI and create `tools/migrate.sh`.
- 🔁 Live Reload (optional): add `air` for rapid Go iteration.
- 📦 Dev Container: enhance with `remoteUser`, `postCreateCommand` (health + optional migrations), pre-install tools.
- 📄 Pre-commit: set up hooks (go fmt, golangci-lint, markdownlint).
- 📊 Observability prep: choose `zerolog` or Go `slog` for structured JSON logs early.

---

## Phase 3 – Local auth (Keycloak)
1. In Keycloak, add a user → log in → get token.
2. In Go, add OIDC middleware using `OIDC_ISSUER` + `OIDC_AUDIENCE`.
3. On first request, upsert player in `players` using token `sub`.
4. ✅ Goal: only authenticated users can create/join tables.

### 🛠 Tooling Tips (Phase 3)
- 🔐 OIDC: use `coreos/go-oidc` for token verification + JWKS caching.
- 🧪 Security Scan: add `gitleaks` pre-commit / CI to catch secrets.
- 🔍 Token Helpers: scripts `tools/get-token.sh` & `tools/jwt-decode.sh`.
- 🛡 Rate Limiting (optional): lightweight limiter middleware for auth endpoints.
- 📦 Dev Container: add OIDC tooling and scripts inside container.

---

## Phase 4 – Scores / Leaderboard
1. Add migrations (golang-migrate or similar) to create:
   - `players`
   - `player_stats`
   - `game_results`
2. On game end, write to `game_results` and update `player_stats`.
3. Expose `GET /api/leaderboard`.
4. ✅ Goal: scoreboard works end-to-end.

### 🛠 Tooling Tips (Phase 4)
- 🗄 Data Access: adopt `sqlc` (type-safe queries) or `gorm` (ORM) – pick one.
- 📈 Metrics: instrument query latency (Prometheus histogram) + counters.
- 🧪 Benchmarks: add Go benchmarks for scoring aggregation (`go test -bench`).
- 📦 Migration Automation: integrate migrations into CI prior to tests.

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

### 🛠 Tooling Tips (Phase 4b)
- 📜 Proto Tooling: introduce `buf` (lint + breaking change detection).
- 🧬 OpenAPI: use `openapi-generator` or `oapi-codegen` for Go types/server stubs.
- 🔁 Codegen Script: `tools/generate.sh` orchestrates proto + OpenAPI generation.
- 🧪 Contract Tests: add a test ensuring generated types match expectations.

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

### 🛠 Tooling Tips (Phase 5)
- 🚦 CI Lint Stage: run golangci-lint, ruff (Python), black, mypy.
- 🗃 Cache: use actions/cache for Go build cache + pip wheels.
- 📊 Coverage: collect Go coverage profile + pytest coverage; optional Codecov upload later.
- 🐳 Image Scan (optional): integrate `trivy` after Docker build.

---

## Phase 6 – Prod-ish auth swap
1. Frontend (Flutter) → switch to Firebase Auth (Google/Facebook/Apple).
2. Go service → change `OIDC_ISSUER` to Firebase/Google Identity Platform.
3. Redeploy containers to Cloud Run.
4. ✅ Goal: prod is “just auth and play,” no Keycloak.

### 🛠 Tooling Tips (Phase 6)
- 🔐 Secrets Mgmt: adopt `sops` or `doppler` for managing non-dev secrets.
- 🌐 Config Strategy: centralize environment variables with validation on startup.
- 📡 Tracing (optional): add OpenTelemetry SDK stubs for future distributed tracing.

---

## Phase 7 – Infra/Terraform (later)
1. Add Terraform modules in `infra/terraform/modules/`:
   - `cloud-run-service`
   - `postgres` / `cloud-sql`
2. Add `infra/terraform/envs/dev/` for project-specific vars.
3. CI: only `terraform fmt` and `terraform validate` (no apply).
4. Later: add GitHub secrets / Workload Identity to allow deploy.
5. ✅ Goal: infra codified, deployments still under your control.

### 🛠 Tooling Tips (Phase 7)
- 🌍 Terraform Quality: add `tflint` + `tfsec` in CI.
- 🧪 Plan Diff (later): use `infracost` for cost estimation (optional).
- 🗺 Graph: generate `terraform graph | dot -Tpng` diagram artifact.
- 🔄 State Strategy: document remote state backend (e.g., GCS bucket) for later production use.

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
* Dev Container enhancement: add `postCreateCommand`, `remoteUser`, optional task to auto-run migrations
* Add `golangci-lint` config (`.golangci.yml`) and pre-commit hook set
* Introduce `sqlc` or `gorm` once schema stabilizes
* Add `buf` + `openapi-generator` scaffolding for contracts
* Add trivy image scan to CI pipeline

### Observability (Later)
* Add basic structured logging & request tracing (OpenTelemetry)
* Add /metrics endpoints (Prometheus) for game & ai services

### Stretch
* WebSocket realtime table updates
* Horizontal scaling considerations (stateless services + shared DB)
* Feature flags for experimental AI assistance
