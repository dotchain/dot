// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package changes

// Move represents a shuffling of some elements (specified  by Offset
// and Count) over to a different spot (specified by Distance, which
// can be negative to indicate a move over to the left).
type Move struct {
	Offset, Count, Distance int
}

// Revert reverts the move.
func (m Move) Revert() Change {
	return Move{m.Offset + m.Distance, m.Count, -m.Distance}
}

// MergeReplace merges a move against a Replace.  The replace always wins
func (m Move) MergeReplace(other Replace) (other1 *Replace, m1 *Splice) {
	other.Before = other.Before.Apply(nil, m)
	return &other, nil
}

// MergeSplice merges a splice with a move.
func (m Move) MergeSplice(o Splice) (Change, Change) {
	x, y := o.MergeMove(m)
	return y, x
}

// MergeMove merges a move against another Move
func (m Move) MergeMove(o Move) (ox []Move, mx []Move) {
	if m == o {
		return nil, nil
	}

	if m.Distance == 0 || m.Count == 0 || o.Distance == 0 || o.Count == 0 {
		return []Move{o}, []Move{m}
	}

	if m.Offset >= o.Offset+o.Count || o.Offset >= m.Offset+m.Count {
		return m.mergeMoveNoOverlap(o)
	}

	if m.Offset <= o.Offset && m.Offset+m.Count >= o.Offset+o.Count {
		return m.mergeMoveContained(o)
	}

	if m.Offset >= o.Offset && m.Offset+m.Count <= o.Offset+o.Count {
		return m.swap(o.mergeMoveContained(m))
	}

	if m.Offset < o.Offset {
		return m.mergeMoveRightOverlap(o)
	}

	return m.swap(o.mergeMoveRightOverlap(m))
}

func (m Move) mergeMoveNoOverlap(o Move) (ox, mx []Move) {
	mdest, odest := m.dest(), o.dest()

	if !m.contains(odest) && !o.contains(mdest) {
		return m.mergeMoveNoOverlapNoDestMixups(o)
	}

	if m.contains(odest) && o.contains(mdest) {
		return m.mergeMoveNoOverlapMixedDests(o)
	}

	if o.contains(mdest) {
		return m.swap(o.mergeMoveNoOverlap(m))
	}

	return m.mergeMoveNoOverlapContainedDest(o)
}

func (m Move) mergeMoveNoOverlapContainedDest(o Move) (ox, mx []Move) {
	mdest, odest := m.dest(), o.dest()

	mdestNew := mdest
	if mdest >= odest && mdest <= o.Offset {
		mdestNew += o.Count
	} else if mdest > o.Offset && mdest <= odest {
		mdestNew -= o.Count
	}

	m1 := m
	if o.Offset <= m.Offset {
		m1.Offset -= o.Count
	}
	m1.Count = m.Count + o.Count
	if mdestNew <= m1.Offset {
		m1.Distance = mdestNew - m1.Offset
	} else {
		m1.Distance = mdestNew - m1.Offset - m1.Count
	}

	if o.Offset > m.Offset && o.Offset < mdest {
		o.Offset -= m.Count
	} else if o.Offset >= mdest && o.Offset < m.Offset {
		o.Offset += m.Count
	}
	odest += m.Distance
	if odest <= o.Offset {
		o.Distance = odest - o.Offset
	} else {
		o.Distance = odest - o.Offset - o.Count
	}

	return []Move{o}, []Move{m1}
}

func (m Move) mergeMoveNoOverlapNoDestMixups(o Move) (ox, mx []Move) {
	mdest, odest := m.dest(), o.dest()
	o1 := Move{Offset: m.mapPoint(o.Offset), Count: o.Count}
	o1 = o1.withDest(m.mapPoint(odest))
	if odest == mdest {
		o1 = o1.withDest(m.Offset + m.Distance)
	}

	m1 := Move{Offset: o.mapPoint(m.Offset), Count: m.Count}
	m1 = m1.withDest(o.mapPoint(mdest))

	return []Move{o1}, []Move{m1}
}

func (m Move) mergeMoveNoOverlapMixedDests(o Move) (ox, mx []Move) {
	var oleft, oright Move
	mdest, odest := m.dest(), o.dest()

	oleft.Count = mdest - o.Offset
	oright.Count = o.Count - oleft.Count
	oleft.Offset = m.Offset + m.Distance - oleft.Count
	oright.Offset = m.Offset + m.Distance + m.Count
	oleft.Distance = odest - m.Offset
	oright.Distance = odest - m.Offset - m.Count
	ox = []Move{oleft, oright}

	distance := o.Offset - m.Offset - m.Count
	if distance < 0 {
		distance = -(m.Offset - o.Offset - o.Count)
	}
	mx = []Move{{
		Offset:   o.Offset + o.Distance - (odest - m.Offset),
		Count:    m.Count + o.Count,
		Distance: distance,
	}}
	return ox, mx
}

func (m Move) mergeMoveRightOverlap(o Move) ([]Move, []Move) {
	overlapSize := m.Offset + m.Count - o.Offset
	overlapUndo := Move{Offset: o.Offset + o.Distance, Count: overlapSize}
	non := Move{Offset: o.Offset + overlapSize, Count: o.Count - overlapSize}

	if o.Distance > 0 {
		overlapUndo.Distance = -o.Distance
		non.Distance = o.Distance
	} else {
		overlapUndo.Distance = o.Count - overlapSize - o.Distance
		non.Distance = o.Distance - overlapSize
	}
	l, r := m.mergeMoveNoOverlap(non)
	return l, append([]Move{overlapUndo}, r...)
}

func (m Move) mergeMoveContained(o Move) ([]Move, []Move) {
	odest := o.dest()
	ox := o
	ox.Offset += m.Distance
	if m.Offset <= odest && odest <= m.Offset+m.Count {
		return []Move{ox}, []Move{m}
	}

	if odest == m.dest() {
		ox = ox.withDest(m.Offset + m.Distance)
		mx := m
		mx.Count -= o.Count
		if o.Distance < 0 {
			mx.Offset += o.Count
		}
		mx = mx.withDest(o.Offset + o.Count + o.Distance)
		return []Move{ox}, []Move{mx}
	}

	ox = ox.withDest(m.mapPoint(odest))

	mx := m
	mx.Offset = o.mapPoint(m.Offset)
	mx.Count = m.Count - o.Count
	mx = mx.withDest(o.mapPoint(m.dest()))

	return []Move{ox}, []Move{mx}
}

// MapIndex maps a particular index to the new location of the index
// after the move
func (m Move) MapIndex(idx int) int {
	switch {
	case idx >= m.Offset+m.Distance && idx < m.Offset:
		return idx + m.Count
	case idx >= m.Offset && idx < m.Offset+m.Count:
		return idx + m.Distance
	case idx >= m.Offset+m.Count && idx < m.Offset+m.Count+m.Distance:
		return idx - m.Count
	}
	return idx
}

func (m Move) mapPoint(idx int) int {
	switch {
	case idx >= m.Offset+m.Distance && idx <= m.Offset:
		return idx + m.Count
	case idx >= m.Offset+m.Count && idx < m.Offset+m.Count+m.Distance:
		return idx - m.Count
	}
	return idx
}

func (m Move) dest() int {
	if m.Distance < 0 {
		return m.Offset + m.Distance
	}
	return m.Offset + m.Distance + m.Count
}

func (m Move) withDest(dest int) Move {
	m.Distance = dest - m.Offset - m.Count
	if m.Distance < 0 {
		m.Distance = dest - m.Offset
	}
	return m
}

func (m Move) contains(idx int) bool {
	return idx > m.Offset && idx < m.Offset+m.Count
}

func (m Move) swap(l, r []Move) ([]Move, []Move) {
	return r, l
}

func movesToChange(m []Move) Change {
	result := make([]Change, len(m))
	for _, mm := range m {
		if mm.Count != 0 && mm.Distance != 0 {
			result = append(result, mm)
		}
	}
	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	}
	return ChangeSet(result)
}

// Merge implements the Change.Merge method
func (m Move) Merge(other Change) (otherx, cx Change) {
	if other == nil {
		return nil, m
	}

	switch o := other.(type) {
	case Replace:
		return change(m.MergeReplace(o))
	case Splice:
		return m.MergeSplice(o)
	case Move:
		l, r := m.MergeMove(o)
		return movesToChange(l), movesToChange(r)
	case Custom:
		return swap(o.ReverseMerge(m))
	}
	panic("Unexpected change")
}

// Change returns either nil or the underlying Move as a change.
func (m *Move) Change() Change {
	if m == nil {
		return nil
	}
	return *m
}

// Normalize ensures that distance is always positive
func (m Move) Normalize() Move {
	if m.Distance < 0 {
		return Move{m.Offset + m.Distance, -m.Distance, m.Count}
	}
	return m
}
