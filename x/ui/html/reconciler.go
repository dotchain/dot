// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"github.com/dotchain/dot/x/ui/dom"
	"golang.org/x/net/html"
)

// Reconciler returns a new reconciler with the provided events and
// keyboard
func Reconciler(events Events, kbd Keyboard) dom.Reconciler {
	if events == nil {
		events = Events{}
	}
	if kbd == nil {
		kbd = Keyboard{}
	}

	return dom.Reconciler(func(tag string, key interface{}) dom.MutableNode {
		n := &html.Node{Type: html.ElementNode, Data: tag}
		if tag == TextTag {
			n.Type = html.TextNode
		}
		return Node{n, events, kbd}
	})
}
