// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package encoding allows different physical representations to provide
// the logical JSON representation needed by DOT.
//
// For example, a sparse array may be represented as a dictionary
// but the SparseArray encoding allows it to provide a regular
// array like interface to DOT.
//
// To identify when an encoding is used, the actual representation
// of an encoding is a map with two keys: "dot:encoding" and
// "dot:encoded" -- the former provides the name of the encoding
// and the later provides the encoded value
//
// Encodings can be registered with Encodings.Register
package encoding

import (
	"strconv"
	"unicode/utf16"
)

// ArrayLike is the default interface to be implemented by
// encodings of collections
type ArrayLike interface {
	Count() int
	Slice(offset, count int) ArrayLike
	Splice(offset int, before, after interface{}) ArrayLike
	RangeApply(offset, count int, fn func(interface{}) interface{}) ArrayLike
	ForEach(func(offset int, val interface{}))
}

// ObjectLike is the default interface to be implemented by
// encodings that behave like objects
type ObjectLike interface {
	Get(key string) interface{}
	Set(key string, value interface{}) ObjectLike
	ForKeys(func(key string, val interface{}))
}

// UniversalEncoding is a combination of both ArrayLike and ObjectLike
type UniversalEncoding interface {
	ArrayLike
	ObjectLike
	IsArray() bool
}

type enrichArray struct{ ArrayLike }

func (e enrichArray) Get(key string) interface{} {
	i, err := strconv.Atoi(key)
	if err != nil {
		panic(errArrayKeyIsNotNumber)
	}

	var result interface{}
	e.Slice(i, 1).ForEach(func(_ int, v interface{}) {
		result = v
	})
	return result
}

func (e enrichArray) Set(key string, value interface{}) ObjectLike {
	i, err := strconv.Atoi(key)
	if err != nil {
		panic(errArrayKeyIsNotNumber)
	}

	r := e.RangeApply(i, 1, func(_ interface{}) interface{} {
		return value
	})
	return enrichArrayIfNeeded(r)
}

func (e enrichArray) ForKeys(fn func(key string, val interface{})) {
	e.ArrayLike.ForEach(func(offset int, val interface{}) {
		fn(strconv.Itoa(offset), val)
	})
}

func (e enrichArray) IsArray() bool {
	return true
}

type enrichObject struct{ ObjectLike }

func (enrichObject) Count() int {
	panic(errMethodNotSupported)
}

func (enrichObject) Slice(offset, count int) ArrayLike {
	panic(errMethodNotSupported)
}

func (enrichObject) Splice(offset int, before, after interface{}) ArrayLike {
	panic(errMethodNotSupported)
}

func (enrichObject) RangeApply(offset, count int, fn func(interface{}) interface{}) ArrayLike {
	panic(errMethodNotSupported)
}

func (enrichObject) ForEach(func(offset int, val interface{})) {
	panic(errMethodNotSupported)
}

func (o enrichObject) IsArray() bool {
	return false
}

func enrichArrayIfNeeded(e ArrayLike) UniversalEncoding {
	if u, ok := e.(UniversalEncoding); ok {
		return u
	}
	return enrichArray{e}
}

func enrichObjectIfNeeded(o ObjectLike) UniversalEncoding {
	if u, ok := o.(UniversalEncoding); ok {
		return u
	}
	return enrichObject{o}
}

// IsString identifies if the interface is nil or a valid string
func IsString(i interface{}) bool {
	if i == nil {
		return true
	}
	switch i := i.(type) {
	case string:
		return true
	case String16:
		return true
	case enrichArray:
		return IsString(i.ArrayLike)
	}
	return false
}

// ToString converts a potentially encoded string into an actual string
func ToString(i interface{}) string {
	if i == nil {
		return ""
	}
	switch i := i.(type) {
	case string:
		return i
	case String16:
		return string(utf16.Decode(i))
	case enrichArray:
		return ToString(i.ArrayLike)
	}
	return ""
}
