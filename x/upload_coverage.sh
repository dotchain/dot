#!/usr/bin/env bash

set -e

for d in $(go list ./... | grep -v vendor | grep -v /cmd/ | grep -v /tools/ ); do
    out=$(echo $d | cut -c21- | sed "s/\//_/g")
    file=coverage$out
    if [ -f $file ]; then
	# remove the ten letter coverage prefix from file name
	mv $file coverage.txt
	echo processing $file $out
	bash <(curl -s https://codecov.io/bash) -F $out
	rm -f coverage.txt
    fi
done