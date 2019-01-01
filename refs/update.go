// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package refs

import "github.com/dotchain/dot/changes"

// Update implements changes.Change interface. All updates
// of references at the container level can only be done  using this
// interface. Nil is not a valid value for a ref and as such can be
// used to indicate the ref didn't exist before or to delete a ref
//
// Unlike most changes, Update is a bit improper. The value of the
// reference is modified by unrelated changes. For example, deleting
// of an element can invalidate or modify the path. This has an
// unfortunate side-effect: Reverts may not restore the references
// perfectly.
type Update struct {
	Key           interface{}
	Before, After Ref
}

// Revert implements changes.Change
func (u Update) Revert() changes.Change {
	return Update{u.Key, u.After, u.Before}
}

// Merge implements changes.Change
func (u Update) Merge(c changes.Change) (cx, ux changes.Change) {
	switch c := c.(type) {
	case nil:
		return nil, u
	case changes.Replace:
		c.Before = c.Before.Apply(nil, u)
		return c, nil
	case changes.PathChange:
		if len(c.Path) == 0 {
			return u.Merge(c.Change)
		}

		before := u.mergeRef(u.Before, c)
		after := u.mergeRef(u.After, c)
		if after == nil && before == nil {
			return c, nil
		}
		return c, Update{u.Key, before, after}
	case Update:
		return u.merge(c)
	case changes.Custom:
		l, r := c.ReverseMerge(u)
		return r, l
	}
	panic("Unknown change type to merge with")
}

// ReverseMerge implements changes.Custom
func (u Update) ReverseMerge(c changes.Change) (cx, ux changes.Change) {
	switch c := c.(type) {
	case nil:
		return nil, u
	case changes.Replace:
		c.Before = c.Before.Apply(nil, u)
		return c, nil
	case changes.PathChange:
		if len(c.Path) == 0 {
			return u.ReverseMerge(c.Change)
		}

		before := u.mergeRef(u.Before, c)
		after := u.mergeRef(u.After, c)
		if after == nil && before == nil {
			return c, nil
		}
		return c, Update{u.Key, before, after}
	case Update:
		l, r := c.merge(u)
		return r, l
	case changes.Custom:
		l, r := c.Merge(u)
		return r, l
	}
	panic("Unknown change type to merge with")
}

// ApplyTo implements changes.Custom
func (u Update) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	if con, ok := v.(Container); ok {
		return con.Apply(ctx, u)
	}
	panic("Update does not implement ApplyTo")
}

func (u Update) merge(c Update) (cx, ux changes.Change) {
	if u.Key != c.Key {
		return c, u
	}
	return Update{c.Key, u.After, c.After}, nil
}

func (u Update) mergeRef(r Ref, c changes.Change) Ref {
	if r != nil {
		r, _ = r.Merge(c)
	}
	if r == InvalidRef {
		return nil
	}
	return r
}
