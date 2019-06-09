// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"strings"

	"github.com/dotchain/dot/x/rich"
)

type textFmt struct{}

func (t textFmt) Open(b *strings.Builder, last, current rich.Attrs, text string) {
	if text != "" {
		must(b.WriteString(text))
	}
}
func (t textFmt) Close(b *strings.Builder, last, current rich.Attrs, text string) {
}

type tagger interface {
	OpenTag() string
	CloseTag() string
}

type simpleFmt struct {
	keys []string
	base Formatter
}

func (s simpleFmt) Open(b *strings.Builder, last, current rich.Attrs, text string) {
	for _, key := range s.keys {
		if after := current[key]; after != nil {
			b.WriteString(after.(tagger).OpenTag())
		}
	}
	s.base.Open(b, last, current, text)
}

func (s simpleFmt) Close(b *strings.Builder, last, current rich.Attrs, text string) {
	for kk := range s.keys {
		if before := last[s.keys[len(s.keys)-kk-1]]; before != nil {
			b.WriteString(before.(tagger).CloseTag())
		}
	}
	s.base.Close(b, last, current, text)
}
