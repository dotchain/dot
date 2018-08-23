# Operational Transforms Package

[![Status](https://travis-ci.org/dotchain/dot.svg?branch=master)](https://travis-ci.org/dotchain/dot?branch=master)
[![GoDoc](https://godoc.org/github.com/dotchain/dot?status.svg)](https://godoc.org/github.com/dotchain/dot)
[![codecov](https://codecov.io/gh/dotchain/dot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotchain/dot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotchain/dot)](https://goreportcard.com/report/github.com/dotchain/dot)

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Documentation

This package provides the core stateless
conflict-free transformations for a few composition-friendly
operations on a virtual JSON type (i.e. the type can be composed using
arrays and maps).  Please see the GoDoc reference.

## Status of the project

The project is mostly in active development but all of the core
transformations are quite stable with only minor tweaks to the API
expected going forward.  The support types of Log and ClientLog are a
bit less stable in how they deal with error conditions.

## Client and Server

A native Golang client is available via the [Ver
package](https://godoc.org/github.com/dotchain/ver). 

A native Golang server implementation (with a variety of backend
storage options) is avilable at
[dotjs](https://github.com/dotchain/dots) 

Please see the [stress
test](https://github.com/dotchain/dots/blob/master/tests/journal_stress.go)
for sample end-to-end usage.
