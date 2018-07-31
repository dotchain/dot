// +build !js,!tiny

// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding

import "reflect"

// RegisterConstructor registers the constructor of an encoding, typically
// called during init of an encoding package.
//
// Note that the constructor provided as argument should have a
// signature of the form:
//     newMyType(c Catalog, m map[string]interface{}) MyType
//
// where MyType implements ArrayLike, ObjectLike or UniversalEncoding.
//
// This functions uses the reflect package to deal with all the
// variations but to avoid the cost of pulling in the reflect package,
// this API is not available in the js builds.
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
