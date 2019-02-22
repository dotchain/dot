# UX

The goals of the UX project here are:

1. Composition friendliness: components should easily be built on top
of other components easily.

2. Isolation: a single component should not be forced to use any
libraries at all, forced to use a specific declarative setup or
immutable types.  Ideally though all components will use shared
libraries for these but the isolation guarantee is based on the
concept that different components can use different versions of
these within the same project without suffering consequences.

3. Strongly typed components: the typical usage of components should
not require type-casting things.  The lack of generics in Go would
force code generation to get around boiler-plate issues and to the
extent possible such code-generation should generate whole structs
that are typically not modified.

4. Clean two-way mutations: where changes can happen due to user
interaction or server-interaction.

## TODO example

The very basic Todo example is
[here](https://github.com/dotchain/dot/tree/master/ux/todo).

It is lacking some features but there is enough to illustrate most of
the points here.

## Proposal

The only shared type between components is a simple mapping onto the
DOM element (and a strongly typed Styles and Props struct associated
tightly with the DOM element)

Every component is expected to be implemented via a single go
struct which can be used to maintain any needed private state. All
components MUST expose a Root DOM element (which for simplicity cannot
change). Parents of components do not use the DOM element except for
assembling them into the children collection of the parent.

Components are created by explicit constructors that take any
necessary props. Every component that takes props should also support
updating the props using the same signature as the constructor.

  TODO: Separate update calls due to server update from those due to
  client updates to enable proper OT-style merging


Components can also provide outputs. An example is a simple TextEdit
control which takes the initial text as a prop but allows users to
edit this. All *output* data can simply be exported by components
however they like but a standard mechanism is for mutable data to use
a *Stream* interface:

    stream.On(func() { ... change notification }) // add listener
    stream.Off(fn) // remove listener
    stream.Value // raw value of a stream (immutable, strongly typed)
    stream.Change  // if the stream value has changed
    stream.Next.Value // the next value on the stream
    stream.Latest() // last entry

The stream interface is a linked list which allows tracking all
changes with notifications.  Components should typically always
maintain the latest value and expose this as most consumers will
typically only want the last value. But any consumer that wants to
track individual changes can simply take snapshots and use that to
walk the list.

Note that a simple streams implementation can be easily generated for
any given base type.

## Events

Components often need to generate callbacks, not store data. The
streams interface works just fine for this for the most part.

Another problem with events is the ability to schedule re-rendering at
animation-refresh time. This can be accomplished by simply deferring
handling on any streams callback on animation-refresh.

  TODO: work out how all components can cooperate in delaying
  re-rendering until animatain refresh.

## Declarative implementation

Some components are much easier and more readably specified with a
declarative syntax.  Since each component is responsible for managing
its own children, one option is to memoize component creation:  on
each update cycle using  the memoizer/cache to fetch-or-create a
child.

The standard **ux** package provide simple ways to reconcile children
nodes and even standard control caches can be implemented using
code-generation. Note that components can use other techniques of
reconciliation so long as they dont impose any additional
requirements on children.

It is possible to also use a React-style DSL for specifying a spect
and have it update a single component and its children (without having
a framework that auto-rerenders).

## Multiple drivers

Currently there is an assumption of a single driver for the whole
package. This is mostly because Go does not have any clean way of
doing injection and so it would force passsing the driver a parameter
to all functions. Given that drivers are expected to be exceedingly
thin implementations, this should be ok.

## Templates

The standard [streams
template](https://github.com/dotchain/dot/blob/master/ux/templates/streams.template)
and [cache
template](https://github.com/dotchain/dot/blob/master/ux/templates/cache.template)
are both relatively straight-forward. Please see
[TODO](https://github.com/dotchain/dot/tree/master/ux/todo) for
example usage.
