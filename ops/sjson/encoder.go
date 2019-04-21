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
)

// Encoder encodes any value
type Encoder struct {
}

// Encode implements ops/nw/Codec Encode method
func (e Encoder) Encode(value interface{}, w io.Writer) error {
	return catch(func() {
		buf := bufio.NewWriter(w)
		e.encode(reflect.ValueOf(value), buf)
		must(buf.Flush())
	})
}

func (e Encoder) encode(v reflect.Value, w *bufio.Writer) {
	if v.IsValid() {
		_, err := w.WriteString("{\"" + typeName(v.Type()) + "\": ")
		must(err)
	}

	e.encodeValue(v, w)
	if v.IsValid() {
		_, err := w.WriteString("}")
		must(err)
	}
}

func (e Encoder) encodeValue(v reflect.Value, w *bufio.Writer) {
	if !v.IsValid() {
		_, err := w.Write([]byte("null"))
		must(err)
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		e.encodeBoolValue(v.Bool(), w)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.encodeIntValue(v.Int(), w)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.encodeUintValue(v.Uint(), w)
	case reflect.Float32:
		e.encodeFloatValue(v.Float(), 32, w)
	case reflect.Float64:
		e.encodeFloatValue(v.Float(), 64, w)
	case reflect.String:
		e.encodeStringValue(v.String(), w)
	case reflect.Ptr:
		e.encodePtrValue(v, w)
	case reflect.Slice:
		e.encodeSliceValue(v, w)
	case reflect.Map:
		e.encodeMapValue(v, w)
	case reflect.Interface:
		e.encode(reflect.ValueOf(v.Interface()), w)
	default:
		panic(errors.New("not yet implemented: " + v.Kind().String()))
	}
}

func (e Encoder) encodePtrValue(v reflect.Value, w *bufio.Writer) {
	if v.IsNil() {
		e.encodeValue(reflect.ValueOf(nil), w)
	} else {
		e.encodeValue(v.Elem(), w)
	}
}

func (e Encoder) encodeSliceValue(v reflect.Value, w *bufio.Writer) {
	if v.IsNil() {
		e.encodeValue(reflect.ValueOf(nil), w)
	} else {
		_, err := w.WriteString("[")
		must(err)
		for idx := 0; idx < v.Len(); idx++ {
			if idx > 0 {
				_, err = w.WriteString(",")
				must(err)
			}

			e.encodeValue(v.Index(idx), w)
		}
		_, err = w.WriteString("]")
		must(err)
	}
}

func (e Encoder) encodeMapValue(v reflect.Value, w *bufio.Writer) {
	if v.IsNil() {
		e.encodeValue(reflect.ValueOf(nil), w)
	} else {
		_, err := w.WriteString("[")
		must(err)
		keys := v.MapKeys()
		for idx := range keys {
			if idx > 0 {
				_, err = w.WriteString(",")
				must(err)
			}
			e.encodeValue(keys[idx], w)
			_, err = w.WriteString(",")
			must(err)
			e.encodeValue(v.MapIndex(keys[idx]), w)
		}
		_, err = w.WriteString("]")
		must(err)
	}
}

func (e Encoder) encodeBoolValue(b bool, w *bufio.Writer) {
	_, err := w.WriteString(strconv.FormatBool(b))
	must(err)
}

func (e Encoder) encodeIntValue(i int64, w *bufio.Writer) {
	_, err := w.WriteString(strconv.FormatInt(i, 10))
	must(err)
}

func (e Encoder) encodeUintValue(ui uint64, w *bufio.Writer) {
	_, err := w.WriteString(strconv.FormatUint(ui, 10))
	must(err)
}

func (e Encoder) encodeFloatValue(f float64, bitSize int, w *bufio.Writer) {
	_, err := w.WriteString(`"` + strconv.FormatFloat(f, 'g', -1, bitSize) + `"`)
	must(err)
}

func (e Encoder) encodeStringValue(s string, w *bufio.Writer) {
	_, err := w.WriteRune('"')
	must(err)

	for _, rn := range s {
		if rn == '"' {
			_, err = w.WriteRune('\\')
			must(err)
		}
		_, err = w.WriteRune(rn)
		must(err)
	}
	_, err = w.WriteRune('"')
	must(err)
}
