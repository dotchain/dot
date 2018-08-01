// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import (
	"errors"
	"github.com/dotchain/dot/conv"
	"github.com/dotchain/dot/encoding"
)

// ErrPathInvalidated is returned by Ref.Update if the operations
// invalidate a particular path
var ErrPathInvalidated = errors.New("Path has been invalidated")

// ErrClientLogNeedsReconcile is returned by Ref.UpdateClient if the
// client log  used has not been reconciled with the provided log.
var ErrClientLogNeedsReconcile = errors.New("Client log needs reconciling")

// Ref holds a reference to some node in the virtual JSON tree
type Ref struct {
	ParentID, BasisID string
	Path              *RefPath
}

// Update transforms the provided ref path to its logical path in the
// server log provided
func (r Ref) Update(l *Log) (Ref, error) {
	id := "something not expected to be used"
	parents := []string{r.BasisID, r.ParentID}
	op := Operation{ID: id, Parents: parents}
	_, merge, err := l.TransformOperation(op)
	if err != nil {
		return r, err
	}
	path := r.Path
	for _, m := range merge {
		p, ok := path.Apply(m.Changes)
		if !ok {
			return r, ErrPathInvalidated
		}
		path = p
	}
	basisID := ""
	if len(l.Rebased) > 0 {
		basisID = l.Rebased[len(l.Rebased)-1].ID
	}
	return Ref{BasisID: basisID, Path: path}, nil
}

// UpdateClient attempts to map the provided ref into a path that
// can be used on the current client state.  It expects that the
// client log has been fully reconciled at this point and the ref has
// been updated against that server log with only local operations
// yet to  be factored in.
func (r Ref) UpdateClient(clog *ClientLog) (*RefPath, error) {
	path := r.Path
	for _, op := range clog.Rebased {
		p, ok := path.Apply(op.Changes)
		if !ok {
			return nil, ErrPathInvalidated
		}
		path = p
	}
	return path, nil
}

// RefIndexType defines whether the index is a pointer, cursor
// start or cursor end type.
type RefIndexType int

const (
	// RefIndexPointer tracks a specific element at that index
	RefIndexPointer RefIndexType = iota
	// RefIndexStart acts as if the index is a selection range start
	RefIndexStart
	// RefIndexEnd acts as if the index is a selection range end
	RefIndexEnd
)

// RefIndex tracks an array item in a ref path
type RefIndex struct {
	Index int
	Type  RefIndexType
}

// String converts an index into a string
func (r *RefIndex) String() string {
	suffix := ""
	switch r.Type {
	case RefIndexStart:
		suffix = "+"
	case RefIndexEnd:
		suffix = "-"
	}
	return conv.FromIndex(r.Index) + suffix
}

var m = map[string]RefIndexType{
	"+": RefIndexStart,
	"-": RefIndexEnd,
}

// NewRefIndex decodes a string version. It return nil if the string
// is not a validly formated index
func NewRefIndex(s string) *RefIndex {
	if len(s) > 0 {
		t := m[s[len(s)-1:]]
		if t != RefIndexPointer {
			s = s[:len(s)-1]
		}
		if conv.IsIndex(s) {
			return &RefIndex{Index: conv.ToIndex(s), Type: t}
		}
	}

	return nil
}

// RefPath represents a path within a virtual JSON object
type RefPath struct {
	key   string
	index *RefIndex
	next  *RefPath
}

// Prepend adds a new path entry before the provided path. Only one of
// key or index must be specified
func (r *RefPath) Prepend(key string, index *RefIndex) *RefPath {
	r.validate(key, index)
	return &RefPath{key: key, index: index, next: r}
}

// Append adds a new entry at the end of the current path but it does
// not modify the currrent path.  Instead it modifies a copy and
// returns that.
func (r *RefPath) Append(key string, index *RefIndex) *RefPath {
	r.validate(key, index)
	result := &RefPath{}
	last := result
	for r != nil {
		*last = *r
		last.next = &RefPath{}
		last, r = last.next, r.next
	}
	last.key, last.index = key, index
	return result
}

// Encode converts the path to an array of strings
func (r *RefPath) Encode() []string {
	result := []string{}
	for r != nil {
		if r.key != "" {
			result = append(result, r.key)
		} else {
			result = append(result, r.index.String())
		}
		r = r.next
	}
	return result
}

// NewRefPath creates a new ref path from the provided array of
// strings. An empty input is valid and returns a nil RefPath which is
// also valid and can be use against all the RefPath methods
func NewRefPath(s []string) *RefPath {
	var result *RefPath
	last := len(s) - 1
	for kk := range s {
		index := NewRefIndex(s[last-kk])
		key := ""
		if index == nil {
			key = s[last-kk]
		}

		result = result.Prepend(key, index)
	}
	return result
}

// Resolve attempts to walk the object for the specified path and
// return the value found.  It returns ok = true if it found the value
// successfully.
func (r *RefPath) Resolve(o interface{}) (interface{}, bool) {
	u := Utils(Transformer{})
	for r != nil {
		v, ok := u.C.TryGet(o)
		if !ok || v == nil {
			return nil, false
		}
		key := r.key
		if r.index != nil {
			key = conv.FromIndex(r.index.Index)
		}
		updated, ok := r.safeGet(v, key)
		if !ok {
			return nil, false
		}
		o, r = updated, r.next

	}

	return o, true
}

// Apply applies a set of changes and returns the effective new path.
// In case the path  was invalidated by the changes, it sets ok to false.
func (r *RefPath) Apply(changes []Change) (result *RefPath, ok bool) {
	result, ok = r, true
	for _, c := range changes {
		if result, ok = result.apply(c.Path, c); !ok {
			break
		}
	}
	return result, ok
}

func (r *RefPath) matches(s string) bool {
	if r.key == s || r.index == nil {
		return r.key == s
	}
	return conv.FromIndex(r.index.Index) == s
}

func (r *RefPath) apply(path []string, c Change) (result *RefPath, ok bool) {
	if r == nil {
		return nil, true
	}

	if len(path) > 0 && r.matches(path[0]) {
		result, ok := r.next.apply(path[1:], c)
		if !ok {
			return nil, false
		}

		if result != r.next {
			return result.Prepend(r.key, r.index), true
		}
		return r, true
	}

	if len(path) > 0 {
		return r, true
	}

	if r.key != "" {
		if c.Set == nil || c.Set.Key != r.key {
			return r, true
		}
		if _, ok = r.next.Resolve(c.Set.After); ok {
			return r, true
		}
		return nil, false
	}

	switch {
	case c.Splice != nil:
		return r.updateIndex(r.getSpliceIndex(c.Splice))
	case c.Move != nil:
		return r.updateIndex(r.getMoveIndex(c.Move))
	}
	return r.applyRange(c.Range)
}

func (r *RefPath) count(i interface{}) int {
	if i == nil {
		return 0
	}
	return Utils(Transformer{}).C.Get(i).Count()
}

func (r *RefPath) getSpliceIndex(s *SpliceInfo) int {
	offset, before, after := s.Offset, s.Before, s.After
	b, a, index := r.count(before), r.count(after), r.index.Index

	if r.index.Type == RefIndexPointer {
		if offset <= index && offset+b > index {
			return -1
		}
		if offset <= index {
			return index + a - b
		}
		return index
	}

	if r.index.Type == RefIndexStart {
		if offset+b <= index {
			return index + a - b
		}

		if offset >= index {
			return index
		}
		return offset + a
	}

	if offset >= index {
		return index
	}

	if offset+b <= index {
		return index + a - b
	}

	return offset
}

func (r *RefPath) getMoveIndex(m *MoveInfo) int {
	offset, count, distance := m.Offset, m.Count, m.Distance
	if distance < 0 {
		offset, count, distance = offset+distance, -distance, count
	}

	index := r.index.Index
	if index < offset || index >= offset+count+distance {
		return index
	}
	if index < offset+count {
		index += distance
	} else {
		index -= count
	}
	return index
}

func (r *RefPath) updateIndex(index int) (*RefPath, bool) {
	if index < 0 {
		return nil, false
	}
	if index == r.index.Index {
		return r, true
	}
	result, rindex := *r, *r.index
	result.index = &rindex
	result.index.Index = index
	return &result, true
}

func (r *RefPath) applyRange(ri *RangeInfo) (*RefPath, bool) {
	offset, count, index := ri.Offset, ri.Count, r.index.Index
	if offset > index || offset+count <= index {
		return r, true
	}

	updated, ok := r.next.Apply(ri.Changes)
	if !ok {
		return nil, false
	}
	if updated == r.next {
		return r, true
	}

	return updated.Prepend("", r.index), true
}

func (r *RefPath) validate(key string, index *RefIndex) {
	keyExists := key != ""
	indexExists := index != nil
	if keyExists == indexExists {
		panic("Specifying both key and index not allowed")
	}
}

func (r *RefPath) safeGet(u encoding.UniversalEncoding, key string) (result interface{}, ok bool) {
	// TODO: avoid this mess by modifying UniversalEncoding
	defer func() {
		if r := recover(); r != nil {
			result = nil
			ok = false
		}
	}()

	return u.Get(key), true
}
