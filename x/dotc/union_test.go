// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc_test

import (
	"testing"

	"github.com/dotchain/dot/x/dotc"
	"github.com/tvastar/test"
)

func TestUnionGenerateApply(t *testing.T) {
	test.File(t.Error, "myunion/input.json", "myunion/generated.go", genUnion)
	test.File(t.Error, "myunion/input2.json", "myunion/generated2.go", genUnion)
	test.File(t.Error, "myunion/input.json", "myunion/generated_test.go", genUnionTests)
	test.File(t.Error, "myunion/input2.json", "myunion/generated2_test.go", genUnionTests)
}

func genUnion(s dotc.Union) (string, error) {
	info := dotc.Info{Package: "myunion", Unions: []dotc.Union{s}}
	info.UnionStreams = info.Unions
	code, err := info.Generate()
	if err != nil {
		logErrorContext(err, code)
	}
	return code, err
}

func genUnionTests(s dotc.Union) (string, error) {
	info := dotc.Info{Package: "myunion", Unions: []dotc.Union{s}}
	info.UnionStreams = info.Unions
	code, err := info.GenerateTests()
	if err != nil {
		logErrorContext(err, code)
	}
	return code, err
}
