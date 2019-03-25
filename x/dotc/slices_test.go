// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc_test

import (
	"testing"

	"github.com/dotchain/dot/x/dotc"
	"github.com/tvastar/test"
)

func TestSliceStreamTests(t *testing.T) {
	test.File(t.Error, "myslice/input.json", "myslice/generated_test.go", genSlicesStreamTests)
}

func genSlicesStreamTests(s []dotc.Slice) (string, error) {
	info := dotc.Info{Package: "myslice", Slices: s}
	code, err := info.GenerateTests()
	if err != nil {
		logErrorContext(err, code)
	}
	return code, err
}
