// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

import (
	"encoding/json"
	"reflect"
	"sync"
)

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
	types map[reflect.Type]constructor
	sync.Mutex
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
		types: map[reflect.Type]constructor{},
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
	case UniversalEncoding, String16:
		b, err := json.Marshal(i)
		if err != nil {
			panic(err)
		}
		var result interface{}
		if err = json.Unmarshal(b, &result); err != nil {
			panic(err)
		}
		return result
	default:
		return i
	}
}

func (c Catalog) getConstructor(name string) constructor {
	if c.catalog == nil {
		return Default.getConstructor(name)
	}

	c.Lock()
	defer c.Unlock()
	ctor := c.names[name]
	if ctor != nil || c.catalog == Default.catalog {
		return ctor
	}
	Default.Lock()
	defer Default.Unlock()
	return Default.names[name]
}

// RegisterTypeConstructor associates a name (such as "dot:utf16")
// with a type and a constructor (which returns that type).
//
// The return type from the constructor must implement either
// ObjectLike or ArrayLike (or both).  The returned
// value from the constructor should also properly deal with JSON
// formatting by implementing MarshalJSON as defined in encoding/json.
func (c Catalog) RegisterTypeConstructor(name string, t reflect.Type, fn func(Catalog, map[string]interface{}) UniversalEncoding) {
	c.Lock()
	defer c.Unlock()

	c.names[name] = fn
	c.types[t] = fn
}

// RegisterConstructor registers the constructor of an encoding, typically
// called during init of an encoding package
func (c Catalog) RegisterConstructor(name string, fn interface{}) {
	fType := reflect.TypeOf(fn)
	if fType.Kind() != reflect.Func {
		panic(errNotFunction)
	}

	if fType.NumIn() != 2 {
		panic(errNumArgs)
	}

	if fType.In(0) != reflect.TypeOf(c) {
		panic(errFirstArgMustBeCatalog)
	}

	var dummy map[string]interface{}
	if fType.In(1) != reflect.TypeOf(dummy) {
		panic(errSecondArgMustBeMap)
	}

	if fType.NumOut() != 1 {
		panic(errSingleReturnValue)
	}

	resultType := fType.Out(0)
	zero := reflect.Zero(resultType).Interface()
	if _, ok := zero.(ArrayLike); ok {
		ctor := func(cat Catalog, m map[string]interface{}) UniversalEncoding {
			args := []reflect.Value{reflect.ValueOf(cat), reflect.ValueOf(m)}
			val := reflect.ValueOf(fn).Call(args)[0].Interface()
			return enrichArrayIfNeeded(val.(ArrayLike))
		}
		c.RegisterTypeConstructor(name, resultType, ctor)
		return
	}

	if _, ok := zero.(ObjectLike); ok {
		ctor := func(cat Catalog, m map[string]interface{}) UniversalEncoding {
			args := []reflect.Value{reflect.ValueOf(cat), reflect.ValueOf(m)}
			val := reflect.ValueOf(fn).Call(args)[0].Interface()
			return enrichObjectIfNeeded(val.(ObjectLike))
		}
		c.RegisterTypeConstructor(name, resultType, ctor)
		return
	}

	panic(errUnexpectedType)
}
