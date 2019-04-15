// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc_test

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/dotchain/dot/x/dotc"
	"github.com/tvastar/test"
)

func TestStructGenerateApply(t *testing.T) {
	test.File(t.Error, "mystruct/input.json", "mystruct/generated.go", genStruct)
	test.File(t.Error, "mystruct/input2.json", "mystruct/generated2.go", genStruct)
	test.File(t.Error, "mystruct/input3.json", "mystruct/generated3.go", genStruct)
	test.File(t.Error, "mystruct/input.json", "mystruct/generated_test.go", genStructTests)
	test.File(t.Error, "mystruct/input2.json", "mystruct/generated2_test.go", genStructTests)
	test.File(t.Error, "mystruct/input3.json", "mystruct/generated3_test.go", genStructTests)
}

func genStruct(s dotc.Struct) (string, error) {
	info := dotc.Info{Package: "mystruct", Structs: []dotc.Struct{s}}
	info.StructStreams = info.Structs
	code, err := info.Generate()
	if err != nil {
		logErrorContext(err, code)
	}
	return code, err
}

func genStructTests(s dotc.Struct) (string, error) {
	info := dotc.Info{Package: "mystruct", Structs: []dotc.Struct{s}}
	info.StructStreams = info.Structs
	code, err := info.GenerateTests()
	if err != nil {
		logErrorContext(err, code)
	}
	return code, err
}

func logErrorContext(e error, code string) {
	re := regexp.MustCompile("[0-9]+:[0-9]+:")
	found := re.FindString(e.Error())
	if found == "" {
		return
	}
	line, err := strconv.Atoi(strings.Split(found, ":")[0])
	if err != nil {
		return
	}
	lines := strings.Split(code, "\n")
	before, after := line-5, line+1
	if before < 0 {
		before = 0
	}
	if after > len(lines) {
		after = len(lines)
	}
	if before >= after {
		log.Println("\nerror:", e)
		return
	}
	log.Println("error", strings.Join(lines[before:after], "\n"), "\nerror:", e)
}
