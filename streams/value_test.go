// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
	"reflect"
	"testing"
)

func TestValueStream(t *testing.T) {
	s := &streams.ValueStream{types.S8(""), streams.New()}
	splice := changes.Splice{0, types.S8(""), types.S8("Hello")}
	var cx changes.Change
	var sx streams.Stream = s
	s.Nextf("key", func() {
		if cx != nil {
			t.Fatal("Unexpected multiple call to Nextf")
		}
		sx, cx = sx.Next()
	})

	s2 := s.Append(splice)
	if s2.(*streams.ValueStream).Value != splice.After {
		t.Fatal("Append does not have the right value")
	}
	s3, c := s.Next()
	if c != splice || !reflect.DeepEqual(s3, s2) {
		t.Error("Next returned unexpected values", c, s3 == s2)
	}
	if cx != c || !reflect.DeepEqual(sx, s3) {
		t.Error("Nextf unexpected behavior")
	}

	if x, y := s3.Next(); x != nil || y != nil {
		t.Error("Next() didn't return nil", x, y)
	}

	s.Nextf("key", nil)
	s.ReverseAppend(splice)
}
