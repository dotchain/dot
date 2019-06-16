// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval

import (
	"strconv"
	"strings"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/data"
)

// Parse converts an expression into the corresponding AST
//
// The result is not yet evaluated, and needs a call to Eval to work.
func Parse(s Scope, code string) changes.Value {
	p := parser{
		Operators: map[string]*opInfo{
			",": {Priority: 1, New: combineArgs},
			"=": {Priority: 2, New: assign},

			// logical ops 1x priorities
			"|": {Priority: 10, New: simpleCall("|")},
			"&": {Priority: 11, New: simpleCall("&")},

			// comparison ops 2x priorities
			"<":  {Priority: 20, New: simpleCall("<")},
			">":  {Priority: 20, New: simpleCall(">")},
			"<=": {Priority: 20, New: simpleCall("<=")},
			">=": {Priority: 20, New: simpleCall(">=")},
			"==": {Priority: 20, New: simpleCall("==")},
			"!=": {Priority: 20, New: simpleCall("!=")},

			// arithmetic ops 3x priorities
			"+": {Priority: 30, Prefix: 0, New: simpleCall("+")},

			// grouping ops 4x priorities
			"(": {Priority: 40, BeginGroup: true},
			")": {Priority: 40, EndGroup: endGroup},

			// dot: max priority
			".": {
				Priority: 100,
				New: func(l, r changes.Value) changes.Value {
					if ref, ok := r.(*data.Ref); ok {
						r = ref.ID.(changes.Value)
					}
					return simpleCall(".")(l, r)
				},
			},
		},
		Error: func(offset int, message string) {
			panic(ParseError{offset, message})
		},
		StringTerm: func(s string, first rune) changes.Value {
			return types.S16(s)
		},
		NumericTerm: func(s string) changes.Value {
			n, _ := strconv.Atoi(strings.TrimSpace(s))
			return changes.Atomic{Value: n}
		},
		NameTerm: func(s string) changes.Value {
			return &data.Ref{ID: types.S16(s)}
		},
		CallTerm: func(fn, args changes.Value) changes.Value {
			a, ok := args.(types.A)
			if !ok {
				a = types.A{args}
			}
			return &Call{A: append(types.A{fn}, a...)}
		},
	}

	v, rest := p.parse(code, 0)
	if rest != "" {
		panic(ParseError{len(code) - len(rest), "unexpected character"})
	}

	return v
}

func endGroup(b *opInfo, isCall bool, term changes.Value, start, end int) changes.Value {
	if term == nil {
		// no term: pretend its an empty array
		return types.A{}
	}

	if !isCall {
		if obj, ok := callObject(term, start); ok {
			term = obj
		}
	}

	// FIX: a.(x) will now incorrectly
	// get evaluated to a.x.
	// Fix is a bit intricate
	// without extra unnecessary nodes
	return term
}

func callObject(args changes.Value, offset int) (changes.Value, bool) {
	errMissingKey := ParseError{offset, "missing key"}
	switch args := args.(type) {
	case vardef:
		return callObject(types.A{args}, offset)
	case types.A:
		foundKey := false
		foundNonKey := false
		dir := &data.Dir{Objects: types.M{}}
		obj := types.M{}
		dir.Root = obj
		for _, arg := range args {
			switch v := arg.(type) {
			case vardef:
				if foundNonKey {
					panic(errMissingKey)
				}
				foundKey = true
				dir.Objects[v.key.ID] = v.value
				obj[v.key.ID] = v.key
			case *data.Ref:
				if foundNonKey {
					panic(errMissingKey)
				}
				foundKey = true
				obj[v.ID] = v
			default:
				if foundKey {
					panic(errMissingKey)
				}
				foundNonKey = true
			}
		}
		if !foundKey {
			return nil, false
		}
		return dir, true
	}
	return nil, false
}

func assign(l, r changes.Value) changes.Value {
	return vardef{l.(*data.Ref), r}
}

func simpleCall(op string) func(l, r changes.Value) changes.Value {
	fn := &data.Ref{ID: types.S16(op)}
	return func(l, r changes.Value) changes.Value {
		return &Call{A: types.A{fn, l, r}}
	}
}

func combineArgs(l, r changes.Value) changes.Value {
	if x, ok := l.(types.A); ok {
		if r == nil {
			return x
		}
		return append(x, r)
	}

	if r == nil {
		return types.A{l}
	}

	return types.A{l, r}
}
