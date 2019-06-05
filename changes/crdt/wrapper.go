// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt

import "github.com/dotchain/dot/changes"

type crdtChange interface {
	ApplyTo(ctx changes.Context, v changes.Value) changes.Value
	Revert() crdtChange
}

type wrapper []crdtChange

func (w wrapper) Merge(o changes.Change) (otherx, cx changes.Change) {
	return o, w
}

func (w wrapper) ReverseMerge(o changes.Change) (otherx, cx changes.Change) {
	return o, w
}

func (w wrapper) Revert() changes.Change {
	result := make(wrapper, len(w))
	for kk, cx := range w {
		result[len(w)-kk-1] = cx.Revert()
	}
	return result
}

func (w wrapper) ApplyTo(ctx changes.Context, v changes.Value) changes.Value {
	for _, cx := range w {
		v = cx.ApplyTo(ctx, v)
	}
	return v
}
