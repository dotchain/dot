// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

// Scope converts a DefMap into a scope/resolver
type Scope struct {
	// DefMap is not anonymous to ensure Scope is not treated as a Val or Def itself
	DefMap *DefMap
}

// Resolve implements the Resolver interface
func (s Scope) Resolve(key interface{}) (Def, Resolver) {
	if s.DefMap == nil {
		return nil, nil
	}

	if def, ok := (*s.DefMap)[key]; ok {
		return def, s
	}
	return nil, nil
}

// ChainResolver returns a resolver that tries the child and then parent.
func ChainResolver(child, parent Resolver) Resolver {
	return chainResolver{child, parent}
}

type chainResolver struct {
	child, parent Resolver
}

func (c chainResolver) Resolve(key interface{}) (Def, Resolver) {
	def, r := c.child.Resolve(key)
	if r == nil {
		def, r = c.parent.Resolve(key)
	}
	return def, r
}
