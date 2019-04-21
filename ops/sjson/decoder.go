// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sjson

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strconv"
	"unicode"
)

// Decoder decodes any value
type Decoder struct {
	types map[string]reflect.Type
}

func (d *Decoder) register(typ reflect.Type) {
	if d.types == nil {
		d.types = map[string]reflect.Type{}
	}

	d.types[typeName(typ)] = typ
	switch typ.Kind() {
	case reflect.Slice:
		d.register(typ.Elem())
	case reflect.Ptr:
		d.register(typ.Elem())
	case reflect.Map:
		d.register(typ.Key())
		d.register(typ.Elem())
	case reflect.Struct:
		for idx := 0; idx < typ.NumField(); idx++ {
			field := typ.Field(idx)
			if field.PkgPath != "" {
				continue
			}
			d.register(field.Type)
		}
	}
}

// Decode implements ops/nw/Codec Decode method
func (d Decoder) Decode(value interface{}, r io.Reader) (err error) {
	return catch(func() {
		val := d.decode(bufio.NewReader(r))
		reflect.ValueOf(value).Elem().Set(val)
	})
}

func (d Decoder) decode(r *bufio.Reader) reflect.Value {
	if d.check("null", r) {
		return d.null()
	}

	if !d.check("{", r) {
		panic(errors.New("misssing {"))
	}

	key := d.readStringValue(r)
	if !d.check(":", r) {
		panic(errors.New("misssing :"))
	}

	result := d.decodeType(typeFromName(key, d.types), r)
	if !d.check("}", r) {
		panic(errors.New("missing }"))
	}
	return result
}

func (d Decoder) decodeType(typ reflect.Type, r *bufio.Reader) reflect.Value {
	switch typ.Kind() {
	case reflect.Bool:
		return d.decodeBoolValue(typ, r)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return d.decodeIntValue(typ, r)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return d.decodeUintValue(typ, r)
	case reflect.String:
		return d.decodeStringValue(typ, r)
	case reflect.Float32, reflect.Float64:
		return d.decodeFloatValue(typ, r)
	case reflect.Ptr:
		return d.decodePtr(typ, r)
	case reflect.Slice:
		return d.decodeSlice(typ, r)
	case reflect.Map:
		return d.decodeMap(typ, r)
	case reflect.Struct:
		return d.decodeStruct(typ, r)
	case reflect.Interface:
		return d.decode(r).Convert(typ)
	}

	panic(errors.New("unknown type " + typeName(typ)))
}

func (d Decoder) decodePtr(typ reflect.Type, r *bufio.Reader) reflect.Value {
	if d.check("null", r) {
		return reflect.Zero(typ)
	}
	inner := d.decodeType(typ.Elem(), r)
	result := reflect.New(typ.Elem())
	result.Elem().Set(inner)
	return result.Convert(typ)
}

func (d Decoder) decodeSlice(typ reflect.Type, r *bufio.Reader) reflect.Value {
	if d.check("null", r) {
		return reflect.Zero(typ)
	}
	if !d.check("[", r) {
		panic(errors.New("missing ["))
	}
	result := reflect.Zero(typ)
	finished := false
	for !finished {
		v := d.decodeType(typ.Elem(), r)
		result = reflect.Append(result, v)
		comma := d.check(",", r)
		finished = !comma && d.check("]", r)
		if !comma && !finished {
			panic(errors.New("missing , or ]"))
		}
	}
	return result
}

func (d Decoder) decodeMap(typ reflect.Type, r *bufio.Reader) reflect.Value {
	if d.check("null", r) {
		return reflect.Zero(typ)
	}
	if !d.check("[", r) {
		panic(errors.New("missing ["))
	}
	result := reflect.MakeMap(typ)
	finished := false
	for !finished {
		key := d.decodeType(typ.Key(), r)
		if !d.check(",", r) {
			panic(errors.New("missing ,"))
		}
		v := d.decodeType(typ.Elem(), r)
		result.SetMapIndex(key, v)
		comma := d.check(",", r)
		finished = !comma && d.check("]", r)
		if !comma && !finished {
			panic(errors.New("missing , or ]"))
		}
	}
	return result
}

func (d Decoder) decodeStruct(typ reflect.Type, r *bufio.Reader) reflect.Value {
	if !d.check("[", r) {
		panic(errors.New("missing ["))
	}
	result := reflect.New(typ).Elem()
	first := true
	for idx := 0; idx < typ.NumField(); idx++ {
		field := typ.Field(idx)
		if field.PkgPath != "" {
			continue
		}

		if !first && !d.check(",", r) {
			panic(errors.New("missing ,"))
		}
		first = false
		result.Field(idx).Set(d.decodeType(field.Type, r))
	}
	if !d.check("]", r) {
		panic(errors.New("missing }"))
	}
	return result
}

func (d Decoder) decodeBoolValue(typ reflect.Type, r *bufio.Reader) reflect.Value {
	if d.check("true", r) {
		return reflect.ValueOf(true).Convert(typ)
	}
	if d.check("false", r) {
		return reflect.ValueOf(false).Convert(typ)
	}
	panic(errors.New("unexpected bool value"))
}

func (d Decoder) decodeIntValue(typ reflect.Type, r *bufio.Reader) reflect.Value {
	return reflect.ValueOf(d.readIntValue(typ.Bits(), r)).Convert(typ)
}
func (d Decoder) decodeUintValue(typ reflect.Type, r *bufio.Reader) reflect.Value {
	return reflect.ValueOf(d.readUintValue(typ.Bits(), r)).Convert(typ)
}

func (d Decoder) decodeStringValue(typ reflect.Type, r *bufio.Reader) reflect.Value {
	return reflect.ValueOf(d.readStringValue(r)).Convert(typ)
}

func (d Decoder) decodeFloatValue(typ reflect.Type, r *bufio.Reader) reflect.Value {
	f, err := strconv.ParseFloat(d.readStringValue(r), typ.Bits())
	must(err)
	return reflect.ValueOf(f).Convert(typ)
}

func (d Decoder) null() reflect.Value {
	var result interface{}
	return reflect.Zero(reflect.ValueOf(&result).Elem().Type())
}

func (d Decoder) skipSpace(r *bufio.Reader) {
	for rn, _, err := r.ReadRune(); err == nil && unicode.IsSpace(rn); rn, _, err = r.ReadRune() {
	}

	must(r.UnreadRune())
}

func (d Decoder) check(s string, r *bufio.Reader) bool {
	d.skipSpace(r)

	if data, err := r.Peek(len(s)); err == nil && string(data) == s {
		_, err = r.Discard(len(s))
		must(err)
		return true
	}

	return false
}

func (d Decoder) readStringValue(r *bufio.Reader) string {
	var err error
	var escape bool
	var rn rune

	if !d.check(`"`, r) {
		panic(errors.New("missing \""))
	}

	runes := []rune{}
	for rn, _, err = r.ReadRune(); err == nil && (escape || rn != '"'); rn, _, err = r.ReadRune() {
		if escape || rn != '\\' {
			runes = append(runes, rn)
		}
		escape = !escape && rn == '\\'
	}
	must(err)
	if rn != '"' {
		panic(errors.New("unexpected char " + string([]rune{rn})))
	}

	return string(runes)
}

func (d Decoder) readIntValue(bitSize int, r *bufio.Reader) int64 {
	var b []rune
	var rn rune
	var err error

	d.skipSpace(r)

	for rn, _, err = r.ReadRune(); err == nil && (rn == '-' || unicode.IsNumber(rn)); rn, _, err = r.ReadRune() {
		b = append(b, rn)
	}
	must(err)
	must(r.UnreadRune())

	result, err := strconv.ParseInt(string(b), 10, bitSize)
	must(err)
	return result
}

func (d Decoder) readUintValue(bitSize int, r *bufio.Reader) uint64 {
	var b []rune
	var rn rune
	var err error

	d.skipSpace(r)

	for rn, _, err = r.ReadRune(); err == nil && (rn == '-' || unicode.IsNumber(rn)); rn, _, err = r.ReadRune() {
		b = append(b, rn)
	}
	must(err)
	must(r.UnreadRune())

	result, err := strconv.ParseUint(string(b), 10, bitSize)
	must(err)
	return result
}
