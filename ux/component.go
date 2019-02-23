// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

import "github.com/dotchain/dot/ux/core"

// Component is a simple helper class that helps with creating
// declarative components.
//
// It is expected to be embedded directly and all updates are
// specified through the Declare method.
//
// Usage:
//
//      type MyComponent struct {
//           ux.Component,
//           ...
//      }
//
//      func NewMyComponent(....) *MyComponent {
//           result := ....
//           result.Declare(props, children...)
//      }
//
//      func (c *MyComponent) Update(...) {
//           c.Declare(updatedProps, updatedChildren...)
//      }
//
// Note that each component still has to create all its children and
// reuse these from previous runs. If the number of children is fixed,
// a component can simply maintain them via named private fields (and
// remember to update them between calls).  If the number of children
// is dynamic, the component cache can be used (see
// templates/cache.template and todo/tasks_view.go for example usage)
type Component struct {
	Root  core.Element
	props core.Props
}

// Declare creates the component if needed or updates it to reflect
// the new set of props and children
func (c *Component) Declare(props core.Props, children ...core.Element) {
	if c.Root == nil {
		c.Root = core.NewElement(props, children...)
		c.props = props
		return
	}

	if c.props != props {
		before, after := propsToMap(c.props), propsToMap(props)
		c.props = props
		for k, v := range after {
			if before[k] != v {
				c.Root.SetProp(k, v)
			}
		}
	}
	UpdateChildren(c.Root, children)
}

func propsToMap(props core.Props) map[string]interface{} {
	return map[string]interface{}{
		"Tag":         props.Tag,
		"Checked":     props.Checked,
		"Type":        props.Type,
		"TextContent": props.TextContent,
		"Styles":      props.Styles,
		"OnChange":    props.OnChange,
	}
}
