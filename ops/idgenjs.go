// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build js,!jsreflect

package ops

import "github.com/gopherjs/gopherjs/js"

// NewID returns a unique ID using crypto
func NewID() interface{} {
	crypto := js.Global.Get("crypto")
	array := js.Global.Get("Uint8Array").New(32)
	crypto.Call("getRandomValues", array)
	result := js.Global.Get("btoa", string(array.Interface().([]byte)))
	return result.Interface()
}
