// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import (
	"github.com/dotchain/dot"
	"sync"
)

// root implements a root container
type root struct {
	// the value is only stored at the last item.
	v interface{}

	// the change that spawned this version
	rebased, compensation []dot.Change

	// the branch where the change was made
	branch *basis

	// the basis of this new version
	own basis

	// the next change. the mutex protects just this value
	// as the rest are immutable
	next *root
	sync.Mutex
}

// Bubble is called by the value with the basis the change was made on
// as well as the basis of the change itself.
func (r *root) Bubble(prev, now *basis, changes []dot.Change) {
	// find the container which has the right basis
	for prev != r.branch && prev != &r.own {
		r = r.next
	}

	compensation := []dot.Change(nil)
	if prev == r.branch {
		// bootstrap with prior compensation if the change was
		// on top of the unmerged basis. if the prev == &r.own
		// that implies all the current changes have been
		// folded in, so no prior compensation to deal with
		l := len(r.compensation)
		compensation = r.compensation[:l:l]
	}

	r.Lock()
	for r.next != nil {
		r.Unlock()
		r = r.next
		compensation = append(compensation, r.compensation...)
		r.Lock()
	}
	defer r.Unlock()

	changes, compensation = x.MergeChanges(compensation, changes)
	v := utils.Apply(r.v, changes)
	next := &root{v: v, rebased: changes, compensation: compensation, branch: now}

	// we clear the r.v so that older values are garbage collected
	r.v, r.next = nil, next
}

// Latest maps the path to any changes (using the basis to detect the
// changes) and returns the latest value, the version with the latest
// value and the new path + basis
func (r *root) Latest(path *dot.RefPath, b *basis) (interface{}, parent, []string, *basis) {
	for b != r.branch && b != &r.own {
		r = r.next
	}

	changes := []dot.Change(nil)
	if b == r.branch {
		// the path is based off the branch, so start with the
		// compensation
		l := len(r.compensation)
		changes = r.compensation[:l:l]
	}

	r.Lock()
	for r.next != nil {
		r.Unlock()
		changes = append(changes, r.rebased...)
		r = r.next
		r.Lock()
	}
	v := r.v
	r.Unlock()
	if path, ok := path.Apply(changes); ok {
		return v, r, path.Encode(), &r.own
	}
	return nil, nil, nil, nil
}
