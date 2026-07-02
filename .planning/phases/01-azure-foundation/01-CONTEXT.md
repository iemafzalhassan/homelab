# Phase 1: Azure Foundation - Context

**Gathered:** 2026-07-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Azure foundation infrastructure using Terraform, including the root module, VNet, subnets, NSGs, and budget alert.
</domain>

<decisions>
## Implementation Decisions

### State Management
- **D-01:** Local state (simpler, no extra cost, fine for a single-user homelab)
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture & Requirements
- `.planning/PROJECT.md` — Project context and constraints
- `.planning/REQUIREMENTS.md` — Phase 1 requirements (INFRA-01 to INFRA-05)
- `.planning/research/ARCHITECTURE.md` — Networking layout, subnets, resource group name
- `.planning/research/PITFALLS.md` — Azure budget alert necessity
</canonical_refs>
