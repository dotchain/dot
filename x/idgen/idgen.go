// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package idgen generates a new unique ID.
//
// It uses crypto/rand to generate 32 random bytes encoded as 64 hex
// characters.
package idgen

import (
	"crypto/rand"
	"encoding/hex"
)

// New returns a unique ID using crypto/rand
func New() interface{} {
	var b [32]byte
	if _, err := rand.Read(b[0:32]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[0:32])
}
