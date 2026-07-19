# Phase 12: Kyros — Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in `12-CONTEXT.md` and `12-MAINTAINER.md` — this log preserves the alternatives considered.

**Date:** 2026-07-12
**Phase:** 12-kyros-platform
**Areas discussed:** Search backend (D-01), Public/Private rollout (D-02), Registry GC & quotas (D-03), Security pipeline scope (D-04), Multi-tenancy isolation (D-05)
**Inputs:** `12-CONTEXT.md` (initial), `12-REVIEWS.md` (Antigravity CLI cross-AI review)

---

## Search Backend (resolution of Antigravity review HIGH-1)

| Option | Description | Selected |
|---|---|---|
| Meilisearch | Rust single binary, ~150 MB RAM, sub-50ms search, official Helm chart, easy backup | |
| OpenSearch (lightweight profile) | AWS-aligned fork of ES, official Helm chart, JVM-based but smaller profile exists | |
| ZincSearch | Go binary, ES-compatible API, ~1 GB RAM, single binary; lighter than Meili but less mature UI | |
| Postgres native (pg_trgm + tsvector + ParadeDB) | Zero new component, reuses Postgres pod memory | |

**User's choice:** Typesense (free-text selection after a follow-up prompt to clarify — selected over Meilisearch as it's equally lightweight and more industry-recognized for cloud-native search)
**Notes:** User's original criteria: "lightweight, cloud-native, industry-grade, follow cloud-native standards". Typesense meets all four. Replaces Elasticsearch from the locked tech stack in the initial CONTEXT.md.
**Antigravity recommendation:** Meilisearch. **Final decision:** Typesense (user's choice).

---

## Public/Private Rollout (resolution of Antigravity review HIGH-3)

| Option | Description | Selected |
|---|---|---|
| Public + private from day 1 | Original CONTEXT.md decision; risky for $100 budget | |
| Private-only MVP, public deferred | Safer, smaller blast radius | ✓ |
| Public-only MVP (no private) | Unusual; only makes sense for a public good | |

**User's choice:** Private-only for MVP (inferred from context — user accepted the Antigravity review's recommendation without re-arguing)
**Notes:** No public OCI in MVP. No self-serve org signup. Public repos require post-MVP abuse-mitigation plan. This decision is binding for sub-phases 12.1-12.5.

---

## Registry GC & Storage Quotas (resolution of Antigravity review MED-1)

| Option | Description | Selected |
|---|---|---|
| Manual GC on operator demand | Forget to run it; blobs grow forever | |
| CronJob-based GC with read-only lock + per-org Azure Blob quota enforced in Go API | D-03 | ✓ |
| Switch to Zot registry (native GC) | Re-architect mid-phase, no Zot on the locked stack | |

**User's choice:** CronJob + quota (D-03) — accepted as a follow-on decision from the review
**Notes:** Daily 03:00 IST. Lock-then-wait-then-GC-then-unlock. Per-org quota default 50 GB. Enforced in Go API before upload.

---

## Security Pipeline Scope (resolution of Antigravity review MED-2 + HIGH-2)

| Option | Description | Selected |
|---|---|---|
| All 5 security features from day 1 (Trivy + Grype + Syft + Cosign + SLSA L2) | Original CONTEXT.md decision; aggressive | |
| Trivy + Syft + Cosign keyless only (drop Grype + SLSA L3+) | D-04 | ✓ |
| Just Trivy + Syft (no signing) | Too thin for "trusted supply chain" branding | |

**User's choice:** Trivy + Syft + Cosign keyless only (D-04) — accepted as a follow-on
**Notes:** Grype dual-scanner deferred to post-MVP. SLSA L2 only for MVP. Cosign needs outbound HTTPS to Sigstore hosts (egress allowlist needed in NSG).

---

## Multi-Tenancy Isolation (resolution of Antigravity review MED-3)

| Option | Description | Selected |
|---|---|---|
| App-layer `WHERE tenant_id = $1` only | Easy to forget; one missed query = cross-tenant leak | |
| Postgres RLS + app-layer defense-in-depth | D-05; structurally impossible to leak at DB layer | ✓ |
| Schema-per-tenant | Heavy ops burden, not worth it at MVP scale | |

**User's choice:** Postgres RLS + app-layer (D-05) — accepted as a follow-on
**Notes:** Every multi-tenant table has `tenant_id` + RLS policy. Go API sets GUC `app.current_tenant` per txn. App-layer WHERE kept as defense-in-depth.

---

## Maintainer Docs (meta-output)

| Option | Description | Selected |
|---|---|---|
| Rewrite CONTEXT.md only | Standard GSD output | |
| Rewrite CONTEXT.md + add separate 12-MAINTAINER.md | LLM-consumable build sheet | ✓ |
| Skip docs, go straight to plan-phase | User explicitly rejected (asked for "maintainer grade" docs) | |

**User's choice:** Both — CONTEXT.md is the audit/decision log; 12-MAINTAINER.md is the dense, LLM-optimized build sheet
**Notes:** User asked twice for "maintainer grade" docs and noted the LLM would use them to build the MVP. Both artifacts produced.

---

## Claude's Discretion

Areas where Claude used best-judgement because the user timed out on a follow-up question:

- Typesense Helm chart values (single-node, PVC, `requests: cpu=100m memory=128Mi`, `limits: cpu=500m memory=512Mi`)
- GC CronJob schedule (daily 03:00 IST)
- Default per-org quota (50 GB)
- All 4 follow-on decisions (D-02, D-03, D-04, D-05) — accepted by user when presented as "resolution of review HIGH/MED items"

---

## Deferred Ideas

All carried forward from CONTEXT.md:

- Public OCI hosting + self-serve org signup
- Grype dual-scanner
- SLSA L3+ (hermetic builds)
- Marketplace / billing
- Multi-cloud storage
- Multi-region deployment
- Helm Chart repo hosting
- WASM modules + AI Model artifacts
- OCI Policies
- kyros.io production domain
- apko/melange declarative builds
- Public trust score leaderboard
- Webhook subscriptions
