// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import "golang.org/x/net/html"

type eventKey struct {
	*html.Node
	Name string
}

// Events manages a map of events on a html "document"
type Events map[eventKey]func(interface{})

// Update adds or removes a handler for an event
func (e Events) Update(n *html.Node, name string, fn func(interface{})) {
	key := eventKey{n, name}
	if fn == nil {
		delete(e, key)
	} else {
		e[key] = fn
	}
}

// Fire fires an event on the provided node
func (e Events) Fire(n *html.Node, name string, event interface{}) {
	key := eventKey{n, name}
	if fn, ok := e[key]; ok {
		fn(event)
	}
}

// Remove removes all event handlers for a node
func (e Events) Remove(n *html.Node) {
	for k := range e {
		if k.Node == n {
			delete(e, k)
		}
	}
}

func (e Events) forEach(n *html.Node, fn func(key string, v interface{})) {
	for k, v := range e {
		if k.Node == n {
			fn(k.Name, v)
		}
	}
}
