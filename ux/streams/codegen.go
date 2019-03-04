// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build ignore

package main

import (
	"github.com/dotchain/dot/compiler"
	"io/ioutil"
)

func main() {
	info := compiler.Info{
		Package: "streams",
		Imports: [][2]string{},
		Streams: []compiler.StreamInfo{
			{
				StreamType: "BoolStream",
				ValueType:  "bool",
			},
			{
				StreamType: "TextStream",
				ValueType:  "string",
			},
		},
	}
	ioutil.WriteFile("generated.go", []byte(compiler.Generate(info)), 0644)
}
