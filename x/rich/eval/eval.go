// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package eval implements expression values.
//
// An expression value can be produced by use of Parse:
//
//   eval.Parse(scope, "list.map(value + 2)")
//   => equivalent to the following:
//
//   // first define all the tokens
//   dot := &data.Ref{ID: types.S16(".")}
//   plus := &data.Ref{ID: types.S16("+")}
//   list := &data.Ref{ID: types.S16("list")}
//   doMap := &data.Ref{ID: types.S16("map")}
//   value := &data.Ref{ID: types.S16("value")}
//   two := changes.Atomic{Value: 2}
//
//   now create a call expression for the above
//   expr := &data.Call{A: types.A{
//       // first element evaluates to list.map
//       &data.Call{A: types.A{dot, list, doMap}},
//       // second element represents value+2
//       &data.Call{A: types.A{plus, value, two}},
//   }}
//
// The scope contains the "globals" for evaluation, including all the
// operators. But methods on objects (such as "map" on arrays) are not
// part of the globals.
//
// An expression can be evaluated with Eval(). Eval converts call
// expression to the actual evaluated values using the provided
// scope.  Eval also walks the contents of any containers, replacing
// any expressions with their evaluated values.  Eval also honors
// references via Dir (so this is an easy way to create scopes).
package eval

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/data"
)

// Scope is any scope lookup function. Scope lookups are expected to
// return evaluated values but nested fields may contain unevaluated
// values.
type Scope func(v interface{}) changes.Value

// Eval evaluates a value within a particular scope
//
// Values of type CallExpr evaluate to the corresponding call.  All
// other values are simply the same values but with any nested call
// expressions evaluated.
//
// The provided scope is used to lookup IDs with Ref{} values. Any Dir
// values automatically contribute by creating a new scope
func Eval(s Scope, v changes.Value) changes.Value {
	switch v := v.(type) {
	case *Call:
		return v.Eval(s)
	case *data.Ref:
		return s(v.ID)
	case types.A:
		return evalArray(s, v)
	case *data.Dir:
		return Eval((&dirScope{s, v.Objects, nil}).lookup, v.Root)
	}
	return v
}

func evalArray(s Scope, v types.A) changes.Value {
	result := make(types.A, len(v))
	for kk, elt := range v {
		result[kk] = Eval(s, elt)
	}
	return result
}

type dirScope struct {
	base       Scope
	objects    types.M
	inProgress map[interface{}]bool
}

func (s *dirScope) lookup(id interface{}) changes.Value {
	val, ok := s.objects[id]
	if !ok {
		return s.base(id)
	}

	if s.inProgress[id] {
		panic("recursion detected")
	}
	if s.inProgress == nil {
		s.inProgress = map[interface{}]bool{}
	}
	s.inProgress[id] = true
	defer func() {
		s.inProgress[id] = false
	}()
	return Eval(s.lookup, val)
}
