// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sjson_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dotchain/dot/ops/sjson"
)

func encode(v interface{}) string {
	var encoded bytes.Buffer
	if err := sjson.Std.Encode(v, &encoded); err != nil {
		panic(err)
	}
	return encoded.String()
}

func decode(s string) interface{} {
	var result interface{}
	if err := sjson.Std.Decode(&result, bytes.NewReader([]byte(s))); err != nil {
		panic(err)
	}
	return result
}

func TestNull(t *testing.T) {
	if x := encode(nil); x != "null" {
		t.Error("failed to encode null", x)
	}

	if x := decode("null"); x != nil {
		t.Error("failed to decode null", x)
	}

	if x := decode("\t    null\n\n    "); x != nil {
		t.Error("failed to decode null", x)
	}
}

type zint32 int32
type myInt32 zint32

func (m myInt32) String() string {
	return "boo"
}

type stringer interface {
	String() string
}

func init() {
	sjson.Std.Register(myInt32(0))
	sjson.Std.Register([]stringer{nil})
}

func TestCases(t *testing.T) {
	var i32 int32 = 5
	mi32 := myInt32(-22)
	var str stringer = mi32
	values := map[string]interface{}{
		// basic
		"{\"bool\": true}":                 true,
		"{\"bool\": false}":                false,
		"{\"uint8\": 5}":                   byte(5),
		"{\"uint16\": 256}":                uint16(256),
		"{\"uint32\": 70000}":              uint32(70000),
		"{\"uint64\": 0}":                  uint64(0),
		"{\"int8\": -3}":                   int8(-3),
		"{\"int16\": -256}":                int16(-256),
		"{\"int32\": -70000}":              int32(-70000),
		"{\"int64\": 9}":                   int64(9),
		"{\"string\": \"hello\\\"world\"}": "hello\"world",
		"{\"string\": \"\"}":               "",
		"{\"float32\": \"-3.1\"}":          float32(-3.1),
		"{\"float64\": \"-2.22\"}":         float64(-2.22),

		// pointers
		"{\"*int32\": 5}":     &i32,
		"{\"*int32\": null}":  (*int32)(nil),
		"{\"**uint8\": null}": (**uint8)(nil),

		// named  types
		"{\"ops/sjson_test.myInt32\": 22}":   myInt32(22),
		"{\"*ops/sjson_test.myInt32\": -22}": &mi32,

		// slices
		"{\"[]string\": [\"hello\"]}": []string{"hello"},
		"{\"[]string\": null}":        []string(nil),

		// maps
		`{"map[string]string": null}`:              map[string]string(nil),
		`{"map[string]string": ["hello","world"]}`: map[string]string{"hello": "world"},

		// map of interface => interface
		`{"map[ops/sjson_test.stringer]ops/sjson_test.stringer": [{"ops/sjson_test.myInt32": 42},{"ops/sjson_test.myInt32": 42}]}`: map[stringer]stringer{myInt32(42): myInt32(42)},

		// slices of interfaces
		`{"[]ops/sjson_test.stringer": [{"ops/sjson_test.myInt32": 42}]}`: []stringer{myInt32(42)},

		// ptr to interface
		`{"*ops/sjson_test.stringer": {"ops/sjson_test.myInt32": -22}}`: &str,

		// slices of pointers of named values
		`{"[]*ops/sjson_test.myInt32": [-22]}`: []*myInt32{&mi32},
	}

	for expect, v := range values {
		data, _ := json.Marshal(v)
		jsonv := string(data)
		t.Run(jsonv, func(t *testing.T) {
			got := encode(v)
			if got != expect {
				t.Fatal("failed to encode", got)
			}
			decoded := decode(got)
			if !reflect.DeepEqual(decoded, v) {
				t.Errorf("failed to decode %v", reflect.TypeOf(int(5)))
			}
		})
	}
}
