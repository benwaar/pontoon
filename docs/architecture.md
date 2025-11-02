# Architecture

## Goals
- Build a simple multiplayer-ish card game (Pontoon).
- Learn auth flows (OIDC, SAML) **locally** using Keycloak.
- Keep production simple: social login (Google/Facebook/Apple) + play.
- Stay mostly provider-agnostic: Docker containers, Terraform, CI.

## High-level

```text
Flutter/Dart (client)
    |
    |  WebSocket / HTTP (with Bearer token)
    v
Go Game Service  <---->  Postgres (scores, players, results)
    |
    | HTTP / events
    v
Python AI Service (optional, slow tasks)

Auth:
- Local: Keycloak (OIDC, SAML experiments)
- Prod: Firebase Auth / social provider → backend verifies JWT
```

## Services

### 1. frontend/
Flutter/Dart app. Talks to the Go game service via HTTP/WS. In local dev, it authenticates against Keycloak and sends the access token to the game service.

### 2. services/game (Go)
- WebSocket room / table management
- Pontoon rules in a `domain/` package (pure Go, no cloud deps)
- DB access through interfaces
- Auth middleware for OIDC (Keycloak locally, Firebase in prod)

Suggested layout:

```text
services/game
├── cmd/server/main.go
├── internal/
│   ├── domain/        # pontoon rules, deck, scoring
│   ├── http/          # REST endpoints
│   ├── ws/            # WebSocket handlers
│   ├── auth/          # OIDC verifier
│   └── storage/       # interfaces + postgres impl
└── go.mod
```

### 3. services/ai (Python / FastAPI)
- Optional helper AI endpoints
- Can call Vertex AI / Gemini / OpenAI
- Keeps HTTP interface so we can swap clouds

### 4. contracts/
- OpenAPI / JSON Schemas / proto files
- Source of truth for client/server
- Codegen step in CI builds Dart/Go/Python models

### 5. infra/
- `docker-compose.yml` for local: keycloak, postgres, game, ai
- `terraform/` for real cloud: provision Cloud Run, DB, service accounts

## Auth story

Local:
1. Flutter → Keycloak → gets token
2. Flutter → Go service (WS/HTTP) with token
3. Go verifies token against Keycloak OIDC endpoint

Prod:
1. Flutter → Firebase Auth (social)
2. Flutter → Go service with token
3. Go verifies Google/Firebase keys

Same Go code, different issuer/audience.

## Data model (minimal)

- `players(id, external_id, display_name, created_at)`
- `player_stats(player_id, games_played, games_won, points_total)`
- `game_results(id, player_id, result, points, played_at)`

Postgres for local; Cloud SQL / other Postgres later.
