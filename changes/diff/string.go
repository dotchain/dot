// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package diff

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// S8 computes the diff between two types.S9 values
func S8(d Differ, old, new changes.Value) changes.Change {
	o := string(old.(types.S8))
	n := string(new.(types.S8))

	result := changes.ChangeSet(nil)
	offset := 0
	for _, diff := range diffmatchpatch.New().DiffMain(o, n, false) {
		splice := changes.Splice{Offset: offset, Before: types.S8(""), After: types.S8("")}
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			splice.Before = types.S8(diff.Text)
		case diffmatchpatch.DiffInsert:
			splice.After = types.S8(diff.Text)
			offset += splice.After.Count()
		default:
			offset += types.S8(diff.Text).Count()
			continue
		}
		result = append(result, splice)
	}

	if result == nil {
		return nil
	}
	return result
}

// S16 computes the diff between two types.S9 values
func S16(d Differ, old, new changes.Value) changes.Change {
	o := string(old.(types.S16))
	n := string(new.(types.S16))

	result := changes.ChangeSet(nil)
	offset := 0
	for _, diff := range diffmatchpatch.New().DiffMain(o, n, false) {
		splice := changes.Splice{Offset: offset, Before: types.S16(""), After: types.S16("")}
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			splice.Before = types.S16(diff.Text)
		case diffmatchpatch.DiffInsert:
			splice.After = types.S16(diff.Text)
			offset += splice.After.Count()
		default:
			offset += types.S16(diff.Text).Count()
			continue
		}
		result = append(result, splice)
	}

	if result == nil {
		return nil
	}
	return result
}
