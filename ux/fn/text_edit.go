// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fn

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

//go:generate go run codegen.go - $GOFILE

// TextEdit implements a text edit control.
//
// codegen: pure
func TextEdit(c *textEditCtx, styles core.Styles, text *streams.TextStream) core.Element {
	var result core.Element

	result = c.Element(
		"root",
		core.Props{
			Tag:         "input",
			Type:        "text",
			TextContent: text.Value,
			Styles:      styles,
			OnChange: &core.EventHandler{func(_ core.Event) {
				text = text.Append(nil, result.Value(), true)
			}},
		},
	)
	return result
}
