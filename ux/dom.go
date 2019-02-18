// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dom implements simple dom widgets
package dom

// Driver represents the interface to be implemented by drivers
type Driver interface {
	NewElement(props Props, children ...Element) Element
}

// Element represents a raw DOM element to be implemented by a
// driver
type Element interface {
	// SetProp updates the prop to the provided value
	SetProp(key string, value interface{})

	// Value is the equivalent of HTMLInputElement.value
	Value() string
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
	OnClick  func(MouseEvent)
	OnChange func(Event)
}

// RegisterDriver allows drivers to register their concrete
// implementation
func RegisterDriver(d Driver) {
	driver = d
}

var driver Driver

// Event is not yet implemented
type Event struct{}

// MouseEvent is not yet implemented
type MouseEvent struct{}

// Change is not yet implemented
type Change interface{}
