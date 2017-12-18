// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package encoding_test

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("String16", String16{}.TestAll)
	t.Run("Array", ArrayTest().TestAll)
	t.Run("Sparse", SparseTest().TestAll)
	t.Run("Dict", DictTest().TestAll)
	t.Run("Set", SetTest().TestAll)
}

func shouldPanic(t *testing.T, message string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Failed to panic", message)
		}
	}()
	f()
}

func marshal(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func unmarshal(s string) interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		panic(err)
	}
	return result
}

func ensureEqual(t *testing.T, s1, s2 interface{}) {
	m1, m2 := marshal(s1), marshal(s2)
	u1, u2 := unmarshal(m1), unmarshal(m2)
	if !reflect.DeepEqual(u1, u2) {
		t.Error("Expected equal values but got", m1, m2)
	}
}
