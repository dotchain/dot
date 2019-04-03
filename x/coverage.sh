#!/usr/bin/env bash

set -e

if [ "$TRAVIS_EVENT_TYPE" == "cron" ]
then
    go test -v ./stress -race -tags stress -run TestStressAndReconnect -rounds 100 -clients 5 -iterations 5
    exit $?
fi

for d in $(go list ./... | grep -v vendor | grep -v /cmd/ | grep -v /tools/ | grep -v /demo/ | grep -v /testing/); do
    out=$(echo $d | cut -c21- | sed "s/\//_/g")
    rm -f coverage$out
    go test -race -coverprofile=coverage$out -covermode=atomic -tags "integration stress" $d
done
