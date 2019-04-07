// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc_test

import (
	"testing"

	"github.com/dotchain/dot/x/dotc"
)

func TestInvalidNames(t *testing.T) {
	info := dotc.Info{Package: "invalid package"}
	if _, err := info.Generate(); err == nil {
		t.Error("Unexpected success with invalid package name", err)
	}

	if _, err := info.GenerateTests(); err == nil {
		t.Error("Unexpected success with invalid package name", err)
	}
}
