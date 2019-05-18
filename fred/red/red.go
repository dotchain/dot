// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package red implements a parser for fred
package red

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/dotchain/dot/fred"
)

// Parse parses an expression and returns a Def for it.
//
// If there are errors it justs returns an expression that evaluates to the error
func Parse(s string) fred.Def {
	var e fred.Def
	p := &parser{
		Operators: defaultOps,
		Error: func(offset int, message string) {
			e = fred.Fixed(fred.Error(strconv.Itoa(offset) + ": " + message))
		},
		StringTerm: func(s string, first rune) fred.Def {
			return fred.Fixed(fred.Text(s))
		},
		NumericTerm: func(s string) fred.Def {
			var n big.Rat
			_, err := fmt.Sscan(s, &n)
			must(err)
			bytes, err := n.MarshalText()
			must(err)
			return fred.Fixed(fred.Num(string(bytes)))
		},
		NameTerm: func(s string) fred.Def {
			return fred.Ref(fred.Fixed(fred.Text(s)))
		},
		CallTerm: func(fn, args fred.Def) fred.Def {
			if args == nil {
				return fred.Call(fn)
			}
			if c, ok := args.(*Comma); ok {
				return fred.Call(fn, c.Args()...)
			}
			return fred.Call(fn, args)
		},
	}

	var retval fred.Def
	func() {
		defer func() {
			if r := recover(); r != nil && e == nil {
				panic(r)
			}
		}()
		result, rest := p.parse(s, 0)
		if rest != "" {
			p.Error(len(s)-len(rest), "unexpected char")
		}
		retval = result
	}()

	if e != nil {
		return e
	}
	if retval == nil {
		return fred.Nil()
	}
	return retval
}

var defaultOps = map[string]*opInfo{
	",": {
		Priority: 2,
		New: func(d1, d2 fred.Def) fred.Def {
			return &Comma{Left: d1, Right: d2}
		},
	},
	"|":  simpleOp("|", 3),
	"&":  simpleOp("&", 5),
	"<":  simpleOp("<", 10),
	">":  simpleOp(">", 10),
	"<=": simpleOp("<=", 10),
	">=": simpleOp(">=", 10),
	"==": simpleOp("==", 10),
	"!=": simpleOp("==", 10),
	"+":  withPrefix(simpleOp("+", 11), fred.Num("0")),
	"-":  withPrefix(simpleOp("-", 11), fred.Num("0")),
	"*":  simpleOp("*", 12),
	"/":  simpleOp("/", 12),
	"(":  {Priority: 20, BeginGroup: true},
	")": {
		Priority: 20,
		EndGroup: func(b *opInfo, term fred.Def, start, end int) fred.Def {
			return term
		},
	},
	".": {
		Priority: 100,
		New: func(d1, d2 fred.Def) fred.Def {
			return fred.Field(d1, unwrapField(d2))
		},
	},
}

func withPrefix(op *opInfo, v fred.Val) *opInfo {
	op.PrefixTerm = func() fred.Def {
		return fred.Fixed(v)
	}
	return op
}

func simpleOp(s string, priority int) *opInfo {
	return &opInfo{
		New: func(d1, d2 fred.Def) fred.Def {
			if def := getCompactedCall(s, d1, d2); def != nil {
				return def
			}
			return fred.Call(fred.Field(d1, fred.Fixed(fred.Text(s))), d2)
		},
		Priority: priority,
	}
}

func getCompactedCall(s string, d1, d2 fred.Def) fred.Def {
	fn, ok := d1.(funk)
	if !ok {
		return nil
	}

	field, ok := fn.Func().(fielder)
	if !ok {
		return nil
	}

	args := field.Args()
	if args == nil || len(*args) != 1 {
		return nil
	}

	arg, ok := (*args)[0].(fixed)
	if !ok {
		return nil
	}

	if arg.Val() != fred.Text(s) {
		return nil
	}

	args = fn.Args()
	clone := append(fred.Defs(nil), *args...)
	clone = append(clone, d2)
	return fred.Call(field, clone...)
}

func unwrapField(d fred.Def) fred.Def {
	p, ok := d.(*fred.Pure)
	if !ok {
		return d
	}

	_, ok = p.Functor.(ref)
	if !ok {
		return d
	}
	f, ok := (*p.Args)[0].(fixed)
	if !ok {
		return d
	}
	_, ok = f.Val().(fred.Text)
	if !ok {
		return d
	}
	return f
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type funk interface {
	Func() fred.Def
	Args() *fred.Defs
}

type fielder interface {
	fred.Def
	FieldBase() fred.Def
	Args() *fred.Defs
}

type fixed interface {
	fred.Def
	Val() fred.Val
}

type ref interface {
	Ref()
}
