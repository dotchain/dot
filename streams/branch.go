// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import "github.com/dotchain/dot/changes"

// Branch returns a new stream based on the provided stream. All
// changes made on the branch are only merged upstream when Push
// is called explicitly and all changes made upstream are only
// brought into the local branch when Pull is called explicitly.
func Branch(upstream Stream) Stream {
	downstream := New()
	b := &branchInfo{upstream, downstream, false}
	return branch{b, downstream}
}

type branch struct {
	info *branchInfo
	s    Stream
}

func (b branch) Push() error {
	return b.info.push()
}

func (b branch) Pull() error {
	return b.info.pull()
}

func (b branch) Undo() {
}

func (b branch) Redo() {
}

func (b branch) Append(c changes.Change) Stream {
	return branch{b.info, b.s.Append(c)}
}

func (b branch) ReverseAppend(c changes.Change) Stream {
	return branch{b.info, b.s.ReverseAppend(c)}
}

func (b branch) Next() (Stream, changes.Change) {
	s, c := b.s.Next()
	if s != nil {
		s = branch{b.info, s}
	}
	return s, c
}

type branchInfo struct {
	up, down Stream
	merging  bool
}

func (b *branchInfo) push() error {
	b.down, b.up = b.merge(b.down, b.up, false)
	return nil
}

func (b *branchInfo) pull() error {
	b.up, b.down = b.merge(b.up, b.down, true)
	return nil
}

func (b *branchInfo) merge(from, to Stream, reverse bool) (fromx, tox Stream) {
	if !b.merging {
		b.merging = true
		next, c := from.Next()
		for next != nil {
			if reverse {
				to = to.ReverseAppend(c)
			} else {
				to = to.Append(c)
			}
			from = next
			next, c = from.Next()
		}
		b.merging = false
	}
	return from, to

}
