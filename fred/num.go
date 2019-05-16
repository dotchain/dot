// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"math/big"

	"github.com/dotchain/dot/changes"
)

// ErrNotNumber is returned when doing arithmetic on non-numbers
var ErrNotNumber = Error("not a number")

// Num implements Val
type Num string

// Apply implements changes.Value
func (n Num) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if c == nil {
		return n
	}
	if replace, ok := c.(changes.Replace); ok {
		return replace.After
	}
	return c.(changes.Custom).ApplyTo(ctx, n)
}

// Text implements Val.Text
func (n Num) Text() string {
	return string(n)
}

// Visit implements Val.Visit
func (n Num) Visit(v Visitor) {
	v.VisitLeaf(n)
}

// Field implements the "method" fields which are +/- etc
func (n Num) Field(e Env, key Val) Val {
	switch key {
	case Text("+"):
		return n.method(func(r1, r2 *big.Rat) { r1.Add(r1, r2) })
	case Text("-"):
		return n.method(func(r1, r2 *big.Rat) { r1.Sub(r1, r2) })
	case Text("/"):
		return n.method(func(r1, r2 *big.Rat) { r2.Inv(r2); r1.Mul(r1, r2) })
	case Text("*"):
		return n.method(func(r1, r2 *big.Rat) { r1.Mul(r1, r2) })
	}
	return ErrNoSuchField
}

func (n Num) method(fn func(r1, r2 *big.Rat)) Val {
	return method(func(e Env, defs *Defs) Val {
		var result big.Rat
		if err := result.UnmarshalText([]byte(string(n))); err != nil {
			return Error(err.Error())
		}
		for _, arg := range *defs {
			r, errval := n.toNum(e, arg)
			if errval != nil {
				return errval
			}
			fn(&result, r)
		}
		s, err := result.MarshalText()
		if err != nil {
			return Error(err.Error())
		}
		return Num(string(s))
	})
}

func (n Num) toNum(e Env, arg Def) (*big.Rat, Val) {
	x := arg.Eval(e)
	nx, ok := x.(Num)
	if !ok {
		if err, ok := x.(Error); ok {
			return nil, err
		}
		return nil, ErrNotNumber
	}
	var r big.Rat
	if err := r.UnmarshalText([]byte(string(nx))); err != nil {
		return nil, Error(err.Error())
	}
	return &r, nil
}
