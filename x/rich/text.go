// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rich

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// NewText initializes a new rich text value
func NewText(s string, attr ...Attr) *Text {
	attrs := Attrs{}
	for _, val := range attr {
		attrs[val.Name()] = val
	}
	return &Text{{Text: s, Attrs: attrs, Size: types.S16(s).Count()}}
}

// Text represents a rich text
type Text []attrText

type attrText struct {
	Text string
	Attrs
	Size int
}

// Count returns the size of the text
func (t *Text) Count() int {
	sum := 0
	for _, x := range *t {
		sum += x.Size
	}
	return sum
}

// PlainText returns the plain text version of the rich text.
func (t *Text) PlainText() string {
	result := ""
	for _, x := range *t {
		result += x.Text
	}
	return result
}

// Slice returns a text within the specified range
func (t *Text) Slice(offset, count int) changes.Collection {
	seen := 0
	s := Text{}
	for _, x := range *t {
		start, end := seen, seen+x.Size
		if start < offset {
			start = offset
		}
		if end > offset+count {
			end = offset + count
		}
		if diff := end - start; diff > 0 {
			text := types.S16(x.Text).Slice(start-seen, diff).(types.S16)
			s = append(s, attrText{string(text), x.Attrs, diff})
		}
		seen += x.Size
	}
	return &s
}

// Concat joins two rich text values
func (t *Text) Concat(o *Text) *Text {
	c := changes.Splice{Offset: t.Count(), Before: &Text{}, After: o}
	return t.Apply(nil, c).(*Text)
}

// SetAttribute returns a change which when applied would set all
// the attributes in the range [offset, offset+count] to the provided
// value.
func (t *Text) SetAttribute(offset, count int, attr Attr) changes.Change {
	before := t.sliceAttr(offset, count, attr.Name())
	return setattr{offset, attr.Name(), before, values{{attr, count}}}

}

// RemoveAttribute returns a change which when applied would remove the
// named attributes in the range [offset, offset+count]
func (t *Text) RemoveAttribute(offset, count int, name string) changes.Change {
	before := t.sliceAttr(offset, count, name)
	return setattr{offset, name, before, values{{changes.Nil, count}}}

}

// Apply implements changes.Value
func (t *Text) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: t.set, Get: t.get, Splice: t.splice}).
		Apply(ctx, c, t)
}

// ApplyCollection implements changes.Collection
func (t *Text) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	return (types.Generic{Set: t.set, Get: t.get, Splice: t.splice}).
		ApplyCollection(ctx, c, t)
}

func (t *Text) sliceAttr(offset, count int, name string) values {
	seen := 0
	s := values{}
	for _, x := range *t {
		start, end := seen, seen+x.Size
		if start < offset {
			start = offset
		}
		if end > offset+count {
			end = offset + count
		}
		if start < end {
			run := valRun{changes.Nil, end - start}
			if attr, ok := x.Attrs[name]; ok {
				run.Value = attr
			}
			if l := len(s); l > 0 && equalAttr(s[l-1].Value, run.Value) {
				s[l-1].Size += run.Size
			} else {
				s = append(s, run)
			}
		}
	}
	return s
}

func (t *Text) splitString(s string, offset, fullSize int) (left, right string) {
	l := types.S16(s).Slice(0, offset).(types.S16)
	r := types.S16(s).Slice(offset, fullSize-offset).(types.S16)
	return string(l), string(r)
}

func (t *Text) applySetAttr(c setattr) *Text {
	offset, after := c.Offset, c.After
	result := &Text{}
	for _, x := range *t {
		text, size := x.Text, x.Size
		var left string

		if offset > 0 && offset < size {
			left, text = t.splitString(text, offset, size)
			result = result.push(left, x.Attrs, offset)
			size, offset = size-offset, 0
		}

		for size > 0 && offset <= 0 && len(after) > 0 {
			updated := x.Attrs.set(c.Name, after[0].Value).(Attrs)
			if after[0].Size > size {
				result = result.push(text, updated, size)
				next := valRun{after[0].Value, after[0].Size - size}
				after = append(values{next}, after[1:]...)
				size = 0
			} else {
				left, text = t.splitString(text, after[0].Size, size)
				size -= after[0].Size
				result = result.push(left, updated, after[0].Size)
				after = after[1:]
			}
		}

		if size > 0 {
			result = result.push(text, x.Attrs, size)
			offset -= size
		}
	}
	return result
}

// note: get returns the actual attrs for that index
func (t *Text) get(key interface{}) changes.Value {
	idx := key.(int)
	seen := 0
	for _, x := range *t {
		if idx >= seen && idx < seen+x.Size {
			return x.Attrs
		}
		seen += x.Size
	}
	return changes.Nil
}

func (t *Text) set(key interface{}, v changes.Value) changes.Value {
	idx := key.(int)
	mid := t.Slice(idx, 1).(*Text)
	(*mid)[0].Attrs = v.(Attrs)
	return t.splice(idx, 1, mid)
}

func (t *Text) splice(offset, remove int, replacement changes.Collection) changes.Collection {
	left := t.Slice(0, offset).(*Text)
	mid := replacement.(*Text)
	right := t.Slice(offset+remove, t.Count()-offset-remove).(*Text)
	return left.concat(mid).concat(right)
}

func (t *Text) concat(o *Text) *Text {
	if l := len(*t); l > 0 && len(*o) > 0 && (*t)[l-1].Attrs.Equal((*o)[0].Attrs) {
		(*t)[l-1].Size += (*o)[0].Size
		(*t)[l-1].Text += (*o)[0].Text
		x := (*o)[1:]
		o = &x
	}
	result := append(*t, *o...)
	return &result
}

func (t *Text) push(text string, attrs Attrs, size int) *Text {
	if l := len(*t); l > 0 && (*t)[l-1].Attrs.Equal(attrs) {
		(*t)[l-1].Text += text
		(*t)[l-1].Size += size
		return t
	}
	result := append(*t, attrText{text, attrs, size})
	return &result
}
