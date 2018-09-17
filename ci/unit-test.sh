#!/usr/bin/env bash

set -euo pipefail

GOCACHE="$PWD/go-build"

cd build-system-buildpack
go test ./...
