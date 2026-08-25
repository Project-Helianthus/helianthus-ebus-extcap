# AGENTS.md

## Purpose and boundaries

`helianthus-ebus-extcap` is a passive eBUS capture client for Wireshark. It
lists the extcap interface, connects to `helianthus-ebusd-proxy` through the
ENS northbound surface, and shapes `ens-events` and `ebus-frames` for capture.

Keep here the extcap CLI, passive ENS ingestion, read-only capture filtering,
and capture-record/pcapng shaping. Do not add Wireshark dissector code, gateway
runtime or semantic publishing surfaces, or direct adapter-class capture.

## Workflow

1. Reconcile `origin/main`, local changes, related issues, branches, PRs,
   reviews, and checks before work.
2. Use one scoped issue, a branch named `issue/<number>-<slug>` from current
   `main`, and one linked PR.
3. Keep changes narrow; add focused tests when behavior changes. Protocol or
   capture-format changes need RED-first evidence where practical.
4. Run `./scripts/ci_local.sh` before pushing. State the exact command and
   result in the PR, then obtain fresh review for the full PR head.
5. Do not merge without green applicable checks, resolved blocking findings,
   and any required public documentation. Stop at the requested boundary.

`ORCHESTRATOR` and `CO_PILOT` are portable reasoning roles. Use a co-pilot for
planning, bounded implementation, adversarial review, or a second opinion; do
not spend that role on routine reads, searches, polling, or shell inspection.
If it is unavailable, continue with an independent fresh review when the risk
justifies it.

## Safety, privacy, and documentation

This client is capture-only: do not introduce bus writes, adapter control, or
active probing. Live capture, credentials, installation changes, and any action
that can affect a device require explicit operator confirmation at action time.
Never commit credentials, personal identifiers, network coordinates, device
fingerprints, or private/raw captures; use sanitized, minimal fixtures.

Keep the opcode catalog explicit and test-covered. Unknown or unsupported data
must remain raw rather than be speculatively decoded. Capture-format or
Wireshark-facing behavior changes need a compatibility note and, when they
change the public eBUS contract, documentation in
[`helianthus-docs-ebus`](https://github.com/Project-Helianthus/helianthus-docs-ebus).
