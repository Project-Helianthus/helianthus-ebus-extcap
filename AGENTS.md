# AGENTS

This repository is part of the **Helianthus Multi-Protocol HVAC Gateway Platform**.

## Dual-AI Operating Model

All development follows the workspace-root
[`AGENTS.md`](../AGENTS.md).

- Role binding is portable between orchestrator and co-pilot.
- Use the co-pilot for reasoning-heavy work only.
- Keep one issue and one PR in flight for this repository.
- Follow doc-gate for any externally visible contract change.

## Repo-Specific Rules

1. This repository is **proxy-only** for passive capture in v1.
2. Direct adapter-class ENS/ENH passive support must not be introduced without
   paired docs and validation evidence.
3. The capture record contract must stay aligned with
   `helianthus-docs-ebus`.
4. The opcode catalog must stay explicit and test-covered.
5. Wireshark-facing behavior changes require a linked change in
   `helianthus-ebus-wireshark` or an explicit compatibility note.

