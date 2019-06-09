// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package html implements rich text to HTML conversion
package html

func must(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
