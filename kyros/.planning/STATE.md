# Kyros — Project State
<!-- GSD:STATE v1 -->

**Project:** Kyros — The Trusted Software Supply Chain Platform  
**Code:** KYROS  
**Version:** v1.0 (MVP)

## Current Status

**Active Phase:** None (project initialized — ready to start Phase 1)  
**Current Milestone:** Milestone 1 — MVP Foundation & Registry  
**Last Session:** 2026-07-12 — Project initialized, CONTEXT.md captured, GSD project bootstrapped  
**Next Action:** `/gsd-plan-phase 1` to begin Phase 1: Monorepo Foundation

## Session History

| Date | Action | Outcome |
|---|---|---|
| 2026-07-12 | `/gsd-discuss-phase` — architecture discovery | All tech stack + architecture decisions locked in 12-CONTEXT.md |
| 2026-07-12 | `/gsd-new-project` — project initialized | PROJECT.md, REQUIREMENTS.md, ROADMAP.md, STATE.md created |

## Resume File

`.planning/ROADMAP.md` — Phase 1 not started

## Stopped At

Project initialized. No phases executed yet.

## Decisions Index

See [PROJECT.md](PROJECT.md) Key Decisions section for all locked decisions.

Key locked decisions:
- Registry engine: cncf/distribution v3
- Backend: Go
- Frontend: Next.js 15 + ShadCN
- Auth: Keycloak (Platform realm reuse)
- Storage: Azure Blob
- Queue: NATS JetStream
- Build engine: Kaniko
- Trust Score: 6-signal weighted composite (0–100)
- All security pipeline features in MVP (Trivy + Grype + Syft + Cosign + SLSA L2)
