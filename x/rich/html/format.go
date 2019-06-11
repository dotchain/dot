// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html

import (
	"strings"

	"github.com/dotchain/dot/x/rich"
)

// Formatter incrementally formats a text segment into html
type Formatter interface {
	Open(b *strings.Builder, last, current rich.Attrs, text string)
	Close(b *strings.Builder, last, current rich.Attrs, text string)
}

// Format formats rich text into html
func Format(t *rich.Text, f Formatter) string {
	var b strings.Builder
	FormatBuilder(&b, t, f)
	return b.String()
}

// FormatBuilder formats rich text into html
func FormatBuilder(b *strings.Builder, t *rich.Text, f Formatter) {
	if f == nil {
		f = DefaultFormatter
	}

	last := rich.Attrs{}
	for _, x := range *t {
		f.Close(b, last, x.Attrs, x.Text)
		f.Open(b, last, x.Attrs, x.Text)
		last = x.Attrs
	}
	if !last.Equal(rich.Attrs{}) {
		f.Close(b, last, rich.Attrs{}, "")
	}
}

// DefaultFormatter formats standard styles such as plain text string,
// bold and italics.
var DefaultFormatter = embedFmt{
	[]string{"Heading", "Link", "Image", "BlockQuote", "List"},
	simpleFmt{
		[]string{"FontStyle", "FontWeight"},
		textFmt{},
	},
}
