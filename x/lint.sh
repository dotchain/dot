#!/bin/bash

golangci-lint run -E goimports -E gosec -E interfacer -E maligned -E misspell -E nakedret -E unconvert -E gocritic -E errcheck $*
