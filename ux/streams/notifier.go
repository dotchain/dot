// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package streams is a collection of frequently used stream types
// and related utilities
package streams

import (
	dots "github.com/dotchain/dot/streams"
)

// generate BoolStream and TextStream
//go:generate go run codegen.go

// Notifier is an alias
type Notifier = dots.Notifier

// Handler is an alias
type Handler = dots.Handler

// Cache is an alias
type Cache = dots.Cache
