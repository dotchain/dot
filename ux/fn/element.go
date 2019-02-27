// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fn

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/simple"
)

//go:generate go run codegen.go - $GOFILE

// Element represents a DOM element
//
// codegen: pure
func Element(c *elementCtx, props core.Props, children ...core.Element) core.Element {
	elt := c.rootElement("root")
	elt.Declare(props, children...)
	return elt.Root
}

// codegen: pure
func rootElement(c *rootElementCtx) *simple.Element {
	return &simple.Element{}
}
