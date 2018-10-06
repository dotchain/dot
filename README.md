# Operational Transforms Package

[![Status](https://travis-ci.org/dotchain/dot.svg?branch=master)](https://travis-ci.org/dotchain/dot?branch=master)
[![GoDoc](https://godoc.org/github.com/dotchain/dot?status.svg)](https://godoc.org/github.com/dotchain/dot)
[![codecov](https://codecov.io/gh/dotchain/dot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotchain/dot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotchain/dot)](https://goreportcard.com/report/github.com/dotchain/dot)

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Documentation

The DOT project is a blend of [Operational
Transformation](https://en.wikipedia.org/wiki/Operational_transformation),
[Persistent
Datastructures](https://en.wikipedia.org/wiki/Persistent_data_structure)
and [reactive](https://en.wikipedia.org/wiki/Reactive_programming)
stream processing.

## Features

1. Small, well tested mutations that compose for rich JSON-like values
2. Immutable, Persistent types for ease of use
3. Rich builtin undo support
4. Folding (committed changes on top of uncommitted changes)
5. Strong references support that are automatically updated with changes
6. Streams and Git-like branching, merging support
7. Customizable rich types for values and changes
8. Simple network support (Gob serialization)

## Demos

See [Demos](https://github.com/dotchain/demos).

### Project status

The whole project is in a refactoring state.  Work items needed before
active release:

1. Migrate storage solutions
2. Complete streams (references list + ChildOf etc)
3. Complete browser-based demo
4. Add performance and stress
