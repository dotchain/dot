// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package ux implements a flexible UX framework.
//
// A UX component is a strongly typed struct with a minimal set of
// constraints:
//
//
// Root Element Constaint
//
//      1. The struct should expose a Root field of core.Element
//      type.
//
// All UX components should expose their Root element so that the
// parent can effectively manage the collection of children.  This
// field is immutable currently (because making it mutable would
// require a mechanism to notify parents of change, not due to any
// deep architectural requirement).
//
// All actual DOM manipulation should be done using the core
// sub-package which is expected to be highly stable as well as
// relatively strongly typed.
//
// A more declarative flavor can be achieved by components embedding
// the simple.Element struct like so:
//
//        type MyComponent struct {
//             simple.Element
//             ... other fields
//        }
//
// The sub-package defines Element which automatically exposes the
// Root field within. In addition, the props and children can be
// specified to simple.Element in a declarative fashion using the
// Declare method.
//
//
// Constructor And Update
//
//        2. Components should implement an Update method that matches
//        the signature of the constructor.
//
// Components are expected to have a constructor when the root element
// and any private state are created. The signature of the constructor
// can be arbitrary and so is expected to be strongly typed.  When the
// values provided as input to the constructor changes, the component
// is informed via the Update method with the exact same signature.
//
// This constraint allows uniformity of components. A stateful
// component can simply detect changes in the props and apply them. A
// declarative component can update the root element declaratively
// using the simple.Element abstraction.  But managing the children is
// still a bit tricky. The ux/templates package provides a template
// for creating a Cache for each sub-component. This can then be used
// in an almost declarative fashion. See todo/tasks_view.go for an
// example.
//
//
// Events And Notifications
//
//         3. Components should expose the latest value of any mutable
//         field with a streams like interface.  Events are simulated
//         with a strongly typed field for the event.
//
// It is expected that a component will expose changeable fields
// (such as text input) via a linked list.  The specific interface can
// be thought of as a generic type like so:
//
//         type Stream<BaseValueType> struct {
//              Value <BaseBalueType>
//              Change dot.Change
//              Next *Stream<BaseValueType>
//              // private fields
//         }
//
// In addition, for notifications support streams should support a
// simple On and Off methods so callers can watch for when a stream is
// changed.
//
// Since Go does not support generics, a template is available in
// templates/streams.template for code generation as needed. This
// template is rather simplistic and does not do OT-style merging
// though this is definitely possible with the same interface
// more-or-less.
//
//
// Interoperability
//
// The contract between components is rather limited and so
// interoperability is mostly a matter of limiting changes to the core
// module. Multiple versions of the simple package can exist
// side-by-side without any issues.
//
//
// Open Issue Around Animation Refresh
//
// Ideally, it would be good if the DOM updates were all scheduled
// around animation-refresh time. The current setup does not have any
// obvious way of doing this. One option is to build a virtual driver
// that does not apply any changes but caches them until the animation
// framework event and then applies everything in one shot.  The
// trouble with such approaches is the fact that code that queries the
// current DOM elements will likely fail.  This is not yet resolved.
//
//
// Containers
//
// A lot of components are generic containers. The current setup
// requires passing the children down as props and updating them via
// the Update method.
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
