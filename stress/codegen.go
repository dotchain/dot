// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build ignore

package main

import (
	"io/ioutil"

	"github.com/dotchain/dot/x/dotc"
)

func main() {
	info.StructStreams = info.Structs
	code, err := info.Generate()
	if err != nil {
		panic(err)
	}
	code = "//+build stress\n\n" + code
	err = ioutil.WriteFile("generated.go", []byte(code), 0644)
	if err != nil {
		panic(err)
	}
}

var info = dotc.Info{
	Package: "stress",
	Structs: []dotc.Struct{{
		Recv: "s",
		Type: "State",
		Fields: []dotc.Field{{
			Name:               "Text",
			Key:                "text",
			Type:               "string",
			ToValueFmt:         "types.S8(%s)%.s",
			FromValueFmt:       "string(%s.(types.S8))%.s",
			FromStreamValueFmt: "%s%.s",
			ToStreamFmt:        "streams.S8%.s",
		}, {
			Name:               "Count",
			Key:                "count",
			Type:               "types.Counter",
			ToValueFmt:         "%s%.s",
			FromValueFmt:       "%s.(types.Counter)%.s",
			FromStreamValueFmt: "int32(%s)%.s",
			ToStreamFmt:        "streams.Counter%.s",
		}},
	}},
	StructStreams: []dotc.Struct{{
		Recv: "s",
		Type: "State",
		Fields: []dotc.Field{{
			Name:               "Text",
			Key:                "text",
			Type:               "string",
			ToValueFmt:         "types.S8(%s)%.s",
			FromValueFmt:       "string(%s.(types.S8))%.s",
			FromStreamValueFmt: "%s%.s",
			ToStreamFmt:        "streams.S8%.s",
		}, {
			Name:               "Count",
			Key:                "count",
			Type:               "types.Counter",
			ToValueFmt:         "%s%.s",
			FromValueFmt:       "%s.(types.Counter)%.s",
			FromStreamValueFmt: "int32(%s)%.s",
			ToStreamFmt:        "streams.Counter%.s",
		}},
	}},
}
