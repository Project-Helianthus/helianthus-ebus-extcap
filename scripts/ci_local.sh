#!/usr/bin/env bash
set -euo pipefail

./scripts/terminology-gate.sh

unformatted="$(find . -name '*.go' -type f -print0 | xargs -0 gofmt -l)"
if [[ -n "${unformatted}" ]]; then
  echo "gofmt required for:"
  echo "${unformatted}"
  exit 1
fi

GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off go build ./...
