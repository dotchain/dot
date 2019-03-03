// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

// Notifier implements standard methods to listen for changes and
// notifications.  Streams typically embed this struct
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

// Handler is a generic structure to hold a function that allows
// function pointers to be properly compared.
type Handler struct {
	Handle func()
}
