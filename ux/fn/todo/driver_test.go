// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package todo_test

import (
	"fmt"
	"github.com/dotchain/dot/ux/core"
)

func init() {
	core.RegisterDriver(driver{})
}

type driver struct{}

func (d driver) NewElement(props core.Props, children ...core.Element) core.Element {
	return &element{props, children}
}

type element struct {
	props    core.Props
	children []core.Element
}

func (e *element) String() string {
	s := "div"
	if e.props.Tag != "" {
		s = e.props.Tag
	}

	props := []string{}
	if e.props.Type != "" {
		props = append(props, fmt.Sprint("type:", e.props.Type))
	}
	if e.props.Styles != (core.Styles{}) {
		props = append(props, fmt.Sprint("styles:", e.props.Styles))
	}
	if e.props.Checked {
		props = append(props, "checked")
	}

	s += fmt.Sprint(props) + "("
	for _, child := range e.children {
		s += " " + child.(*element).String()
	}
	s += e.props.TextContent
	return s + ")"
}

func (e *element) SetProp(key string, value interface{}) {
	switch key {
	case "Checked":
		e.props.Checked = value.(bool)
	case "TextContent":
		e.props.TextContent = value.(string)
	case "Styles":
		e.props.Styles = value.(core.Styles)
	case "OnChange":
		e.props.OnChange = value.(*core.EventHandler)
	default:
		panic("Unknown key: " + key)
	}
}

func (e *element) Value() string {
	switch {
	case e.props.Type != "checkbox":
		return e.props.TextContent
	case e.props.Checked:
		return "on"
	}
	return "off"
}

/**
func (e *element) setValue(s string) {
	if e.props.Type == "checkbox" {
		e.props.Checked = s == "on"
	} else {
		e.props.TextContent = s
	}
	if cx := e.props.OnChange; cx != nil {
		cx.Handle(core.Event{})
	}
}
**/

func (e *element) Children() []core.Element {
	return e.children
}

func (e *element) RemoveChild(index int) {
	c := make([]core.Element, len(e.children)-1)
	copy(c, e.children[:index])
	copy(c[index:], e.children[index+1:])
	e.children = c
}

func (e *element) InsertChild(index int, elt core.Element) {
	c := make([]core.Element, len(e.children)+1)
	copy(c, e.children[:index])
	c[index] = elt
	copy(c[index+1:], e.children[index:])
	e.children = c
}
