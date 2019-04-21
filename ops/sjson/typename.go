// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sjson

import (
	"errors"
	"reflect"
	"strings"
)

func typeName(v reflect.Type) string {
	if path, name := v.PkgPath(), v.Name(); path != "" && name != "" {
		path = strings.TrimPrefix(path, "github.com/dotchain/dot/")
		path = strings.TrimPrefix(path, "github.com/")
		return path + "." + name
	}

	switch v.Kind() {
	case reflect.Ptr:
		return "*" + typeName(v.Elem())
	case reflect.Slice:
		return "[]" + typeName(v.Elem())
	case reflect.Map:
		return "map[" + typeName(v.Key()) + "]" + typeName(v.Elem())
	}
	return v.Kind().String()
}

func typeFromName(name string, types map[string]reflect.Type) reflect.Type {
	if v, ok := types[name]; ok {
		return v
	}

	if v, ok := typesDefault[name]; ok {
		return v
	}

	if strings.HasPrefix(name, "*") {
		return reflect.PtrTo(typeFromName(name[1:], types))
	}

	if strings.HasPrefix(name, "[]") {
		return reflect.SliceOf(typeFromName(name[2:], types))
	}

	if strings.HasPrefix(name, "map") {
		return typeFromMap(name, types)
	}

	panic(errors.New("unknown type " + name))
}

func typeFromMap(name string, types map[string]reflect.Type) reflect.Type {
	count := 0
	for idx, rn := range name[4:] {
		switch rn {
		case '[':
			count++
		case ']':
			count--
		}
		if count < 0 {
			keyType := typeFromName(name[4:4+idx], types)
			elemType := typeFromName(name[idx+5:], types)
			return reflect.MapOf(keyType, elemType)
		}
	}
	panic(errors.New("invalid type name: " + name))
}

var typesDefault = map[string]reflect.Type{
	"bool":    reflect.TypeOf(false),
	"int8":    reflect.TypeOf(int8(0)),
	"int16":   reflect.TypeOf(int16(0)),
	"int32":   reflect.TypeOf(int32(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"uint8":   reflect.TypeOf(uint8(0)),
	"uint16":  reflect.TypeOf(uint16(0)),
	"uint32":  reflect.TypeOf(uint32(0)),
	"uint64":  reflect.TypeOf(uint64(0)),
	"string":  reflect.TypeOf(""),
	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),
}
