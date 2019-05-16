// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"strconv"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// ErrInvalidArgs is a generic error for mismatched number or types of args
var ErrInvalidArgs = Error("invalid args")

// Text uses types.S16 to implement string-based values
type Text string

// Apply implements changes.Value
func (t Text) Apply(ctx changes.Context, c changes.Change) changes.Value {
	if splice, ok := c.(changes.Splice); ok {
		splice.Before = types.S16(string(splice.Before.(Text)))
		splice.After = types.S16(string(splice.After.(Text)))
		c = splice
	}
	if custom, ok := c.(changes.Custom); ok {
		return custom.ApplyTo(ctx, t)
	}

	v := types.S16(string(t)).Apply(ctx, c)
	if x, ok := v.(types.S16); ok {
		return Text(x)
	}
	return v
}

// Count implements changes.Collection
func (t Text) Count() int {
	return types.S16(string(t)).Count()
}

// Slice implements changes.Collection
func (t Text) Slice(offset, count int) changes.Collection {
	v := types.S16(string(t)).Slice(offset, count)
	return Text(string(v.(types.S16)))
}

// ApplyCollection implements changes.Collection
func (t Text) ApplyCollection(ctx changes.Context, c changes.Change) changes.Collection {
	if splice, ok := c.(changes.Splice); ok {
		splice.Before = types.S16(string(splice.Before.(Text)))
		splice.After = types.S16(string(splice.After.(Text)))
		c = splice
	}

	v := types.S16(string(t)).ApplyCollection(ctx, c)
	return Text(string(v.(types.S16)))
}

// Text implements Val.Text
func (t Text) Text() string {
	return string(t)
}

// Visit implements Val.Visit
func (t Text) Visit(v Visitor) {
	v.VisitLeaf(t)
}

// Field implements the "method" fields which only has "concat" at this point.
func (t Text) Field(e Env, key Val) Val {
	switch key {
	case Text("concat"):
		return t.method(t.concatMethod)
	case Text("length"):
		return Num(strconv.Itoa(types.S16(string(t)).Count()))
	case Text("splice"):
		return t.method(t.spliceMethod)
	case Text("slice"):
		return t.method(t.sliceMethod)
	}
	return ErrNoSuchField
}

func (t Text) method(fn func(args Vals) Val) Val {
	return method(func(e Env, args *Defs) Val {
		others := args.Eval(e).(*Vals)
		var vals Vals
		if others != nil {
			vals = *others
		}
		return fn(vals)
	})
}

func (t Text) concatMethod(args Vals) Val {
	result := string(t)
	for _, arg := range args {
		result += arg.Text()
	}
	return Text(result)
}

func (t Text) sliceMethod(args Vals) Val {
	x := types.S16(string(t))

	if len(args) != 1 && len(args) != 2 {
		return ErrInvalidArgs
	}

	offset, err := ToInt(args[0])
	if err != nil {
		return err
	}

	count := int64(x.Count()) - offset
	if len(args) == 2 {
		var err Val
		count, err = ToInt(args[1])
		if err != nil {
			return err
		}
	}

	if offset < 0 || count < 0 || offset+count > int64(x.Count()) {
		return ErrInvalidArgs
	}

	return Text(string(x.Slice(int(offset), int(count)).(types.S16)))
}

func (t Text) spliceMethod(args Vals) Val {
	x := types.S16(string(t))

	if len(args) != 3 {
		return ErrInvalidArgs
	}

	offset, err := ToInt(args[0])
	if err != nil {
		return err
	}

	count, err := ToInt(args[1])
	if err != nil {
		return err
	}

	r := args[2].Text()
	if offset < 0 || count < 0 || offset+count > int64(x.Count()) {
		return ErrInvalidArgs
	}

	x = x.Apply(nil, changes.Splice{
		Offset: int(offset),
		Before: x.Slice(int(offset), int(count)),
		After:  types.S16(r),
	}).(types.S16)

	return Text(string(x))
}
