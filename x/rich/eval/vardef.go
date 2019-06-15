// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich/data"
)

// vardef must only appear with do() or object() blocks
// and is a fake "changes.Value" that is not meant to be persisted
type vardef struct {
	key   *data.Ref
	value changes.Value
}

func (v vardef) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return nil
}
