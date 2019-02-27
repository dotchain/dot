// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
//
//
// This code is generated by github.com/dotchain/dot/ux/fn/codegen.go

package fn

import (
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/streams"
)

type textEditCtx struct {
	ElementCache
	memoInitialized bool
	memoizedParams  struct {
		styles  core.Styles
		text    *streams.TextStream
		result1 core.Element
	}
}

func (c *textEditCtx) areArgsSame(styles core.Styles, text *streams.TextStream) bool {
	if styles != c.memoizedParams.styles {
		return false
	}
	if text != c.memoizedParams.text {
		return false
	}
	return true
}

func (c *textEditCtx) refreshIfNeeded(styles core.Styles, text *streams.TextStream) (result1 core.Element) {
	if !c.memoInitialized || !c.areArgsSame(styles, text) {
		return c.refresh(styles, text)
	}
	return c.memoizedParams.result1
}

func (c *textEditCtx) refresh(styles core.Styles, text *streams.TextStream) (result1 core.Element) {
	c.memoInitialized = true
	c.memoizedParams.styles, c.memoizedParams.text = styles, text
	c.ElementCache.Begin()
	defer c.ElementCache.End()
	c.memoizedParams.result1 = TextEdit(c, styles, text)
	return c.memoizedParams.result1
}

// TextEditCache is generated from TextEdit.  Please see that for
// documentation
type TextEditCache struct {
	old, current map[interface{}]*textEditCtx
}

// Begin starts the round
func (c *TextEditCache) Begin() {
	c.old, c.current = c.current, map[interface{}]*textEditCtx{}
}

// End ends the round
func (c *TextEditCache) End() {
	// TODO: deliver Close() handlers if they exist
	c.old = nil
}

// TextEdit implements the cache create or fetch method
func (c *TextEditCache) TextEdit(key interface{}, styles core.Styles, text *streams.TextStream) core.Element {
	cOld, ok := c.old[key]

	if ok {
		delete(c.old, key)
	} else {
		cOld = &textEditCtx{}
	}
	c.current[key] = cOld
	return cOld.refreshIfNeeded(styles, text)
}
