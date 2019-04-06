// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package testops

import (
	"github.com/dotchain/dot/ops"
)

// NullCache implements ops.Cache but does not save anything
func NullCache() ops.Cache {
	return nullCache{}
}

// MemCache implements ops.Cache
func MemCache() ops.Cache {
	return &memCache{map[int]ops.Op{}, map[int][]ops.Op{}}
}

type nullCache struct{}

func (nc nullCache) Load(ver int) (ops.Op, []ops.Op) {
	return nil, nil
}

func (nc nullCache) Store(key int, op ops.Op, merge []ops.Op) {
}

type memCache struct {
	x     map[int]ops.Op
	merge map[int][]ops.Op
}

func (mc *memCache) Load(ver int) (ops.Op, []ops.Op) {
	return mc.x[ver], mc.merge[ver]
}

func (mc *memCache) Store(ver int, op ops.Op, merge []ops.Op) {
	if _, ok := mc.x[ver]; ok {
		panic("recalculaed version")
	}
	mc.x[ver] = op
	mc.merge[ver] = merge
}
