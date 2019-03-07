# Contributing to DOT

DOT is the core library for operational transformtion.  This is
implemented in Go and is also meant to be the reference
implementation.

The code is mostly idiomatic go.  In addtion, there is a tendency to
write functional code (i.e very low side-effects -- preferring to
return new slices instead of mutating input arg, for example). There
is a lot of immutable types used.

The core library is quite critical -- the whole DOT system will
require upgrades if code here changes in some backwards incomptabile
way. So, this code is intentionally likely to grow slowly or only grow
in an additive process (where backwards compatiblity is not an
issue).

## Documentation

The code is somewhat sparsely documented but pleqse feel free to file
issues for even simple questions.

## Code organization

* The changes directory contains the core **changes** package
which implements the basic transformations that is the basis for
everything else.

* The ops directory implements **operations** and mainly deals with
the network/protocol aspects.

* The streams and refs package implement the concept of *streams* of
data for managing changes.

## Building, testing, linting

While standard `go get -u ./...` and `go test ./...` should work, all
pull requests to this project will be tested against ./x/lint.sh and
./x/coverage.sh.

```
go test --coverprofile=cover.out
go tool cover --html=cover.out
```

Linting is done using [gometalinter](https://github.com/alecthomas/gometalinter) but with
a very specific set of lint rules.  Please run `./x/lint.sh` to lint the project.


```
go get -u github.com/alecthomas/gometalinter
gometalinter --install --update
./x/lint.sh
```

## Filing issues

Please feel free to file issues whether it is a simple matter of
trying understand code or project ideas or if it is an actual bug
report.


"There are no stupid questions."


## Developing

Pull requests are welcome and appreciated.
