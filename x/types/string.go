// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package types

import (
	"github.com/dotchain/dot/changes"
	"unicode/utf16"
)

// S8 implements a string whose offsets are UTF8 bytes
type S8 string

// Slice implements Value.Slice.  Offset and count is based on UTF8
func (s S8) Slice(offset, count int) changes.Value {
	return s[offset : offset+count]
}

// Count returns the number of UTF8 bytes in the string
func (s S8) Count() int {
	return len(s)
}

// Apply implements Value.Apply
func (s S8) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return s
	case changes.Replace:
		if c.IsDelete {
			return changes.Nil
		}
		return c.After
	case changes.Splice:
		o := c.Offset
		remove := string(c.Before.(S8))
		replace := c.After.(S8)
		return s[:o] + replace + s[o+len(remove):]
	case changes.Move:
		ox, cx, dx := c.Offset, c.Count, c.Distance
		if c.Distance < 0 {
			ox, cx, dx = ox+dx, -dx, cx
		}
		x1, x2, x3 := ox, ox+cx, ox+cx+dx
		return s[:x1] + s[x2:x3] + s[x1:x2] + s[x3:]
	case changes.Custom:
		return c.ApplyTo(s)
	}
	panic("Unknown change type.  Cannot apply")
}

// S16 implements string with offsets and counts referring to utf16
// runes. The UTF16 offsets map to native Javascript string offsets.
type S16 string

// Slice implements changes.Value.Slice.  Offset and count are in
// UTF16 units.
func (s S16) Slice(offset, count int) changes.Value {
	return s[s.FromUTF16(offset):s.FromUTF16(offset+count)]
}

// Count returns the number of UTF16 characters
func (s S16) Count() int {
	return s.ToUTF16(len(s))
}

// Apply implements Value.Apply
func (s S16) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return s
	case changes.Replace:
		if c.IsDelete {
			return changes.Nil
		}
		return c.After
	case changes.Splice:
		o := s.FromUTF16(c.Offset)
		remove := string(c.Before.(S16))
		replace := c.After.(S16)
		return s[:o] + replace + s[o+len(remove):]
	case changes.Move:
		ox, cx, dx := c.Offset, c.Count, c.Distance
		if c.Distance < 0 {
			ox, cx, dx = ox+dx, -dx, cx
		}
		x1, x2, x3 := s.FromUTF16(ox), s.FromUTF16(ox+cx), s.FromUTF16(ox+cx+dx)
		return s[:x1] + s[x2:x3] + s[x1:x2] + s[x3:]
	case changes.Custom:
		return c.ApplyTo(s)
	}
	panic("Unknown change type.  Cannot apply")
}

// FromUTF16 converts an UTF16 offset into a regular string offset
// It panics if the UTF16 offset lies within a UTF32 character
func (s S16) FromUTF16(idx int) int {
	seen := 0
	for kk, r := range s {
		if idx == seen {
			return kk
		}
		seen += s.utf16Count(r)
	}
	if idx == seen {
		return len(s)
	}
	panic("Unexpected idx value")
}

// ToUTF16 converts a regular index to UTF16 offsets
func (s S16) ToUTF16(idx int) int {
	seen := 0
	for kk, r := range s {
		if idx == kk {
			return seen
		}
		seen += s.utf16Count(r)
	}
	if idx == len(s) {
		return seen
	}
	panic("Unexpected idx value")
}

func (s S16) utf16Count(r rune) int {
	r1, r2 := utf16.EncodeRune(r)
	if r1 == '\uFFFD' && r2 == '\uFFFD' {
		return 1
	}
	return 2
}
