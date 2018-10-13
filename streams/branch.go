// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams

import "github.com/dotchain/dot/changes"

// Branch returns a new stream based on the provided stream. All
// changes made on the branch are only merged upstream when an
// an explicit call Push and all changes made upstream are only
// brought into the local branch on an explicit call to Pull.
func Branch(upstream Stream) Stream {
	downstream := New()
	b := &branchInfo{upstream, downstream, false}
	return branch{b, downstream}
}

// Push can be called on any branch to push all local changes
// upstream.  It will panic if called on a non-branch
func Push(s Stream) {
	s.(branch).info.push()
}

// Pull can be called on any branch to pull all upstream changes over
// locally. It will panic if called on a non-branch
func Pull(s Stream) {
	s.(branch).info.pull()
}

// Connect connects two streams as if downstream was created by a call
// to Branch with the caveat that connected streams do not support
// calls to Push or Pull -- they are effectively always in merge
// mode.
func Connect(up, down Stream) {
	b := &branchInfo{up, down, false}
	up.Nextf(b, b.pull)
	down.Nextf(b, b.push)
}

type branch struct {
	info *branchInfo
	s    Stream
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

func (b branch) Nextf(key interface{}, fn func()) {
	b.s.Nextf(key, fn)
}

type branchInfo struct {
	up, down Stream
	merging  bool
}

func (b *branchInfo) push() {
	b.down, b.up = b.merge(b.down, b.up, false)
}

func (b *branchInfo) pull() {
	b.up, b.down = b.merge(b.up, b.down, true)
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
