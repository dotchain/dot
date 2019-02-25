// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

// Subs help track multiple subscriptions to notifiers. It
// allows a more declarative use of the Notifiers which support a more
// imperative use.
type Subs struct {
	old, current map[*Notifier]*Handler
}

// Begin starts a round of subscriptions
func (s *Subs) Begin() {
	s.old, s.current = s.current, map[*Notifier]*Handler{}
}

// End ends a round of subscriptions
func (s *Subs) End() {
	for nn, hh := range s.old {
		nn.Off(hh)
	}
	s.old = nil
}

// On adds a subscription. Note that there is no way to remove
// subscriptions directly. Any subscriptions not re-established in the
// current round are automatically dropped, so directly removing
// subscriptions is not supported.
func (s *Subs) On(n *Notifier, fn func()) {
	if old, ok := s.old[n]; ok {
		s.current[n] = old
		delete(s.old, n)
	} else if _, ok := s.current[n]; !ok {
		s.current[n] = &Handler{fn}
		n.On(s.current[n])
	}
	s.current[n].Handle = fn
}
