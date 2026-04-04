# helianthus-ebus-extcap

`helianthus-ebus-extcap` is the Helianthus Wireshark `extcap` client for
passive ENS capture over `helianthus-ebusd-proxy`.

## Purpose and Scope

### What belongs in this repository

- Wireshark `extcap` CLI entrypoint for live passive capture.
- ENS session ingestion from `helianthus-ebusd-proxy`.
- Capture record shaping for `ens-events` and `ebus-frames`.
- Optional offline export scaffolding for `pcapng`.
- Filter parsing for read-only capture selectors (`src`, `dst`, `opcode`,
  `family`).

### What does not belong in this repository

- Wireshark dissector code. That belongs in `helianthus-ebus-wireshark`.
- Gateway runtime, semantic publishing, or GraphQL/MCP surfaces.
- Direct adapter-class passive capture support.

## Status

- Bootstrap repository.
- CLI surface and opcode catalog are in place.
- Live ENS ingestion and `pcapng` emission are not implemented yet.
- Licensing follows the Wireshark lane: `GPL-2.0-or-later`.

## Quickstart

### Build

```bash
go build ./cmd/helianthus-ebus-extcap
```

### List extcap interfaces

```bash
go run ./cmd/helianthus-ebus-extcap --extcap-interfaces
```

### List extcap config

```bash
go run ./cmd/helianthus-ebus-extcap --extcap-interface helianthus-ebus-proxy --extcap-config
```

### Validate repo

```bash
./scripts/ci_local.sh
```

## Planned V1 Behavior

- Connect passively to `helianthus-ebusd-proxy` over ENS northbound.
- Emit two logical streams:
  - `ens-events`
  - `ebus-frames`
- Mirror live capture to `pcapng`.
- Attach semantic metadata only for the locked v1 opcode catalog.

## Related Repositories

- `helianthus-ebus-wireshark`
- `helianthus-docs-ebus`
- `helianthus-ebus-adapter-proxy`
