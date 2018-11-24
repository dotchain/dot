// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Splice represents an array edit change.  A set of elements from the
// specified offset are removed and replaced with a new set of elements.
type Splice struct {
	Offset        int
	Before, After Collection
}

// Revert inverts the effect of the splice.
func (s Splice) Revert() Change {
	return Splice{s.Offset, s.After, s.Before}
}

// MergeReplace merges a move against a Replace.  The replace always wins
func (s Splice) MergeReplace(other Replace) (other1 *Replace, s1 *Splice) {
	other.Before = other.Before.Apply(nil, s)
	return &other, nil
}

// MergeSplice merges a splice against another Splice.
func (s Splice) MergeSplice(other Splice) (other1, s1 *Splice) {
	sstart, send := s.Offset, s.Offset+s.Before.Count()
	ostart, oend := other.Offset, other.Offset+other.Before.Count()

	switch {
	case send <= ostart: // [  ] <  >
		other.Offset += s.After.Count() - (send - sstart)
		return &other, &s
	case sstart >= oend: // < > [ ]
		sx := s
		sx.Offset += other.After.Count() - (oend - ostart)
		return &other, &sx
	case sstart < ostart && send < oend: // [  < ]  >
		other.Offset = sstart + s.After.Count()
		other.Before = other.Before.Slice(send-ostart, oend-send)

		s.Before = s.Before.Slice(0, ostart-sstart)
		return &other, &s

	case sstart == ostart && send < oend: // [<  ]   >
		other.Before = other.Before.ApplyCollection(nil, Splice{0, s.Before, s.After})
		return &other, nil

	case sstart <= ostart && send >= oend: // [ < > ]
		sliced := s.Before.Slice(ostart-sstart, oend-ostart)
		s.Before = s.Before.ApplyCollection(nil, Splice{ostart - sstart, sliced, other.After})
		return nil, &s

	case sstart > ostart && send <= oend: // < [ ]>
		sliced := other.Before.Slice(sstart-ostart, send-sstart)
		other.Before = other.Before.ApplyCollection(nil, Splice{sstart - ostart, sliced, s.After})
		return &other, nil
	default: // sstart < oend: // < [ > ]
		other.Before = other.Before.Slice(0, sstart-ostart)
		sx := s
		sx.Offset = ostart + other.After.Count()
		sx.Before = s.Before.Slice(oend-sstart, send-oend)
		return &other, &sx
	}
}

// MergeMove merges a splice against a move
func (s Splice) MergeMove(o Move) (ox, sx Change) {
	beforeSize := s.Before.Count()
	if o.Offset >= s.Offset && o.Offset+o.Count <= s.Offset+beforeSize {
		return s.mergeContainedMove(o)
	}

	if o.Offset <= s.Offset && o.Offset+o.Count >= s.Offset+beforeSize {
		o.Count += s.After.Count() - beforeSize
		s.Offset += o.Distance
		return o, s
	}

	if o.Offset >= s.Offset+beforeSize || s.Offset >= o.Offset+o.Count {
		return s.mergeNonOverlappingMove(o)
	}

	// first undo the intersection and then merge as before
	rest, undo := o, Move{Offset: o.Offset + o.Distance}
	if o.Offset > s.Offset {
		left := s.Offset + beforeSize - o.Offset
		rest.Offset += left
		rest.Count -= left
		undo.Count = left
		if o.Distance < 0 {
			rest.Distance -= left
			undo.Distance = o.Count - left - o.Distance
		} else {
			undo.Distance = -o.Distance
		}
	} else {
		right := o.Offset + o.Count - s.Offset
		rest.Count -= right
		undo.Count = right
		undo.Offset += rest.Count
		if o.Distance < 0 {
			undo.Distance = -o.Distance
		} else {
			rest.Distance += right
			undo.Distance = right - o.Distance - o.Count
		}
	}

	ox, sx = s.mergeNonOverlappingMove(rest)
	if a, ok := sx.(ChangeSet); ok {
		sx = ChangeSet(append([]Change{undo}, a...))
	} else {
		sx = ChangeSet{undo, sx}
	}
	return ox, sx
}

func (s Splice) mergeNonOverlappingMove(o Move) (ox, sx Change) {
	odest, beforeSize := o.dest(), s.Before.Count()
	diff := s.After.Count() - beforeSize
	if odest > s.Offset && odest < s.Offset+beforeSize {
		right := s.Offset + beforeSize - odest
		s1 := s
		s1.Before = s.Before.Slice(0, odest-s.Offset)
		empty := s.Before.Slice(0, 0)
		s2 := Splice{o.Offset + o.Count + o.Distance, empty, empty}
		s2.Before = s.Before.Slice(odest-s.Offset, right)

		if o.Offset < s.Offset {
			s1.Offset -= o.Count
			o.Distance += right + diff
		} else {
			o.Distance += right
			o.Offset += diff
		}
		return o, ChangeSet([]Change{s2, s1})
	}

	if odest <= s.Offset {
		if o.Offset > s.Offset {
			o.Offset += diff
			o.Distance -= diff
			s.Offset += o.Count
		}
	} else if odest >= s.Offset+s.Before.Count() {
		if o.Offset > s.Offset {
			o.Offset += diff
		} else {
			o.Distance += diff
			s.Offset -= o.Count
		}
	}
	return o, s
}

func (s Splice) mergeContainedMove(o Move) (ox, sx Change) {
	beforeSize, odest := s.Before.Count(), o.dest()
	if odest >= s.Offset && odest <= s.Offset+beforeSize {
		s.Before = s.Before.ApplyCollection(nil, Move{o.Offset - s.Offset, o.Count, o.Distance})
		return nil, s
	}

	sliced := s.Before.Slice(o.Offset-s.Offset, o.Count)
	empty := sliced.Slice(0, 0)
	spliced := s.Before.ApplyCollection(nil, Splice{o.Offset - s.Offset, sliced, empty})

	if odest < s.Offset {
		ox = Splice{odest, empty, sliced}
		sx = Splice{s.Offset + o.Count, spliced, s.After}
		return ox, sx
	}
	ox = Splice{odest + s.After.Count() - beforeSize, empty, sliced}
	s.Before = spliced
	return ox, s
}

// Merge implements the Change.Merge method
func (s Splice) Merge(other Change) (otherx, cx Change) {
	if other == nil {
		return nil, s
	}

	switch o := other.(type) {
	case Replace:
		return change(s.MergeReplace(o))
	case Splice:
		return change(s.MergeSplice(o))
	case Move:
		return s.MergeMove(o)
	case Custom:
		return swap(o.ReverseMerge(s))
	}
	panic("Unexpected change")
}

// MapIndex maps an index to the new location of the index after the
// splice. It also returns whether the item at that index has been
// modified by the splice change.
func (s Splice) MapIndex(idx int) (int, bool) {
	switch {
	case idx < s.Offset:
		return idx, false
	case idx >= s.Offset+s.Before.Count():
		return idx + s.After.Count() - s.Before.Count(), false
	}
	return s.Offset, true
}

// Change returns nil or the underlying Splice
func (s *Splice) Change() Change {
	if s == nil {
		return nil
	}
	return *s
}
