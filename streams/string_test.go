// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package streams_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/streams"
)

func TestS16Stream(t *testing.T) {
	s := streams.New()
	strong := &streams.S16{Stream: s, Value: "10"}
	if strong.Stream != s {
		t.Fatal("Unexpected Stream()", strong.Stream)
	}

	strong = strong.Update("15")
	if strong.Value != "15" {
		t.Error("Update did not change value", strong.Value)
	}
	s, c := s.Next()

	before, after := types.S16("10"), types.S16("15")
	if !reflect.DeepEqual(c, changes.Replace{Before: before, After: after}) {
		t.Error("Unexpected change on main stream", c)
	}

	c = changes.Splice{Offset: 1, Before: types.S16(""), After: types.S16("2")}
	s = s.Append(c)
	c = changes.Move{Offset: 1, Count: 1, Distance: 1}
	s = s.Append(c)
	strong = strong.Latest()

	if strong.Value != "152" {
		t.Error("Unexpected change on string stream", strong.Value)
	}

	if _, c := strong.Next(); c != nil {
		t.Error("Unexpected change on string stream", c)
	}

	strong = strong.Splice(1, 1, "9") // now 192
	strong = strong.Move(0, 1, 1)     // now 912
	if strong.Value != "912" {
		t.Error("Unexpected strong.Value", strong.Value)
	}

	s, _ = streams.Latest(s)
	s = s.Append(changes.Replace{Before: types.S16("152"), After: changes.Nil})

	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on string stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: types.S16("99")})
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on string stream", c, strong)
	}

}

func TestS8Stream(t *testing.T) {
	s := streams.New()
	strong := &streams.S8{Stream: s, Value: "10"}
	if strong.Stream != s {
		t.Fatal("Unexpected Stream()", strong.Stream)
	}

	strong = strong.Update("15")
	if strong.Value != "15" {
		t.Error("Update did not change value", strong.Value)
	}
	s, c := s.Next()

	before, after := types.S8("10"), types.S8("15")
	if !reflect.DeepEqual(c, changes.Replace{Before: before, After: after}) {
		t.Error("Unexpected change on main stream", c)
	}

	c = changes.Splice{Offset: 1, Before: types.S8(""), After: types.S8("2")}
	s = s.Append(c)
	c = changes.Move{Offset: 1, Count: 1, Distance: 1}
	s = s.Append(c)
	strong = strong.Latest()

	if strong.Value != "152" {
		t.Error("Unexpected change on string stream", strong.Value)
	}

	if _, c := strong.Next(); c != nil {
		t.Error("Unexpected change on string stream", c)
	}

	strong = strong.Splice(1, 1, "9") // now 192
	strong = strong.Move(0, 1, 1)     // now 912
	if strong.Value != "912" {
		t.Error("Unexpected strong.Value", strong.Value)
	}

	s, _ = streams.Latest(s)
	s = s.Append(changes.Replace{Before: types.S8("152"), After: changes.Nil})

	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on string stream", c)
	}

	s.Append(changes.Replace{Before: changes.Nil, After: types.S8("99")})
	if strong, c = strong.Next(); c != nil {
		t.Error("Unexpected change on string stream", c, strong)
	}

}
