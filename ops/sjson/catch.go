// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
package sjson

func catch(fn func()) (err error) {
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			if err, ok = r.(error); !ok {
				panic(r)
			}
		}
	}()
	fn()
	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
