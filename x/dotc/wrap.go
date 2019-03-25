// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"fmt"
	"unicode"
)

func wrapValue(v, vType string, atomic bool) string {
	if atomic || atomicTypes[vType] {
		return "changes.Atomic{" + v + "}"
	}
	if wrapper := wrappers[vType]; wrapper != "" {
		return fmt.Sprintf(wrapper, v)
	}

	return v
}

func unwrapValue(v, vType string, atomic bool) string {
	if atomic || atomicTypes[vType] {
		return v + ".(changes.Atomic).Value.(" + vType + ")"
	}

	if unwrapper := unwrappers[vType]; unwrapper != "" {
		return fmt.Sprintf(unwrapper, v)
	}

	return v + ".(" + vType + ")"
}

func streamType(s string) string {
	runes := []rune(s)
	for !unicode.IsLetter(runes[0]) {
		runes = runes[1:]
	}
	s = string(runes)
	if x, ok := streamTypes[s]; ok {
		return x
	}
	return s + "Stream"
}

var atomicTypes = map[string]bool{
	"bool": true,
	"int":  true,
}

var wrappers = map[string]string{
	"string": "types.S16(%s)",
}

var unwrappers = map[string]string{
	"string": "string(%s.(types.S16))",
}

var streamTypes = map[string]string{
	"bool":      "streams.Bool",
	"int":       "streams.Int",
	"string":    "streams.S16",
	"types.S16": "streams.S16",
	"types.S8":  "streams.S8",
}
