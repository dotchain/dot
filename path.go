// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

func commonPathLength(p1, p2 []string) int {
	if len(p2) < len(p1) {
		return commonPathLength(p2, p1)
	}

	for ii, ss := range p1 {
		if ss != p2[ii] {
			return ii
		}
	}
	return len(p1)
}
