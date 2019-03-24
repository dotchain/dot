// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import "fmt"

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
