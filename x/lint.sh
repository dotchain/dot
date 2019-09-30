#!/bin/bash

export GO111MODULE=on

curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.19.1

$(go env GOPATH)/bin/golangci-lint run -E goimports -E gosec -E maligned -E misspell -E nakedret -E unconvert -E gocritic -E errcheck $*
