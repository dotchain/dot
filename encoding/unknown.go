// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

type unknown struct {
	encoding string
	val      UniversalEncoding
}

func newUnknown(c Catalog, m map[string]interface{}) UniversalEncoding {
	name := m["dot:encoding"].(string)
	val, _ := c.TryGet(m["dot:encoded"])
	return unknown{name, val}
}

func (u unknown) NormalizeDOT() interface{} {
	return map[string]interface{}{
		"dot:encoding": u.encoding,
		"dot:generic":  true,
		"dot:encoded":  Normalize(u.val),
	}
}

func (u unknown) Count() int {
	return u.val.Count()
}

func (u unknown) Slice(offset, count int) ArrayLike {
	val := enrichArrayIfNeeded(u.val.Slice(offset, count))
	return unknown{u.encoding, val}
}

func (u unknown) Splice(offset int, before, after interface{}) ArrayLike {
	val := enrichArrayIfNeeded(u.val.Splice(offset, before, after))
	return unknown{u.encoding, val}
}

func (u unknown) RangeApply(offset, count int, fn func(interface{}) interface{}) ArrayLike {
	val := enrichArrayIfNeeded(u.val.RangeApply(offset, count, fn))
	return unknown{u.encoding, val}
}

func (u unknown) ForEach(fn func(offset int, val interface{})) {
	u.val.ForEach(fn)
}

func (u unknown) Get(key string) interface{} {
	return u.val.Get(key)
}

func (u unknown) Set(key string, value interface{}) ObjectLike {
	val := enrichObjectIfNeeded(u.val.Set(key, value))
	return unknown{u.encoding, val}
}

func (u unknown) ForKeys(fn func(key string, val interface{})) {
	u.val.ForKeys(fn)
}

func (u unknown) Contains(key string) bool {
	return u.val.Contains(key)
}

func (u unknown) IsArray() bool {
	return u.val.IsArray()
}
