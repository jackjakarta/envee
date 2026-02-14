#!/bin/bash

VERSION="${1:-dev}"
LDFLAGS="-s -w -X main.version=${VERSION}"

GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-linux-amd64.tar.gz envee && \
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-linux-arm64.tar.gz envee && \
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-darwin-amd64.tar.gz envee && \
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-darwin-arm64.tar.gz envee

rm envee
echo "Build complete!"
echo "Checksums:"

shasum -a 256 envee-darwin-arm64.tar.gz
shasum -a 256 envee-darwin-amd64.tar.gz

shasum -a 256 envee-linux-arm64.tar.gz
shasum -a 256 envee-linux-amd64.tar.gz
