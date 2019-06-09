// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package rich implements rich text data types
package rich

import (
	"github.com/dotchain/dot/changes"
)

type valRun struct {
	changes.Value
	Size int
}

type values []valRun

func (a values) count() int {
	sum := 0
	for _, x := range a {
		sum += x.Size
	}
	return sum
}

func (a values) slice(offset, count int) values {
	seen := 0
	result := values{}
	for _, x := range a {
		start, end := seen, seen+x.Size
		if offset > start {
			start = offset
		}
		if offset+count < end {
			end = offset + count
		}
		if start < end {
			e := valRun{x.Value, end - start}
			result = append(result, e)
		}
		seen += x.Size
	}
	return result
}
