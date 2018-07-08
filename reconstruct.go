// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

import "github.com/dotchain/dot/encoding"

// Reconstruct creates a minimal sequence of changes that will
// effectively produce the same model as the provided input
func (u Utils) Reconstruct(model interface{}) []Change {
	if model == nil {
		return nil
	}
	if encoding.IsString(model) {
		return []Change{{Splice: &SpliceInfo{After: model}}}
	}

	x, ok := u.C.TryGet(model)
	if !ok || x == nil {
		return nil
	}

	if x.IsArray() {
		return []Change{{Splice: &SpliceInfo{After: x}}}
	}
	changes := []Change{}
	x.ForKeys(func(key string, val interface{}) {
		set := &SetInfo{Key: key, After: val}
		changes = append(changes, Change{Set: set})
	})
	return changes

}
