# Manifesto

## Easy conflict-free data synchronization

Any consumer of the platform should achieve conflict-free data 
synchronization with minimal effort.  In particular, consumers of
the platform should not have to write operational transforms or 
even have a deep knowledge of the technology.

## Use JSON

JSON should be used as the transport protocol.  It should also
be used for the virtual object model.

## Design for interoperability

1. Support Javascript, Python, Go languages with idioms natural 
to each language.
2. Support Browser (JS) clients as well as Server scenarios such
as bots and data integrations.
3. Support thin clients (little to no reconciliation) as well as
fat clients (with reconciliation and powerful undo etc).

Objective-C, Swift and Java are other contenders for languages.

A couple of thorns for interoperability are floating point arithmetic
and unicode representations as each subsystem has a different
natural representation.  These will have to be tightly specified.

## Simple protocols, Clear documentation

## Rich types

Platform consumers should be able to use native types (such as structs
 or classes) and not be restricted to special data types provided by
the platform or be stuck with using a raw JSON object as the state.

## Type composition

Rich types, as described above, should be able to use other rich types
in their fields (or as elements of arrays).  For example, one should
be able to use any array-like containers so long as a 1-1 mapping
exists between this and the virtual JSON array.

There is no particular effort to support inheritance -- composition
is the preferred route.

## Custom encodings

Rich types should be able to provide custom on-the-wire encodings for
performance reasons.  For example, sparse arrays can benefit from
run length encodings.

## Easy Undo/Redo 

Undo/redo should be implemented out of the box.  No special development
effort should be needed. The three principles of Undo/Redo:

1. All clients converge in the face of undo/redo.
2. In the absence of intervening remote operations, undo/redo should be
perfect.
3. In the presence of intervening remote operations that do not conflict,
undo should be perfect but otherwise undo can be a bit more lax.

## Authentication and authorization granularity at model level

The actual authentication/authorization is expected to be custom though
a default implementation of the common authentication/authorization
models should be supported.

## Reference backend implementation with pluggable storage

Out of the box support for MySQL, Postgres and MySQL.  Snapshot storage
support via local file system and S3.

## Well documented limitations and constraints

For example, integrity constraints (such as sum of two separate fields 
cannot be a particular constant when clients operate in parallel) are 
not possible right now so consumers should be comfortable with loose 
consistency.  Another similar issue is around validation of operations.

## Custom mutation types

While developers should not have to do this, document the steps needed to
support custom mutation types beyond the default four.

## Large number of model objects

Since the authentication/authorization granularity is at the level of a
model, all protocols should natively and naturally support multiple models.
In particular, it should not be required to open a new connection for each
model.  The platform should also support the ability to construct rich
models using composition.  This effectively implies that authentication
and authorization should be somewhat built a bit deeper than at the 
connection level.

## Support Agents, Catalog and higher level constructs

The longer term goal is to support agents (which work with some class of
models/clients and use the same protocol).  One such agent would be
a "snapshot" service, for example.

Similarly, most rich applications benefit from the addition of a catalog
(which maintains all objects and does garbage collection on them).

These are less understood but well within the scope of the project
