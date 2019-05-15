// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package fred

// Scope converts a DefMap into a scope/resolver
type Scope struct {
	// DefMap is not anonymous to ensure Scope is not treated as a Val or Def itself
	DefMap *DefMap
}

// Resolve implements part of the Resolver interface
func (s Scope) Resolve(key interface{}) Def {
	if s.DefMap == nil {
		return nil
	}

	if def, ok := (*s.DefMap)[key]; ok {
		return def
	}
	return nil
}
