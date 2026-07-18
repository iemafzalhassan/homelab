---
status: passed
phase: 12
sub-phase: 12-01
---

# Phase 12 Sub-phase 01 Verification

## Goals Achieved

- CI workflow is defensive: web-lint is a documented no-op (Next 16 has no
  `next lint`), web-build uses `npx` (not `bun --filter`).
- `golangci-lint` config migrated from v1 to v2 schema. 0 issues on full
  `golangci-lint run ./...`.
- Registry auth has a single source of truth: `*auth.Validator` (Keycloak
  JWKS, RS256). The dev-mode-bypass that accepted any token when
  `KYROS_AUTH_SECRET` was empty is removed. The middleware is fail-closed:
  nil validator → 503.
- `httpserver` is real production hardening: `ReadHeaderTimeout: 5s` (gosec
  G112), `middleware.RealIP` removed (staticcheck SA1019, deprecated and
  exploitable).
- All 4 `internal/registry` tests pass, including the new
  `TestAuthMiddleware_NilValidator` regression test.

## Verification Steps

1. `go build ./...` from `kyros/` — exit 0.
2. `go vet ./...` from `kyros/` — exit 0.
3. `go test ./...` from `kyros/` — 1 package with tests, 4 tests pass.
4. `golangci-lint run ./...` from `kyros/` — 0 issues.
5. `npx --no-install next build` from `kyros/apps/web/` — 5 routes built
   (4 static + 1 dynamic + middleware proxy).
6. `helm lint deploy/helm/kyros/` — 1 chart linted, 0 failed (icon warning
   only).

## Outcome

Phase 12-01 foundation is complete. The Kyros monorepo is now in a state
where every future 12-NN plan can assume:
- A green `go test`, `go build`, `golangci-lint`, and `next build`.
- A real, fail-closed OCI registry auth path that won't accept any token
  just because a config var was forgotten.
- A CI that won't silently break on iCloud-synced working directories.

## Self-Check

- [x] Goal achieved: CI is green, auth is production-safe, lint is clean.
- [x] No D-01..D-05 decisions changed.
- [x] No new REQ-IDs added (this is infra, not feature work).
- [x] All 5 worker services (scanner/sbom/signer/builder/notifier) remain
      skeleton for 12-03 to implement.
- [x] The webhook handler in `internal/api/webhooks/registry.go` remains a
      stub for 12-02 to implement.
