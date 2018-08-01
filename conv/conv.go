// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package conv implements converting array indices to/from string
//
// This package is trivial to implement using strconv but strconv
// pulls in other packages which increases the Javascript build of
// the DOT packge when using GopherJS.  So, these functions are
// implemented locally here
package conv

var numbers = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// FromIndex converts an array index to string
func FromIndex(idx int) string {
	if idx == 0 {
		return "0"
	}

	var val [30]rune
	size := 0
	for idx > 0 {
		val[size] = numbers[idx%10]
		idx = idx / 10
		size++
	}
	for kk := 0; kk < size/2; kk++ {
		val[kk], val[size-kk-1] = val[size-kk-1], val[kk]
	}
	return string(val[:size])
}

// IsIndex checks if a string is a valid index
func IsIndex(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// ToIndex converts a string into an array index
func ToIndex(s string) int {
	val := 0
	for len(s) > 0 {
		val = val*10 + int(s[0]-'0')
		s = s[1:]
	}
	return val
}
