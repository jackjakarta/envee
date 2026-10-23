#!/bin/bash

VERSION="${1:-dev}"
LDFLAGS="-s -w -X main.version=${VERSION}"

if command -v sha256sum >/dev/null 2>&1; then
    SHA="sha256sum"
else
    SHA="shasum -a 256"
fi

rm -f envee-*.tar.gz checksums.txt

GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-linux-amd64.tar.gz envee && \
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-linux-arm64.tar.gz envee && \
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-darwin-amd64.tar.gz envee && \
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o envee && tar -czf envee-darwin-arm64.tar.gz envee

rm envee

$SHA envee-*.tar.gz > checksums.txt

echo "Build complete!"
echo "Checksums (written to checksums.txt — upload it with the release for install.sh):"
cat checksums.txt
