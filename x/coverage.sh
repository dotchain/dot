#!/usr/bin/env bash

set -e

if [ "$TRAVIS_EVENT_TYPE" == "cron" ]
then
    # Add -race back to this once pq data races get fixed
    GO111MODULE=on go test -v ./stress -tags stress -rounds 20 -clients 5 -iterations 5
    exit $?
fi

for d in $(go list ./... | grep -v vendor | grep -v /cmd/ | grep -v /tools/ | grep -v /demo/ | grep -v /testing/); do
    out=$(echo $d | cut -c21- | sed "s/\//_/g")
    rm -f coverage$out
    # Add -race back to this once pq data races get fixed
    GO111MODULE=on go test -coverprofile=coverage$out -covermode=atomic -tags "integration stress" $d
done
