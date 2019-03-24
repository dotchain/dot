// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package dotc implements code-generation tools for dot.changes
package dotc

import (
	"bytes"
	"go/format"

	"golang.org/x/tools/imports"
)

// Info tracks all information used for code generation
type Info struct {
	Package string
	Imports [][2]string
	Structs []Struct
	Unions  []Union
}

// Generate implements the helper methods for the provided types
func (info Info) Generate() (string, error) {
	var buf bytes.Buffer

	info.Imports = append(
		[][2]string{
			{"", "github.com/dotchain/dot/changes"},
			{"", "github.com/dotchain/dot/changes/types"},
		}, info.Imports...)

	if err := infoTpl.Execute(&buf, info); err != nil {
		return "", err
	}

	for _, s := range info.Structs {
		if err := s.GenerateApply(&buf); err != nil {
			return "", err
		}
		if err := s.GenerateSetters(&buf); err != nil {
			return "", err
		}
	}

	for _, u := range info.Unions {
		if err := u.GenerateApply(&buf); err != nil {
			return "", err
		}
		if err := u.GenerateSetters(&buf); err != nil {
			return "", err
		}
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), err
	}

	p2, err := imports.Process("generated.go", p, nil)
	if err != nil {
		return string(p), err
	}

	return string(p2), nil
}
