// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ux

// Notifier implements standard methods used to notify mutations
type Notifier struct {
	handlers []*Handler
}

// On registers a handler to be notified on change.
func (n *Notifier) On(h *Handler) {
	n.handlers = append(n.handlers, h)
}

// Off deregisters the handler. Any pending notifications may still be delivered.
func (n *Notifier) Off(h *Handler) {
	for kk, hh := range n.handlers {
		if hh != h {
			continue
		}
		handlers := make([]*Handler, len(n.handlers)-1)
		copy(handlers, n.handlers[:kk])
		copy(handlers[kk:], n.handlers[kk+1:])
		n.handlers = handlers
	}
}

// Notify notifies all registered handlers of a change
func (n *Notifier) Notify() {
	for _, h := range n.handlers {
		h.Handle()
	}
}

type Handler struct {
	Handle func()
}
