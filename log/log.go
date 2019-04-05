// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package log defines the interface for loging within the DOT
// project.
//
// It is mainly used by the ops and ops/nw packages
package log

// Log is the default interface for logging used through out the DOT
// project. This allows callers to provide their implementation if
// needed.
type Log interface {
	Printf(fmt string, v ...interface{})
	Println(v ...interface{})
}

// Default returns a default logger that does not print anything
func Default() Log {
	return nolog{}
}

type nolog struct{}

func (n nolog) Printf(fmt string, v ...interface{}) {}
func (n nolog) Println(v ...interface{})            {}
