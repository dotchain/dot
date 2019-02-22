// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package ux implements the basic UX controls and helpers
//
// A UX component is a strongly typed pointer value. It is created
// with a strongly typed constructor of the shape:
//
//     NewComponent(styles core.Styles, props...)
//
// The styles are intended for the root DOM node of the component. The
// props be any mutable or immutable value (though it is simpler to
// detect changes if these were immutable).
//
// Every UX component exposes the root DOM node via a Root field:
//
//     type MyComponent {
//          Root *core.Element,
//          ...
//     }
//
// This allows the parent to add the root DOM node to the children
// collection.  Components are free to maintain whatever state they
// need and expose any "output" values.
//
// Generated Code
//
// The lack of generics in Go leads to using generated code for useful
// helper types. The templates sub-package has two such examples:
// streams.template and cache.template.  The former implemennts
// streams (see bottom up updates below) and the latter implements a
// cache (which is a form of memoization -- see the Containers section
// below).  Both of these templates are used  in the todo package.
//
// Top Down Updates
//
// Every component must implement a Update method with the same
// signature as the constructor. This allows the parent of the
// component to force changes to the props.  While components can
// have model props that have their own listeners for changes,
// components still have to implement the Update method at the very
// least to support modifying the styles.
//
// Bottom up updates
//
// While components can set up their own mechanisms for changes and
// notifications, this package provides a notifier type to help ease
// this.  Components are encourage to implement streams for mutable
// public state (see Checkbox for sample implementation).  Note that
// events can also be expressed in similar fashion with named strongly
// typed stream objects representing the sequence of events.
//
// Containers
//
// A lot of components are generic containers. The children for these
// containers are best treated as arrays of elements with the
// UpdateChildren helper method providing an efficient way for
// containers to maintain the children collection each time.
//
// A lot of container components also need to manage the
// creating/deletion of child components. It is much simpler to do
// this declaratively and  this can be accomplished with a memoizer.
// The templates/ folder contains a sample template for caching a
// collection of components of a specific type.  Please see the
// TodoTasks example in the todo/ folder for an implementation that
// looks a lot more declarative.
//
// Declarative syntax
//
// The package here offers no help for declarative syntax right now
// but the example use in TodoTasks in sub-package todo is probably
// a step in that direction.
//
// Animation Refresh
//
// How components coordinate to make sure DOM mutations happen only
// during animation-refresh is TBD. One approach is to have every
// component implement a refresh() method which only uses cached
// render state to update the whole DOM.  This will probably tie in
// with the declarative syntax. The lack of dependency injection would
// require the scheduler (which is used to register for refresh) to be
// passed as a prop.
//
// Collaborative components
//
// Collaborative components involve state that can be modified
// top-down as well as bottom-up.  These are a bit tricky to do. The
// built-in streams implemnetation in the templates/ directory does
// not handle merging these changes but it could be extended in a
// straight-forward way based on dot/streams package.
package ux

import "github.com/dotchain/dot/ux/core"

// the following types are derived from core which should change very
// rarely, if at all.

// Driver is a type alias
type Driver = core.Driver
// Element is a type alias
type Element = core.Element
// Styles is a type alias
type Styles = core.Styles
// Props is a type alias
type Props = core.Props
// EventHandler is a type alias
type EventHandler = core.EventHandler
// Event is a type alias
type Event = core.Event
// Change is a type alias
type Change = core.Change

// NewElement is an alias
var NewElement = core.NewElement
// RegisterDriver is an alias
var RegisterDriver = core.RegisterDriver
