// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dom implements simple dom widgets
package dom

// Driver represents the interface to be implemented by drivers
type Driver interface {
	Raw(props Props) Raw
}

// Raw represents a raw DOM element to be implemented by a driver
type Raw interface {
	SetProp(key string, value interface{})
}

// Styles represents a set of CSS Styles
type Styles struct {
	Color string
}

// Props represents the props of an element
type Props struct {
	Checked bool
	Styles
	OnClick func(MouseEvent)
}

// RegisterDriver is meant to be called by a driver to register
// itself. There can only be one driver at given time.
func RegisterDriver(d Driver) {
	driver = d
}

var driver Driver

// NYI
type MouseEvent struct{}
type Change interface{}
