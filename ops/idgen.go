// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package ops

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID returns a unique ID using crypto/rand
func NewID() interface{} {
	var b [32]byte
	_, err := rand.Read(b[0:32])
	must(err)
	return hex.EncodeToString(b[0:32])
}
