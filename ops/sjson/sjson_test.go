// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sjson_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/dotchain/dot/ops/sjson"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

type myStruct struct {
	Boo        myInt32
	unexported float32
	Hoo        string
}

type myMap map[*[]int]int

func init() {
	sjson.Std.Register(myInt32(0))
	sjson.Std.Register([]stringer{nil})
	sjson.Std.Register(myStruct{})
	sjson.Std.Register(myMap{})
}

func TestSuccess(t *testing.T) {
	_ = myStruct{unexported: 52} // keep lint happy

	var i32 int32 = 5
	mi32 := myInt32(-22)
	var str stringer = mi32
	values := map[string]interface{}{
		// basic
		"{\"bool\": true}":                 true,
		"{\"bool\": false}":                false,
		"{\"uint\": 19}":                   uint(19),
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
		"{\"[]string\": [\"hello\",\"world\"]}": []string{"hello", "world"},
		"{\"[]string\": null}":                  []string(nil),

		// arrays
		`{"[2]int": [2,3]}`: [2]int{2, 3},

		// maps
		`{"map[string]string": null}`:              map[string]string(nil),
		`{"map[string]string": ["hello","world"]}`: map[string]string{"hello": "world"},

		// structs
		`{"ops/sjson_test.myStruct": [99,"balloons"]}`: myStruct{Hoo: "balloons", Boo: myInt32(99)},

		// map of interface => interface
		`{"map[ops/sjson_test.stringer]ops/sjson_test.stringer": [{"ops/sjson_test.myInt32": 42},{"ops/sjson_test.myInt32": 42}]}`: map[stringer]stringer{myInt32(42): myInt32(42)},

		// slices of interfaces
		`{"[]ops/sjson_test.stringer": [{"ops/sjson_test.myInt32": 42}]}`: []stringer{myInt32(42)},
		`{"[]ops/sjson_test.stringer": [null]}`:                           []stringer{nil},

		// ptr to interface
		`{"*ops/sjson_test.stringer": {"ops/sjson_test.myInt32": -22}}`: &str,

		// slices of pointers of named values
		`{"[]*ops/sjson_test.myInt32": [-22]}`: []*myInt32{&mi32},

		// type of map
		`{"ops/sjson_test.myMap": [[0],1]}`: myMap{&[]int{0}: 1},
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
			opt1 := cmpopts.IgnoreUnexported(myStruct{})
			opt2 := cmp.Comparer(func(v1, v2 myMap) bool {
				entries1 := []interface{}{}
				entries2 := []interface{}{}
				for k, v := range v1 {
					entries1 = append(entries1, *k, v)
				}
				for k, v := range v2 {
					entries2 = append(entries2, *k, v)
				}
				return cmp.Equal(entries1, entries2, opt1)
			})
			if !cmp.Equal(decoded, v, opt1, opt2) {
				t.Errorf("failed to decode %s %#v", expect, decoded)
			}
		})
	}
}

func TestMapMulti(t *testing.T) {
	v := map[int]int{1: 10, 2: 20}
	s := encode(v)
	if s != `{"map[int]int": [1,10,2,20]}` && s != `{"map[int]int": [2,20,1,10]}` {
		t.Fatal("Unexpected encoding", s)
	}

	if x := decode(s); !reflect.DeepEqual(x, v) {
		t.Error("Decoding multi error", x)
	}
}

func TestMapNested(t *testing.T) {
	p1 := []int{0}
	p2 := []int{1}
	v := map[*[]int]int{&p1: 5, &p2: 10}
	s := encode(v)
	expected := map[string]bool{
		`{"map[*[]int]int": [[0],5,[1],10]}`: true,
		`{"map[*[]int]int": [[1],10,[0],5]}`: true,
	}
	if !expected[s] {
		t.Fatal("Unexpected encoding", s)
	}

	if x := decode(s); !expected[encode(x)] {
		t.Errorf("Decoding nested map error %#v", x)
	}
}

func TestEncodeFailChannel(t *testing.T) {
	var encoded bytes.Buffer
	if err := sjson.Std.Encode(make(chan bool), &encoded); err == nil {
		t.Fatal("encoded channel", encoded.String())
	}
}

func TestEncodeFailWriter(t *testing.T) {
	if err := sjson.Std.Encode("", (failWriter{})); err == nil {
		t.Fatal("unexpected success with nil writer")
	}
}

type failWriter struct{}

func (w failWriter) Write(p []byte) (int, error) {
	return 0, errors.New("fail writes")
}

func TestDecodeMalformed(t *testing.T) {
	malformed := []string{
		`int`,
		`{int`,
		`{"int"}`,
		`{"int":`,
		`{"int":55`,
		`{"int":55,`,
		`{"int":"`,
		`{"[]int":55}`,
		`{"[]int":[55`,
		`{"[]int":[55?`,
		`{"[2]int": 55}`,
		`{"[2]int": [55]}`,
		`{"[2]int": [55,10}`,
		`{"map[int]int":55}`,
		`{"map[int]int":[`,
		`{"map[int]int":[0]}`,
		`{"map[int]int":[0,]}`,
		`{"map[int]int":[0,6?]}`,
		`{"map[int]int":[0,6,7]}`,
		`{"map[": 22}`,
		`{"ops/sjson_test.myStruct": {}}`,
		`{"ops/sjson_test.myStruct": 55}`,
		`{"ops/sjson_test.myStruct": []}`,
		`{"ops/sjson_test.myStruct": ["hello"]}`,
		`{"ops/sjson_test.myStruct": [42]}`,
		`{"ops/sjson_test.myStruct": [42, "helloo", 2]}`,
		`{"boo":"`,
		`{"bool":5}`,
		`{"string": "hello`,
		`{"func": ""}`,
	}

	sjson.Std.Register(func() {})
	for _, test := range malformed {
		t.Run(test, func(t *testing.T) {
			var result interface{}
			err := sjson.Std.Decode(&result, bytes.NewReader([]byte(test)))
			if err == nil {
				t.Error("Failed to detect  malformed value", result)
			}
		})
	}
}
