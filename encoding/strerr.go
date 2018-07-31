// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

type strerr string

func (s strerr) Error() string {
	return string(s)
}

var errArrayKeyIsNotNumber = strerr("array key is not a number")
var errMethodNotSupported = strerr("method is not supported")
var errIndexOutOfBounds = strerr("array index is out of bounds")
var errStringIndexOutOfBounds = strerr("string index is out of bounds")
var errUnknownEncoding = strerr("could not find a json-like encoding")
var errNotFunction = strerr("arg is not a function")
var errNumArgs = strerr("arg is not a 2-param function")
var errFirstArgMustBeCatalog = strerr("function should use first arg as catalog")
var errSecondArgMustBeMap = strerr("function second arg must be map[string]interface{}")
var errSingleReturnValue = strerr("function must return a single value")
var errUnexpectedType = strerr("unexpected return value type")
var errNoSuchField = strerr("did not find any field for the specified key")
var errNotStringType = strerr("found a non-string type")
