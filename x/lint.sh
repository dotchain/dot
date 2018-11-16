#!/bin/bash

gometalinter ./... --disable=vet --disable=gotypex --disable=vetshadow --cyclo-over=15
