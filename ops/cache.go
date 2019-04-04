// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops

// Cache represents a generic cache layer. This is intentionally a
// subset of sync.Map
type Cache interface {
	Load(key interface{}) (interface{}, bool)
	Store(key, value interface{})
}
