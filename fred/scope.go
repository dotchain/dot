// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Scope implements a simple resolver
type Scope map[interface{}]Def

func (s Scope) set(key interface{}, value changes.Value) changes.Value {
	clone := Scope{}
	for k, v := range s {
		if k != key {
			clone[k] = v
		}
	}

	if value != changes.Nil {
		clone[key] = value.(Def)
	}
	return clone
}

func (s Scope) get(key interface{}) changes.Value {
	if v, ok := s[key]; ok {
		return v
	}
	return changes.Nil
}

// Apply implements the changes.Value interface
func (s Scope) Apply(ctx changes.Context, c changes.Change) changes.Value {
	return (types.Generic{Set: s.set, Get: s.get}).Apply(ctx, c, s)
}

// Resolve implements the Resolver interface
func (s Scope) Resolve(key interface{}) (Def, Resolver) {
	if def, ok := s[key]; ok {
		return def, s
	}
	return nil, nil
}

// WithParent returns a chained resolver
func (s Scope) WithParent(r Resolver) Resolver {
	return withParentScope{s, r}
}

type withParentScope struct {
	self   Scope
	parent Resolver
}

func (s withParentScope) Resolve(key interface{}) (Def, Resolver) {
	def, r := s.self.Resolve(key)
	if r == nil {
		def, r = s.parent.Resolve(key)
	}
	return def, r
}
