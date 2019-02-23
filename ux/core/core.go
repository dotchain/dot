// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package core is the minimal shared UX infrastructure.
//
// This is a thin minimal strongly-typed wrapper around native browser
// DOM implementations. In particular, the imperative nature of the
// native DOM interface is maintained here. The higher order ux
// library has wrappers to help make things less imperative.
//
// All changes to this have to be backwards compatible.
package core

// Driver represents the interface to be implemented by drivers. This
// allows testing in non-browser environments
type Driver interface {
	NewElement(props Props, children ...Element) Element
}

// NewElement creates a new element using the registered driver.
//
// While the children can be specified here, they can also be modified
// via AddChild/RemoveChild APIs
func NewElement(props Props, children ...Element) Element {
	return driver.NewElement(props, children...)
}

// Element represents a raw DOM element to be implemented by a
// driver
type Element interface {
	// SetProp updates the prop to the provided value
	SetProp(key string, value interface{})

	// Value is the equivalent of HTMLInputElement.value
	Value() string

	// Children returns a readonly slice of children
	Children() []Element

	// RemoveChild remove a child element at the provided index
	RemoveChild(index int)

	// InsertChild inserts a child element at the provided index
	InsertChild(index int, elt Element)
}

// Styles represents a set of CSS Styles
type Styles struct {
	Color string
}

// Props represents the props of an element
type Props struct {
	Tag         string
	Checked     bool
	Type        string
	TextContent string
	Styles
	OnChange *EventHandler
}

// ToMap returns the map version of props (useful for diffs)
func (p Props) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Tag":         p.Tag,
		"Checked":     p.Checked,
		"Type":        p.Type,
		"TextContent": p.TextContent,
		"Styles":      p.Styles,
		"OnChange":    p.OnChange,
	}
}

// EventHandler is struct to hold a callback function
//
// This is needed simply to make Props be comparable (which makes it
// easier to see if anything has changed)
type EventHandler struct {
	Handle func(Event)
}

// RegisterDriver allows drivers to register their concrete
// implementation
func RegisterDriver(d Driver) {
	driver = d
}

var driver Driver

// Event is not yet implemented
type Event struct{}

// Change is not yet implemented
type Change interface{}
