// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package rt

import "github.com/dotchain/dot/changes"

// Formatted implements the basic rich text methods. It stores the
// rich text in a sequence of Segments where each Segment has the text
// value and styles common to the whole segment.  This is an immutable
// type with all the methods returning new value instead of modifying
// things in place
//
// Inline styles are stored as a key value collection.  Block styles
// (such as list type) etc are also key vaue collections but every
// segment has a sequence of block styles to capture the fact that
// block styles can be nested.  For example, a region of text which
// is inside a numbered list which is itself inside a bulleted list
// would have its block styles as:
//
//     Block{{"list": "bulleted"}, {"list": "ordered"}}
//
// A single block may be split across many segments
//
// The actual text requires a changes.Value so that the formatting can
// be usesd with UTF8 or UTF16 strings.  See
// "github.com/dotchain/dot/changes/changes/types" for "S8" and "S16"
// implementations.
type Formatted struct {
	Text     changes.Collection
	Segments []Segment
}

// Equal identifies if two formatted texts are exactly the same
func (f Formatted) Equal(o Formatted) bool {
	if f.Text != o.Text {
		return false
	}
	own := f.Segments
	other := o.Segments

	if len(own) == 0 {
		own = []Segment{{Count: f.Text.Count()}}
	}
	if len(other) == 0 {
		other = []Segment{{Count: o.Text.Count()}}
	}

	if len(own) != len(other) {
		return false
	}

	for kk, seg := range own {
		if !seg.Equal(other[kk]) {
			return false
		}
	}
	return true
}

// Slice removes a sub-section of the string bringing along with it
// the inline and block styles.  Note that Slice does not return any
// styles for empty slices.
func (f Formatted) Slice(offset, count int) Formatted {
	seen, s := 0, []Segment(nil)

	for kk := 0; kk < len(f.Segments) && seen < offset+count; kk++ {
		seg := f.Segments[kk]
		start, end := seen, seen+seg.Count
		if start < offset {
			start = offset
		}
		if end > offset+count {
			end = offset + count
		}
		if start < end {
			seg.Count = end - start
			s = append(s, seg)
		}
		seen += f.Segments[kk].Count
	}

	return Formatted{f.Text.Slice(offset, count), s}
}

// Concat appends the provided strings onto the receiver
func (f Formatted) Concat(args ...Formatted) Formatted {
	text := f.Text
	s := append([]Segment(nil), f.Segments...)
	if len(s) == 0 && text.Count() > 0 {
		s = []Segment{{Count: text.Count()}}
	}
	zero := text.Slice(0, 0)
	for _, arg := range args {
		if arg.Text.Count() == 0 {
			continue
		}

		text = text.ApplyCollection(changes.Splice{text.Count(), zero, arg.Text})
		seg := arg.Segments
		if len(seg) == 0 {
			seg = []Segment{{Count: arg.Text.Count()}}
		}
		if l := len(s); l > 0 && s[l-1].Equal(seg[0]) {
			s[l-1].Count += seg[0].Count
			seg = seg[1:]
		}
		s = append(s, seg...)
	}
	return Formatted{text, s}
}

// Splice removes elements at the provided offset and replaces that
// with the provided replacement
func (f Formatted) Splice(offset, remove int, replace Formatted) Formatted {
	right := f.Slice(offset+remove, f.Text.Count()-offset-remove)
	return f.Slice(0, offset).Concat(replace, right)
}

// UpdateSlice replaces the identified slice (offset and count) with
// the value returned by the provided function.  The function is
// called with the slice as a parameter.
func (f Formatted) UpdateSlice(offset, count int, fn func(Formatted) Formatted) Formatted {
	return f.Splice(offset, count, fn(f.Slice(offset, count)))
}

// Move shuffles the elements over to the right by the provided
// distance. If the distance is negative it, shuffles left instead
func (f Formatted) Move(offset, count, distance int) Formatted {
	o, c, d := offset, count, distance
	if distance < 0 {
		o, c, d = o+d, -d, c
	}
	x1 := f.Slice(0, o)
	x2 := f.Slice(o, c)
	x3 := f.Slice(o+c, d)
	x4 := f.Slice(o+c+d, f.Text.Count()-o-c-d)
	return x1.Concat(x3, x2, x4)
}

func (f Formatted) mapSegment(fn func(Segment) Segment) Formatted {
	if len(f.Segments) == 0 {
		result := f
		result.Segments = []Segment{fn(Segment{Count: f.Text.Count()})}
		return result
	}

	s := make([]Segment, 0, len(f.Segments))
	for _, seg := range f.Segments {
		seg = fn(seg)
		if l := len(s); l > 0 && s[l-1].Equal(seg) {
			s[l-1].Count += seg.Count
		} else {
			s = append(s, seg)
		}
	}
	return Formatted{f.Text, s}
}

// UpdateInline updates the inline styles of the rich text. Removed
// can be used to specify styles that need to be removed.  Use
// UpdateSlice to modify the styles of a slice.
func (f Formatted) UpdateInline(updated Styles, removed []string) Formatted {
	return f.mapSegment(func(seg Segment) Segment {
		result := seg
		result.Inline = seg.Inline.Update(updated, removed)
		return result
	})
}

// UpdateBlock updates the block styles of the rich text. Removed
// can be used to specify styles that need to be removed.  Use
// UpdateSlice to modify the styles of a slice.
//
// Creating or Deleting a block styles can be done with InsertBlock
// and RemoveBlock respectively.  UpdateBlock expects the block at the
// specified index to exist.
func (f Formatted) UpdateBlock(idx int, updated Styles, removed []string) Formatted {
	return f.mapSegment(func(seg Segment) Segment {
		result := seg
		result.Block = append([]Styles(nil), seg.Block...)
		result.Block[idx] = seg.Block[idx].Update(updated, removed)
		return result
	})
}

// RemoveBlock removes the block specified by the index.  Use
// UpdateSlice to modify the styles of a slice of the rich text.
func (f Formatted) RemoveBlock(idx int) Formatted {
	return f.mapSegment(func(seg Segment) Segment {
		result := seg
		result.Block = append([]Styles(nil), seg.Block[:idx]...)
		result.Block = append(result.Block, seg.Block[idx+1:]...)
		return result
	})
}

// InsertBlock inserts a new block at the provided index. Use
// UpdateSlice to modify the styles of a slice of the rich text.
func (f Formatted) InsertBlock(idx int, styles Styles) Formatted {
	return f.mapSegment(func(seg Segment) Segment {
		result := seg
		result.Block = append([]Styles(nil), seg.Block[:idx]...)
		result.Block = append(result.Block, styles)
		result.Block = append(result.Block, seg.Block[idx:]...)
		return result
	})
}

// Segment captures a slice of the rich text with the same inline and
// block styles.
type Segment struct {
	Count  int
	Inline Styles
	Block  []Styles
}

// Equal checks if two segments have the same inline and block styles
func (s Segment) Equal(o Segment) bool {
	if s.Count != o.Count || len(s.Block) != len(o.Block) || !s.Inline.Equal(o.Inline) {
		return false
	}

	kk := 0
	for kk < len(s.Block) && s.Block[kk].Equal(o.Block[kk]) {
		kk++
	}

	return kk == len(s.Block)
}

// Styles capture both inline and block styles.  The value stored
// in the map must be "comparable"
type Styles map[string]interface{}

// Equal checks if two styles are the same
func (s Styles) Equal(o Styles) bool {
	return s.Contains(o) && o.Contains(s)
}

// Contains checks if all of the styles in the arg are present (and
// the same) in the receiver
func (s Styles) Contains(o Styles) bool {
	for k, v := range o {
		if v2, ok := s[k]; !ok || v2 != v {
			return false
		}
	}
	return true
}

// Update returns a new set of styles merging the current styles with
// the provided styles (the latter takes precedence).  The removed
// list can be used to provide keys that must be deleted.
func (s Styles) Update(o Styles, removed []string) Styles {
	result := Styles{}
	rm := map[string]bool{}
	for _, key := range removed {
		rm[key] = true
	}
	for k, v := range s {
		if _, ok := o[k]; !ok && !rm[k] {
			result[k] = v
		}
	}
	for k, v := range o {
		result[k] = v
	}
	return result
}
