// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// templates generates golang code from simple text templates.
package main

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main() {
	data := map[string]string{}
	for _, arg := range os.Args[1:] {
		pair := strings.Split(arg, "=")
		if len(pair) == 2 {
			data[pair[0]] = pair[1]
		} else {
			data["template"] = arg
		}
	}

	text, err := ioutil.ReadFile(data["template"])
	if err != nil {
		panic(err)
	}

	t, err := template.New(data["template"]).Parse(string(text))
	if err != nil {
		panic(err)
	}

	out, err := os.Create(data["out"])
	if err != nil {
		panic(err)
	}

	if err := t.Execute(out, data); err != nil {
		panic(err)
	}
}
