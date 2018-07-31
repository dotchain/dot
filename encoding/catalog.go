// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

import "unicode/utf16"

// Default is the default catalog that encodings register with.
var Default = NewCatalog()

// Get takes an arbitrary interface and returns a UniversalEncoding
// type that DOT can work with
func Get(i interface{}) UniversalEncoding {
	return Default.Get(i)
}

// Unget does the reverse of Get
func Unget(i interface{}) interface{} {
	return Default.Unget(i)
}

type constructor func(Catalog, map[string]interface{}) UniversalEncoding

type catalog struct {
	names map[string]constructor
	types map[interface{}]constructor
}

// Catalog is a copyable thread-safe collection of encodings.
//
// Encodings can be easily registered via #Catalog.RegisterConstructor.
//
// Catalogs inherit from the #Default catalog.
type Catalog struct {
	*catalog
}

// NewCatalog creates a new catalog.
func NewCatalog() Catalog {
	c := &catalog{
		names: map[string]constructor{},
		types: map[interface{}]constructor{},
	}
	return Catalog{catalog: c}
}

// TryGet attempts to convert the provide interface to a UniversalEncoding
// if it is either an array-like or object-like type.  If not, it sets ok
// to false to indicate it could not find a good type
func (c Catalog) TryGet(i interface{}) (UniversalEncoding, bool) {
	if i == nil {
		return nil, true
	}

	switch i := i.(type) {
	case UniversalEncoding:
		return i, true
	case ArrayLike:
		return enrichArray{i}, true
	case ObjectLike:
		return enrichObject{i}, true
	case string:
		return enrichArray{NewString16(i)}, true
	case []interface{}:
		return enrichArray{NewArray(c, i)}, true
	case map[string]interface{}:
		if name, ok := i["dot:encoding"].(string); ok {
			if ctor := c.getConstructor(name); ctor != nil {
				return ctor(c, i), true
			}
		}
		return enrichObject{Dict(i)}, true
	default:
		return nil, false
	}
}

// Get takes an arbitrary interface and returns a UniversalEncoding
// type that DOT can work with
func (c Catalog) Get(i interface{}) UniversalEncoding {
	if result, ok := c.TryGet(i); ok {
		return result
	}
	panic(errUnknownEncoding)
}

// Unget reverses any wrapping done by Get but does not do this recursively.
func (c Catalog) Unget(i interface{}) interface{} {
	if i == nil {
		return nil
	}

	switch i := i.(type) {
	case enrichArray:
		return c.Unget(i.ArrayLike)
	case Array:
		return c.Unget(i.v)
	case enrichObject:
		return c.Unget(i.ObjectLike)
	case Dict:
		return map[string]interface{}(i)
	case String16:
		return string(utf16.Decode(i))
	default:
		return i
	}
}

func (c Catalog) getConstructor(name string) constructor {
	if c.catalog == nil {
		return Default.getConstructor(name)
	}

	ctor := c.names[name]
	if ctor != nil || c.catalog == Default.catalog {
		return ctor
	}
	return Default.names[name]
}

// RegisterTypeConstructor associates a name (such as "dot:utf16")
// with a type and a constructor (which returns that type).
//
// To avaid importing the reflect package unnecessarily, the
// returnType arg is weakly typed.
func (c Catalog) RegisterTypeConstructor(name string, returnType interface{}, fn func(Catalog, map[string]interface{}) UniversalEncoding) {
	c.names[name] = fn
	c.types[returnType] = fn
}

// RegisterArrayConstructor is a minor variation on
// RegisterTypeConstrutor where the constructor function returns an
// ArrayLike type instead of UniversalEncoding
func (c Catalog) RegisterArrayConstructor(name string, returnType interface{}, fn func(Catalog, map[string]interface{}) ArrayLike) {
	f := func(c Catalog, m map[string]interface{}) UniversalEncoding {
		return enrichArrayIfNeeded(fn(c, m))
	}
	c.RegisterTypeConstructor(name, returnType, f)
}

// RegisterObjectConstructor is a minor variation on
// RegisterTypeConstrutor where the constructor function returns an
// ObjectLike type instead of UniversalEncoding
func (c Catalog) RegisterObjectConstructor(name string, returnType interface{}, fn func(c Catalog, m map[string]interface{}) ObjectLike) {
	f := func(c Catalog, m map[string]interface{}) UniversalEncoding {
		return enrichObjectIfNeeded(fn(c, m))
	}
	c.RegisterTypeConstructor(name, returnType, f)
}
