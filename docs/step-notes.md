## Phase 1 – Local setup (evolved)
Initial manual steps have been replaced by scripted automation.

### Current Ports & Credentials
| Service | Host Port | Container Port | Notes |
|---------|-----------|----------------|-------|
| Keycloak | 8081 | 8080 | Admin: admin / admin |
| Postgres (dev) | 55432 | 5432 | Roles: postgres, game; DB: pontoon |
| Game Service | 9000 | 9000 | Health: /healthz |
| AI Service | 9001 | 9001 | Health: /healthz |

Postgres host port changed from 5432 -> 55432 to avoid collision with a locally installed Postgres server.

### One-Line Bring Up
From project root:
```
./tools/up.sh
```
This builds (if needed) and starts all containers:
* Extended Postgres (custom init scripts)
* Extended Keycloak (theme placeholder + realm import capable)
* Game service (Go)
* AI service (Python)

### One-Line Tear Down / Reset
```
./tools/down.sh          # stop & remove containers
./tools/db-reset.sh      # drop Postgres volume & restart stack (or --no-up)
./tools/drop-infra.sh    # nuke everything: containers, images, volumes, network
```

### Scripts Overview
| Script | Purpose |
|--------|---------|
| build.sh | Build all service images (now 4 build contexts: game, ai, postgres, keycloak) |
| up.sh | Start stack (network + containers) |
| down.sh | Stop & remove containers + network |
| db-reset.sh | Recreate Postgres volume and restart stack |
| clean-docker.sh | Prune dangling images / cache (add --all for aggressive) |
| psql.sh | Convenience wrapper for container psql (defaults pontoon/postgres) |
| configure-realm.sh | Import realm from `infra/realm-pontoon.json` (use --force to re-import) |
| add-realm-user.sh | (Optional) programmatic user creation; now embedded user makes this mostly redundant |

### Keycloak Realm Automation
Realm definition lives in `infra/realm-pontoon.json`:
* Realm: pontoon
* Client: pontoon-api (confidential, secret: pontoon-secret)
* Embedded Dev User: username `user`, password `user`

Import / re-import:
```
./tools/configure-realm.sh --force
```
Script waits for host port 8081, authenticates with admin/admin, deletes existing realm (force), then imports new one.

### Obtaining a Token (Client Credentials)
```
curl -s -X POST \
   -d 'grant_type=client_credentials' \
   -d 'client_id=pontoon-api' \
   -d 'client_secret=pontoon-secret' \
   http://localhost:8081/realms/pontoon/protocol/openid-connect/token
```

### Embedded Postgres Initialization
Files: `infra/postgres/init/001-role-and-db.sql`
Creates:
* Role game (SUPERUSER for local convenience)
* Database pontoon owned by game

Connect from host:
```
PGPASSWORD=game psql -h localhost -p 55432 -U game -d pontoon
```
Or via helper:
```
./tools/psql.sh pontoon game
```

### Health Checks
```
curl -s http://localhost:9000/healthz
curl -s http://localhost:9001/healthz
```

### Shortcut Reference (conversation triggers)
| Shortcut | Action |
|----------|--------|
| up | ./tools/up.sh |
| down | ./tools/down.sh |
| db-reset | ./tools/db-reset.sh |
| configure-realm | ./tools/configure-realm.sh |
| add-user | ./tools/add-realm-user.sh user user --force-password |
| psql | ./tools/psql.sh pontoon game |
| clean-docker | ./tools/clean-docker.sh |
| drop-infra | ./tools/drop-infra.sh |

### Troubleshooting Notes
| Issue | Cause | Fix |
|-------|-------|-----|
| Cannot connect as game on 5432 | Host Postgres conflict | Use new host port 55432 |
| Realm import exits 141 | Grep broken pipe on verification | (Non-fatal) Improve script by JSON parsing if needed |
| Keycloak not ready | Timing during startup | Script waits up to 60s; re-run if exceeded |


## Keycloak pontoon-api client configuration notes

- Client Type: OpenID Connect (OIDC)
- Access Type: confidential

### Capability configuration
- Client authentication: On
- Authorization: Off
- Standard flow: On
- Direct access grants: On (optional, for dev/testing)
- Implicit flow: Off
- Service accounts roles: On (if using client credentials grant)
- OAuth 2.0 Device Authorization Grant: Off
- OIDC CIBA Grant: Off

### Rationale

### URL and origin settings (local dev)
- Root URL: http://localhost:9000
- Home URL: http://localhost:9000
- Valid redirect URIs: http://localhost:9000/*
- Valid post logout redirect URIs: http://localhost:9000/*
- Web origins: http://localhost:9000

> Adjust ports/URLs if your frontend or API runs on a different port. For multiple services, add each required origin/redirect URI.