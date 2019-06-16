// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval

import (
	"errors"
	"strconv"
)

// the following are runtime errors
var errRecursion = errors.New("recursion detected")
var errUnknownField = errors.New("unknown field")
var errUnknownReceiver = errors.New("unknown receiver")
var errNotCallable = errors.New("not a function")
var errInvalidArgs = errors.New("invalid arguments")

// ParseError captures the error state
type ParseError struct {
	Offset  int
	Message string
}

// Error implements the Error interface
func (p ParseError) Error() string {
	return strconv.Itoa(p.Offset) + ": " + p.Message
}
