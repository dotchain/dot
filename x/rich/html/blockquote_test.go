// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package html_test

import (
	"reflect"
	"testing"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/x/rich"
	"github.com/dotchain/dot/x/rich/html"
)

func TestBlockQuoteApply(t *testing.T) {
	s1, s2 := rich.NewText("quote1"), rich.NewText("quote2")
	bq1 := html.BlockQuote{Text: &s1}
	bq2 := html.BlockQuote{Text: &s2}

	if x := bq1.Apply(nil, nil); !reflect.DeepEqual(x, bq1) {
		t.Error("Unexpected apply", x)
	}

	replace := changes.Replace{Before: bq2, After: bq2}
	if x := bq1.Apply(nil, replace); !reflect.DeepEqual(x, bq2) {
		t.Error("Unexpected replace", x)
	}

	c := changes.PathChange{
		Path:   []interface{}{"Text"},
		Change: changes.Replace{Before: s1, After: s2},
	}
	if x := bq1.Apply(nil, c).(html.BlockQuote); !reflect.DeepEqual(*x.Text, s2) {
		t.Error("Unexpected change", x)
	}
}

func TestBlockQuoteFormatHTML(t *testing.T) {
	bq := html.NewBlockQuote(rich.NewText("hello", html.FontBold))

	if x := html.Format(bq, nil); x != "<blockquote><b>hello</b></blockquote>" {
		t.Error("Unexpected", x)
	}
}
