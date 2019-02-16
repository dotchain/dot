// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dom

// ChangeHandler is a simple change handling interface
type ChangeHandler interface {
	HandleChange()
}

// Notifier implements standard methods used to notify mutations
type Notifier struct {
	handlers []ChangeHandler
}

// On registers a handler to be notified on change.
func (n *Notifier) On(h ChangeHandler) {
	n.handlers = append(n.handlers, h)
}

// Off deregisters the handler. Any pending notifications may still be delivered.
func (n *Notifier) Off(h ChangeHandler) {
	for kk, hh := range n.handlers {
		if hh != h {
			continue
		}
		handlers := make([]ChangeHandler, len(n.handlers)-1)
		copy(handlers, n.handlers[:kk])
		copy(handlers[kk:], n.handlers[kk+1:])
		n.handlers = handlers
	}
}

// Notify notifies all registered handlers of a change and unregisters them
func (n *Notifier) Notify() {
	handlers, n.handlers = n.handlers, nil
	for _, h := range handlers {
		h.HandleChange()
	}
}
