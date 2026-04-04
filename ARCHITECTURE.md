# Architecture

## Goal

`helianthus-ebus-extcap` provides a Wireshark `extcap` entrypoint for passive
ENS capture over `helianthus-ebusd-proxy`.

## Current Bootstrap Layout

- `cmd/helianthus-ebus-extcap`: CLI entrypoint.
- `internal/extcap`: argument parsing and extcap response rendering.
- `internal/capture`: shared record types, filters, and semantic opcode catalog.
- `scripts/`: local CI and terminology gate.

## Planned Runtime Path

1. Wireshark invokes the `extcap` binary.
2. The binary advertises one passive interface backed by
   `helianthus-ebusd-proxy`.
3. Capture mode connects to an ENS northbound endpoint.
4. ENS bytes are normalized into:
   - `ens-events`
   - `ebus-frames`
5. The capture path emits live packets to Wireshark and optionally mirrors them
   to `pcapng`.

## Invariants

- V1 is proxy-only.
- Semantic metadata is attached only for the locked opcode catalog.
- Unknown opcode values fall back to raw display.

