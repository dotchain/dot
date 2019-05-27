// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

type S = types.S8

func TestStream(t *testing.T) {
	s := streams.New()

	s1 := s.Append(changes.Replace{Before: changes.Nil, After: S("Hello World")})

	c1_1 := s1.Append(changes.Splice{Offset: 0, Before: S(""), After: S("A ")})
	c1_2 := c1_1.Append(changes.Splice{Offset: 2, Before: S(""), After: S("B ")})

	c2_1 := s1.Append(changes.Splice{Offset: 0, Before: S(""), After: S("X ")})
	c2_1_merged, _ := streams.Latest(s)

	c2_2 := c2_1.Append(changes.Splice{Offset: 2, Before: S(""), After: S("Y ")})
	c2_2_with_c1_1, _ := c2_2.Next()

	c1_3 := c2_1_merged.Append(changes.Splice{Offset: 6, Before: S(""), After: S("C ")})
	c2_3 := c2_2_with_c1_1.Append(changes.Splice{Offset: 6, Before: S(""), After: S("Z ")})

	_, c := streams.Latest(s)
	if v := S("").Apply(nil, c); v != S("A B X Y C Z Hello World") {
		t.Error("Merge failed: ", v)
		t.Error("changes", c1_1, c1_2, c1_3)
		t.Error("changes", c2_1, c2_2, c2_3)
	}
}

func TestStreamPushPullUndoRedo(t *testing.T) {
	s := streams.New()
	s.Undo()
	s.Redo()
	if err := s.Push(); err != nil {
		t.Error("push", err)
	}
	if err := s.Pull(); err != nil {
		t.Error("pull", err)
	}
}
