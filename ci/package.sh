#!/usr/bin/env bash

set -euo pipefail

GOCACHE="$PWD/go-build"

OUTPUT="$PWD/artifactory"

cd build-system-buildpack
go build -i -ldflags='-s -w' -o bin/package package/main.go
bin/package $OUTPUT