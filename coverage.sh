#!/usr/bin/env bash

set -e

for d in $(go list ./... | grep -v vendor | grep -v /cmd/ | grep -v /tools/ | grep -v /demo/ | grep -v /testing/); do
    out=$(echo $d | cut -c21- | sed "s/\//_/g")
    rm -f coverage$out
    go test -coverprofile=coverage$out -covermode=atomic $d
done