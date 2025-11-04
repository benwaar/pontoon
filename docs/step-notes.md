## Phase 1 – Local skeleton
1. `cd infra && docker compose up --build`
   - keycloak (http://localhost:8080, admin/admin)
   - postgres (game/game, db=pontoon)
   - game-service (Go) → http://localhost:9000/healthz
   - ai-service (Python) → http://localhost:9001/healthz
2. Create Keycloak realm `pontoon`, client `pontoon-api`.

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