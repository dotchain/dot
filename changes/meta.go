// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Meta wraps a change with some metadata that is maintained as the
// change is merged. This is useful for carrying contexts with
// changes. One example is the current user making the change
type Meta struct {
	Data interface{}
	Change
}

// Merge merges the change preserving the meta data
func (m Meta) Merge(other Change) (otherx, cx Change) {
	if m.Change != nil {
		other, m.Change = m.Change.Merge(other)
	}
	return other, m
}

// Revert reverts the change preserving the meta data
func (m Meta) Revert() Change {
	if m.Change != nil {
		m.Change = m.Change.Revert()
	}
	return m
}

// ReverseMerge implements Custom.ReverseMerge
func (m Meta) ReverseMerge(c Change) (Change, Change) {
	if c != nil {
		m.Change, c = c.Merge(m.Change)
	}
	return c, m
}

// ApplyTo implements Custom.ApplyTo
func (m Meta) ApplyTo(ctx Context, v Value) Value {
	return v.Apply(metacontext{m.Data, ctx}, m.Change)
}

// MetaValue fetches the meta value and the previous context
// associated with the current change context.
func MetaValue(ctx Context) (v interface{}, p Context) {
	if vx := ctx.Value(metakey); vx != nil {
		v = vx.(metacontext).value
		p = vx.(metacontext).parent
	}
	return v, p
}

type metacontext struct {
	value  interface{}
	parent Context
}

type mckey int

var metakey = mckey(0)

func (mc metacontext) Value(key interface{}) interface{} {
	switch {
	case key == metakey:
		return mc
	case mc.parent == nil:
		return nil
	}
	return mc.parent.Value(key)
}
