// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"time"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
)

// Config defines the configuration options for synchronization
type Config struct {
	// A reliable ops store.
	ops.Store

	// AutoTransform is off by default but can be set to true to
	// automatically transform the provided store
	AutoTransform bool
	ops.Cache

	// Session state
	Version int
	Pending []ops.Op

	// logger
	log.Log

	// Session state notifier
	Notify func(version int, pending []ops.Op)

	// Backoff configures the exponential backoff settings
	Backoff struct {
		Rand         func() float64
		Initial, Max time.Duration
	}
}

// Option can be used to configure Sync behavior
//
// See WithSession and WithCloseNotify
type Option func(c *Config)

// WithSession configures the connector to start with the provided
// version and pending instead of starting from scrach
func WithSession(version int, pending []ops.Op) Option {
	return func(c *Config) {
		c.Version = version
		c.Pending = pending
	}
}

// WithNotify configures a callback to be called when the
// session state changes. In particular, this is called when the
// session is closed.
func WithNotify(fn func(version int, pending []ops.Op)) Option {
	return func(c *Config) {
		c.Notify = fn
	}
}

// WithLog configures the logger to use
func WithLog(l log.Log) Option {
	return func(c *Config) {
		c.Log = l
	}
}

// WithBackoff configures the binary-exponential backoff settings
func WithBackoff(rng func() float64, initial, max time.Duration) Option {
	return func(c *Config) {
		c.Backoff.Rand = rng
		c.Backoff.Initial = initial
		c.Backoff.Max = max
	}
}

// WithAutoTransform specifies that the initial store yields
// untransformed operations and must be automatically transformed.
//
// The cache is required
func WithAutoTransform(cache ops.Cache) Option {
	return func(c *Config) {
		c.AutoTransform = true
		c.Cache = cache
	}
}
