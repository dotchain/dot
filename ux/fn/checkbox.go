// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fn

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

//go:generate go run codegen.go - $GOFILE

// Checkbox implements a checkbox control.
//
// codegen: pure
func Checkbox(c *checkboxCtx, styles core.Styles, checked *streams.BoolStream) core.Element {

	var result core.Element

	result = c.Element(
		"root",
		core.Props{
			Tag:     "input",
			Type:    "checkbox",
			Checked: checked.Value,
			Styles:  styles,
			OnChange: &core.EventHandler{func(_ core.Event) {
				v := result.Value() == "on"
				checked = checked.Append(nil, v, true)
			}},
		},
	)
	return result
}
