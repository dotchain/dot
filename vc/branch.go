// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc

import (
	"github.com/dotchain/dot"
	"sync"
)

// Branch provides the methods to synchronize a specific branch
type Branch interface {
	// Push takes any that have happened on the branch and applies
	// it to the parent branch. The current branch itself is not
	// updated as such, so calls to Latest on the current branch
	// will not reflect any changes on the main branch
	Push()
}

type branch struct {
	sync.Mutex
	// control on the parent branch where the branch started
	parent Control
	// child container
	*root
}

func (b *branch) Push() {
	b.Lock()
	defer b.Unlock()

	r := b.root
	changes := []dot.Change(nil)

	r.Lock()
	for r.next != nil {
		r.Unlock()
		r = r.next
		changes = append(changes, r.rebased...)
		r.Lock()
	}
	r.Unlock()
	b.root = r

	if len(changes) != 0 {
		b.parent.UpdateSync(changes)
	}
}
