// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package riched

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/x/rich"
)

// Editor wraps rich.Text with selection state and richer APIs
//
// The current cursor position is stored in Focus and Achor, both of
// which must have atleast one element (the index into the root
// rich.Text) but can have more elements in the path when referring to
// embedded objects.  Note that Anchor may be later in the rich.Text
// than Focus.  Start() and End() return the earlier and later
// (respectively) of Focus and Anchor
//
// The current set of inline styles (such as FontWeight) are stored in
// Overrides.
//
// Editor itself is  an immutable type which supports the
// changes.Value interface but it is not meant to be sent across the
// network as only the Text is meant to be collaborative.
//
// All mutation methods return changes instead of mutating the
// underlying state. These can be applied to get the new value.
// Note that navigation is considered a mutation type as it modifies
// the state.  Every mutation method here is available on
// riched.Stream with the same method name but with the return value
// being the next stream instance which makes it easier to use.
type Editor struct {
	Text          *rich.Text
	Focus, Anchor []interface{}
	Overrides     rich.Attrs
}

// NewEditor properly initializes an Editor instance.
func NewEditor(t *rich.Text) *Editor {
	p := []interface{}{0}
	return &Editor{Text: t, Focus: p, Anchor: p}
}

// Apply implements changes.Value
func (e *Editor) Apply(ctx changes.Context, c changes.Change) changes.Value {
	// accepts only changes.Replace, changes.PathChange and
	// changes.ChangeSet
	switch c := c.(type) {
	case nil:
		return e
	case changes.Replace:
		return c.After
	case changes.PathChange:
		if len(c.Path) == 0 {
			return e.Apply(ctx, c.Change)
		}
		if c.Path[0] == "Text" {
			return e.applyText(ctx, c)
		}
		return e.applyNonText(ctx, c)
	}
	return c.(changes.Custom).ApplyTo(ctx, e)
}

func (e *Editor) applyText(ctx changes.Context, c changes.PathChange) changes.Value {
	inner := changes.PathChange{Path: c.Path[1:], Change: c.Change}
	clone := *e
	clone.Text = clone.Text.Apply(ctx, inner).(*rich.Text)

	// patch up focus and anchor based on edits to text
	focus, _ := refs.Path(clone.Focus).Merge(inner)
	fpath, _ := focus.(refs.Path)
	clone.Focus = []interface{}(fpath)
	anchor, _ := refs.Path(clone.Anchor).Merge(inner)
	apath, _ := anchor.(refs.Path)
	clone.Anchor = []interface{}(apath)

	return &clone
}

func (e *Editor) applyNonText(ctx changes.Context, c changes.PathChange) changes.Value {
	inner := changes.PathChange{Path: c.Path[1:], Change: c.Change}
	clone := *e
	switch c.Path[0] {
	case "Focus":
		f := changes.Atomic{Value: clone.Focus}
		clone.Focus = f.Apply(ctx, inner).(changes.Atomic).Value.([]interface{})
	case "Anchor":
		a := changes.Atomic{Value: clone.Anchor}
		clone.Anchor = a.Apply(ctx, inner).(changes.Atomic).Value.([]interface{})
	case "Overrides":
		clone.Overrides = clone.Overrides.Apply(ctx, inner).(rich.Attrs)
	}
	return &clone
}

// SetSelection updates selection state
func (e *Editor) SetSelection(focus, anchor []interface{}) changes.Change {
	fBefore, fAfter := changes.Atomic{Value: e.Focus}, changes.Atomic{Value: focus}
	aBefore, aAfter := changes.Atomic{Value: e.Anchor}, changes.Atomic{Value: anchor}
	return changes.ChangeSet{
		changes.PathChange{
			Path:   []interface{}{"Focus"},
			Change: changes.Replace{Before: fBefore, After: fAfter},
		},
		changes.PathChange{
			Path:   []interface{}{"Anchor"},
			Change: changes.Replace{Before: aBefore, After: aAfter},
		},
	}
}

// SetOverride update an override that is used for text insertion
//
// A negative override can be created by using NoAttribute{"name"} as
// the attribute
func (e *Editor) SetOverride(attr rich.Attr) changes.Change {
	var b changes.Value = changes.Nil
	if v, ok := e.Overrides[attr.Name()]; ok {
		b = v
	}
	if b == attr {
		return nil
	}
	return changes.PathChange{
		Path:   []interface{}{"Overrides", attr.Name()},
		Change: changes.Replace{Before: b, After: attr},
	}
}

// RemoveOverride removes an override
func (e *Editor) RemoveOverride(name string) changes.Change {
	b := e.Overrides[name]
	if b == nil {
		return nil
	}
	return changes.PathChange{
		Path:   []interface{}{"Overrides", name},
		Change: changes.Replace{Before: b, After: changes.Nil},
	}
}

// ClearOverrides removes any overrides if present
func (e *Editor) ClearOverrides() changes.Change {
	if len(e.Overrides) == 0 {
		return nil
	}
	return changes.PathChange{
		Path:   []interface{}{"Overrides"},
		Change: changes.Replace{Before: e.Overrides, After: rich.Attrs{}},
	}
}

// CurrentAttributes returns the attributes at the current selection.
//
// It does not factor in any overrides or remove Embed styles. Use
// CurrentEffectiveAttributes for that.
func (e *Editor) CurrentAttributes() rich.Attrs {
	r, idx := e.resolve(e.Text, e.Focus)
	seen := 0
	for _, x := range *r {
		seen += x.Size
		if seen >= idx {
			return x.Attrs
		}
	}
	return nil
}

// CurrentEffectiveAttributes returns the attributes at the current selection.
//
// It factors in any overrides and removes any "Embed" attributes
func (e *Editor) CurrentEffectiveAttributes() rich.Attrs {
	attrs := rich.Attrs{}
	for k, v := range e.CurrentAttributes() {
		if k != "Embed" {
			attrs[k] = v
		}
	}
	for k, v := range e.Overrides {
		if _, ok := v.(NoAttribute); ok {
			delete(attrs, k)
		} else {
			attrs[k] = v
		}
	}
	return attrs
}

// InsertString inserts a string at the current selection
func (e *Editor) InsertString(s string) changes.Change {
	attrs := []rich.Attr{}
	for _, v := range e.CurrentEffectiveAttributes() {
		attrs = append(attrs, v)
	}
	after := rich.NewText(s, attrs...)
	return changes.PathChange{
		Path: append([]interface{}{"Text"}, e.Focus[:len(e.Focus)-1]...),
		Change: changes.Splice{
			Offset: e.Focus[len(e.Focus)-1].(int),
			Before: &rich.Text{},
			After:  after,
		},
	}
}

func (e *Editor) resolve(r *rich.Text, path []interface{}) (*rich.Text, int) {
	// NYI
	return r, path[0].(int)
}
