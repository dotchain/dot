// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc_test

import (
	"testing"

	"github.com/dotchain/dot/x/dotc"
	"github.com/tvastar/test"
)

func TestStructStreamGenerateInterface(t *testing.T) {
	test.File(t.Error, "mystruct/input3.json", "mystruct/generated3.go", genStructStream)
}

func genStructStream(s dotc.Struct) (string, error) {
	info := dotc.Info{Package: "mystruct", Structs: []dotc.Struct{s}}
	code, err := info.Generate()
	if err != nil {
		logErrorContext(err, code)
	}
	return code, err
}
