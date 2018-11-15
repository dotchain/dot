// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reservet.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"github.com/dotchain/dot/ui/input"
	"golang.org/x/net/html"
)

type focusKeyboard struct {
	FocusID int
	input.Keyboard
}

// Keyboard manages a map of keyboard event handlers on a html "document"
type Keyboard map[*html.Node]focusKeyboard

// Update adds or removes a keyboard handler for a node
func (kbd Keyboard) Update(n *html.Node, id int, k input.Keyboard) {
	if id == 0 && k == nil {
		delete(kbd, n)
	} else {
		kbd[n] = focusKeyboard{id, k}
	}
}

// Get returns the keyboard, id associated with a node
func (kbd Keyboard) Get(n *html.Node) (id int, k input.Keyboard) {
	result := kbd[n]
	return result.FocusID, result.Keyboard
}

// Focus returns the current keyboard handler, if any
func (kbd Keyboard) Focus() input.Keyboard {
	var id int
	var handler input.Keyboard

	for _, v := range kbd {
		if v.FocusID > 0 && v.FocusID > id {
			id = v.FocusID
			handler = v.Keyboard
		}
	}
	return handler
}
