# Pontoon

Learning / playground repo for a realtime card game (Pontoon) with:

- **Dart/Flutter** frontend
- **Go** realtime game service (WebSocket/API)
- **Python** AI/helper service
- **Keycloak** locally for OIDC/SAML experiments
- **Firebase Auth** (or other social) in production
- **Docker + docker-compose** for local (dev containers)
- **Terraform + GitHub Actions** for infra/CI

This repo is meant to be **cloud-agnostic-ish**: local uses Keycloak + Postgres, prod can use GCP (Cloud Run, Firebase Auth, Cloud SQL), but the app code talks to OIDC + Postgres and can be pointed elsewhere.
