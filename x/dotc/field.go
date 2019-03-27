// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dotc

import (
	"fmt"
	"strings"
)

var toValueFormats = map[string]string{
	"bool":    "changes.Atomic{%s}%.s",
	"int":     "changes.Atomic{%s}%.s",
	"string":  "types.S16(%s)%.s",
	"default": "%s%.s",
}

var fromValueFormats = map[string]string{
	"bool":    "(%s).(changes.Atomic).Value.(%s)",
	"int":     "(%s).(changes.Atomic).Value.(%s)",
	"string":  "string((%s).(types.S16))%.s",
	"default": "(%s).(%s)",
}

var fromStreamValueFormats = map[string]string{
	"bool":      "%s%.s",
	"int":       "%s%.s",
	"string":    "%s%.s",
	"types.S16": "string(%s)%.s",
	"types.S8":  "string(%s)%.s",
	"default":   "%s%.s",
}

var toStreamTypeFormats = map[string]string{
	"bool":      "streams.Bool%.s",
	"int":       "streams.Int%.s",
	"string":    "streams.S16%.s",
	"types.S16": "streams.S16%.s",
	"types.S8":  "streams.S8%.s",
	"default":   "%sStream",
}

// Field holds info for a struct field
type Field struct {
	Name, Key, Type                                           string
	Atomic                                                    bool
	ToValueFmt, FromValueFmt, ToStreamFmt, FromStreamValueFmt string
}

func (f Field) format(key, format string, m map[string]string) string {
	if format != "" {
		return fmt.Sprintf(format, key, f.Type)
	}

	if f.Atomic {
		return fmt.Sprintf(m["bool"], key, f.Type)
	}
	if x := m[f.Type]; x != "" {
		return fmt.Sprintf(x, key, f.Type)
	}
	return fmt.Sprintf(m["default"], key, f.Type)
}

// ToValue converts a strongly typed field to changes.Value
func (f Field) ToValue(recv, field string) string {
	key := recv
	if field != "" {
		key += "." + field
	}
	return f.format(key, f.ToValueFmt, toValueFormats)
}

// FromValue converts a changes.Value to the type of the field
func (f Field) FromValue(recv, field string) string {
	key := recv
	if field != "" {
		key += "." + field
	}
	return f.format(key, f.ToValueFmt, fromValueFormats)
}

// FromStreamType returns the name of the associated stream type
func (f Field) FromStreamValue(recv, field string) string {
	key := recv
	if field != "" {
		key += "." + field
	}
	return f.format(key, f.FromStreamValueFmt, fromStreamValueFormats)
}

// ToStreamType returns the name of the associated stream type
func (f Field) ToStreamType() string {
	format := f.ToStreamFmt
	if format == "" {
		format = toStreamTypeFormats[f.Type]
	}
	if format == "" {
		format = toStreamTypeFormats["default"]
	}
	t := f.Type
	for t[0] == '*' {
		t = t[1:]
	}
	return fmt.Sprintf(format, t)
}

// Setter returns the method name of the field setter
func (f Field) Setter() string {
	title := strings.Title(f.Name)
	if title == f.Name {
		return "Set" + title
	}
	return "set" + title
}
