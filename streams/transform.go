// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import "github.com/dotchain/dot/changes"

// Transform creates a new stream where changes are modified before
// being appended or before being fetched.
//
// Both next and append are optional.
//
// Note that nil is allowed both as argument and return for next and
// append.
func Transform(parent Stream, append, next func(changes.Change) changes.Change) Stream {
	if append == nil {
		append = func(c changes.Change) changes.Change { return c }
	}
	if next == nil {
		next = func(c changes.Change) changes.Change { return c }
	}
	return transform{parent, append, next}
}

type transform struct {
	Stream
	append, next func(changes.Change) changes.Change
}

func (t transform) Next() (Stream, changes.Change) {
	next, nextc := t.Stream.Next()
	if next == nil {
		return nil, nil
	}

	nextc = t.next(nextc)
	return transform{next, t.append, t.next}, nextc
}

func (t transform) Append(c changes.Change) Stream {
	if c = t.append(c); c != nil {
		return transform{t.Stream.Append(c), t.append, t.next}
	}
	return t
}

func (t transform) ReverseAppend(c changes.Change) Stream {
	return transform{t.Stream.ReverseAppend(c), t.append, t.next}
}
