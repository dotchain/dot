// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package types implements OT-compatible immutable values.
//
// Strings are implemented by S8 and S16.   S8 implements the standard
// Go string where offsets refer to bytes.  S16 is a better choice for
// working with UFT16 encoded values such as is native  to
// Javascript.  In this case, the offsets count the number of UF16
// units (which maps to native JS string offsets).
//
// General arrays values can be represented by A while M implements
// maps.
//
// Counter implements a 32-bit integer counter. This also serves as an
// example of an interesting data structure as it uses a virtual array
// as far as OT is concerned but only stores the accumuated count.
//
// A much richer type is available at
// https://godoc.org/github.com/dotchain/dot/x/rt which also
// demonstrates how to implement a custom change that is applicable
// to the type.
package types
