#!/usr/bin/env sh

set -e

_arch=$(go env GOARCH)

# tar platform
for os in linux darwin windows; do
  _target="migrate-$os-$_arch.tar.gz"
  CGO_ENABLED=0 GOOS=$os GOARCH=$_arch \
    go build -o backend
  tar czf $_target app.yml backend
done

rm backend
