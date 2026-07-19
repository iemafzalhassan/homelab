---
phase: 12
reviewers: [antigravity]
reviewed_at: 2026-07-12T01:55:00Z
plans_reviewed: []
notes:
  - "Phase 12 has no PLAN.md files yet (planning has not started). Review covers the LOCKED decisions in 12-CONTEXT.md."
  - "Only Antigravity CLI (agy 1.1.0) was available; other CLIs (gemini/claude/codex/opencode/qwen/cursor) are not installed. Ollama is running at localhost:11434 but has zero models pulled — skipped."
  - "Antigravity reviewer noted it could not access the homelab repo files from its sandbox and downgraded file-grounded claims accordingly. No file:line citations were possible."
  - "Ollama was detected (HTTP 200 on /v1/models) but returned an empty model list, so no local fallback was possible."
---

# Cross-AI Plan Review — Phase 12

## Antigravity Review

> *Reviewer note: I am running in an empty scratch workspace and cannot access the existing homelab repository files. I cannot verify existing configurations (like the Traefik Gateway or Workload Identity integrations). Findings below are treated as open questions or risks based on the `CONTEXT.md` decisions rather than asserted code flaws.*

### 1. Summary
Kyros is an ambitious and well-conceptualized OCI-compliant registry platform with an excellent focus on supply chain security and a unique "Trust Score" differentiator. The architecture patterns (reusing Keycloak, Traefik Gateway API, and LGTM observability) correctly leverage the existing homelab foundation. However, the plan suffers from severe scope bloat and architectural choices that are fundamentally incompatible with the stated homelab constraints (an 8 GB system node with ~4.5 GB remaining headroom). Developing this full feature set (registry, UI, RBAC, NATS, Elasticsearch, dual-scanners, Kaniko builds, and Cosign integration) from scratch in 12 weeks is highly unrealistic for a single DevOps engineer, and hosting it publicly on day one introduces existential risks to the $100 budget.

### 2. Strengths
- **Ecosystem Reuse:** Excellent decision to reuse the existing Keycloak Platform realm, LGTM observability stack, and Traefik Gateway API. This prevents reinventing the wheel for foundational services.
- **Cost-Aware Compute:** Pushing Kaniko build jobs and ephemeral scanners to the scale-to-zero Spot node pool is a perfect architectural fit for the budget constraints.
- **Trust Score Product Vision:** The composite Trust Score is a strong, clear product differentiator that justifies building a custom registry wrapper rather than just using Harbor or GitLab.
- **Standardized Backend:** Choosing `cncf/distribution` backed by Azure Blob Storage provides a robust, proven OCI compliance layer without having to write registry protocol parsing from scratch.

### 3. Concerns

- **HIGH: Compute Headroom & Elasticsearch Incompatibility**
  - *Context:* The system node only has ~4.5 GB of RAM available.
  - *Risk:* The LOCKED tech stack requires PostgreSQL, Redis, NATS JetStream, Trivy/Grype, Next.js, a Go Backend, and **Elasticsearch**. Elasticsearch alone requires minimum 1-2 GB of heap to function reliably and often spikes higher. Running this alongside NATS, Postgres, and in-memory vulnerability DBs (Trivy/Grype) will almost certainly cause OOM kills on the system node, destabilizing ArgoCD and Jenkins.
- **HIGH: Unrealistic MVP Timeline & Scope**
  - *Context:* A 12-week roadmap for 1 DevOps engineer.
  - *Risk:* 12 weeks to build a multi-tenant SaaS, wire up NATS for asynchronous security pipelines, write a dynamic Kaniko spot-node orchestrator, integrate Fulcio/Rekor keyless signing, and build a Trust Score engine is extremely aggressive. Phase 12.3 and 12.4 alone contain multi-month enterprise-grade features.
- **HIGH: Public Access & Budget DoS Vectors**
  - *Context:* "Public + private images from day 1" and Azure Blob Storage backend.
  - *Risk:* Allowing public image pushing/pulling without explicitly defined quotas, rate limiting, and automated abuse moderation makes the registry an immediate target for hosting malware or crypto-miner images. This could cause Azure Blob bandwidth/storage costs to easily blow past the $100 ceiling.
- **MEDIUM: Registry Garbage Collection & Quotas**
  - *Context:* Using `cncf/distribution`.
  - *Risk:* `cncf/distribution` does not perform background garbage collection automatically. Unreferenced blobs and aborted uploads will accumulate indefinitely. There is no mention of a GC cronjob or how org-level storage quotas will be enforced before blobs are written to Azure.
- **MEDIUM: Security Pipeline Operational Overhead**
  - *Context:* Trivy + Grype, Cosign (Sigstore/Fulcio).
  - *Risk:* Dual scanners mean maintaining and downloading two separate CVE vulnerability databases constantly, which uses network bandwidth and memory. Additionally, Fulcio/Rekor keyless signing requires outbound internet access for OIDC identity verification and transparency log entries, which might require specific network egress rules in the homelab.

### 4. Suggestions

- **Swap Elasticsearch for Postgres Full-Text Search (or Meilisearch):**
  - Drop Elasticsearch from the LOCKED tech stack immediately. For an MVP, Postgres `pg_trgm` and full-text search capabilities are more than sufficient for indexing images and tags, and consume zero additional memory footprint. If dedicated search is absolutely required, use a lightweight alternative like Meilisearch or ZincSearch.
- **Phase the Rollout (Private First):**
  - Restrict the MVP to *private* images and invited users only. Defer "Public images" and arbitrary Org sign-ups to a post-MVP phase to mitigate budget-busting abuse vectors.
- **Simplify the Scanner Pipeline:**
  - Drop Grype for the MVP. Run only Trivy to reduce memory overhead, DB download bandwidth, and integration complexity. A dual-scanner setup is an enterprise feature that complicates the MVP Trust Score logic unnecessarily.
- **Implement Explicit GC and Rate Limiting:**
  - Add explicit tasks in Phase 12.2 to configure Traefik rate-limiting middleware for the registry ingress.
  - Design a CronJob that locks the registry to read-only mode, runs the `registry garbage-collect` binary, and unlocks it. (Alternatively, explore if you can use the newer Zot registry which handles GC natively better than Docker Distribution).
- **Define Multi-Tenancy Boundaries:**
  - Explicitly define how PostgreSQL will handle multi-tenancy (e.g., Row-Level Security vs. explicit `org_id` filtering in the Go ORM) to prevent cross-tenant data leaks.

### 5. Risk Assessment
**Overall Risk Level: HIGH**

*Justification:* While the vision is excellent, the architectural choices (specifically Elasticsearch) physically conflict with the homelab's hardware constraints. The MVP roadmap is vastly over-scoped for a 12-week timeframe and a single engineer. Furthermore, launching a publicly writable OCI registry on a strict $100 budget without pre-planned quotas, garbage collection, and moderation is a severe financial risk. Trimming the tech stack and restricting the MVP to private tenants will bring this down to a manageable MEDIUM risk.

---

## Consensus Summary

Only one reviewer was available (Antigravity CLI). Consensus synthesis is therefore not statistically meaningful — the report below surfaces the reviewer's findings as the single signal, with priority flags for items the user should resolve before invoking `/gsd-plan-phase 12`.

### Single-Reviewer Findings (treated as highest-priority since they are unverified by peer)

**HIGH severity (must address before planning):**
1. **Elasticsearch vs 8 GB system node** — physical conflict. The locked stack is incompatible with homelab compute budget. Action: re-evaluate the "LOCKED" tech stack before planning, or commit to a dedicated user-pool node for Kyros and update node-sizing rationale in PROJECT.md.
2. **Public OCI registry day 1 = budget DoS vector** — without quotas, rate limiting, and moderation, the $100 budget is at risk. Action: explicitly defer "public + private" to post-MVP, OR define a strict abuse-mitigation plan.
3. **12-week MVP scope for 1 engineer is aggressive** — Phase 12.3 and 12.4 each look like 6+ week features. Action: either trim MVP, extend timeline, or accept that MVP will be narrower than the 5-feature list.

**MEDIUM severity (should address during planning):**
4. cncf/distribution needs explicit GC cronjob and blob-link config
5. Dual-scanner (Trivy + Grype) is heavy for MVP; pick one
6. Fulcio/Rekor keyless signing needs outbound egress planning (Sigstore public hosts)
7. Multi-tenancy isolation strategy in Postgres (RLS vs app-layer org_id filter) not yet defined

### Open Questions (could not be verified)
- Does the existing system node have any free capacity headroom? (Requires `kubectl top node` against the live cluster — not done here.)
- Is outbound internet for Sigstore/Fulcio currently allowed from cluster workloads? (Need to check NetworkPolicies / NSG egress rules once Phase 9 lands.)
- Are the .planning/PROJECT.md node-sizing figures still accurate, or has the platform already grown past the 4.5 GB estimate?

### Divergent Views
N/A — single reviewer.

### To incorporate feedback into planning
Run:
  /gsd-plan-phase 12 --reviews
This passes `12-REVIEWS.md` back into the planner as a `reviews` input, which will surface the HIGH-severity items as constraints the next PLAN.md drafts must address.
