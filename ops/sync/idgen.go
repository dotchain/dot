// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package sync

import (
	"crypto/rand"
	"encoding/hex"
)

// newID returns a unique ID using crypto/rand
func (s *session) newID() (interface{}, error) {
	var b [32]byte
	_, err := rand.Read(b[:])
	return hex.EncodeToString(b[:]), err
}
