// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package seqtest contains a table of human readable sequence tests
// and some helper methods to run tests from it.
//
// Each test is represented by three strings: the left change, the
// right change and the final output.
//
// Each change string stores a string representation of the change.
// Splices are represented by using square brackets to indicate
// the region of text that is being mutated.  Within the square
// brackets the original text before mutation is on the left of
// the vertical | sign with the mutated value to the right.
//
// So, a splice operation that takes:
//   Hello World
// and changes it to
//   Hello cruel world
// would look like this
//   Hello [W|cruel w]orld
//
// Insertions would basically have a blank value to the left of
// the or | sign and deletions would have a blank value to the
// right.
// Moves use a similar scheme of tagging the items being moved
// with square brackets but the location of the destination is
// shown with the | sign.
//
// So, a move operation that takes:
//   Hello bad big world!
// and changes it to
//   Hello big bad wrold:
// would look something like:
//   Hello |bad [big ]world!
//
//
// Range operations are a bit more complicated to encode. A range
// change looks like a splice except the replacement string is one of
// the following four: "strike", "unstrike", "underline" and
// "ununderline".  These replacements indicate the actual inner
// operation in the range: setting or unsetting the specified
// attributes for the range.

// Package seqtest implements a standard suite of validations.
package seqtest

import (
	"regexp"
	"strings"
)

// StringMutator is the interface to create changes
type StringMutator interface {
	Splice(initial string, offset, count int, replacement string) interface{}
	Move(initial string, offset, count, distance int) interface{}
	Range(initial string, offset, count int, attribute string) interface{}
}

// Validator is the interface to actually implement changes
type Validator func(name, initial string, left, right interface{}, merged string)

// ForEachTest iterates through all the tests
func ForEachTest(left, right StringMutator, validate Validator) {
	for _, test := range sequenceTests {
		l, initial := parse(test[0], left)
		r, _ := parse(test[1], right)

		name := test[0] + " x " + test[1]
		validate(name, initial, l, r, test[2])
	}
}

var mutationRE = regexp.MustCompile(`\[.*\]`)
var attributes = map[string]bool{
	"strike":      true,
	"unstrike":    true,
	"underline":   true,
	"ununderline": true,
}

func parse(s string, m StringMutator) (interface{}, string) {
	indices := mutationRE.FindStringIndex(s)
	inner := s[indices[0]+1 : indices[1]-1]
	left, right := s[:indices[0]], s[indices[1]:]

	if pipe := strings.IndexAny(inner, "|"); pipe >= 0 {
		before, after := inner[:pipe], inner[pipe+1:]
		initial := left + before + right
		if attributes[after] {
			return m.Range(initial, indices[0], len(before), after), initial
		}
		return m.Splice(initial, indices[0], len(before), after), initial
	}

	if pipe := strings.IndexAny(left, "|"); pipe >= 0 {
		initial := left[:pipe] + left[pipe+1:] + inner + right
		distance := len(left) - 1 - pipe
		return m.Move(initial, len(left)-1, len(inner), -distance), initial
	}

	pipe := strings.IndexAny(right, "|")
	initial := left + inner + right[:pipe] + right[pipe+1:]
	return m.Move(initial, len(left), len(inner), pipe), initial
}

var sequenceTests = [][3]string{
	// Series A: non conflicting splice splice actions
	{"-abc[123|d]efg-", "-abc123[efg|EFG]-", "-abcdEFG-"},
	{"-abc[123|d]efg-", "-abc123[|EFG]efg-", "-abcdEFGefg-"},
	{"-abc[123|d]efg-", "-abc123[efg|]-", "-abcd-"},

	{"-abc123[|d]efg-", "-abc123[efg|EFG]-", "-abc123dEFG-"},
	{"-abc123[|d]efg-", "-abc123[|EFG]efg-", "-abc123dEFGefg-"},
	{"-abc123[|d]efg-", "-abc123[efg|]-", "-abc123d-"},

	{"-abc[123|]efg-", "-abc123[efg|EFG]-", "-abcEFG-"},
	{"-abc[123|]efg-", "-abc123[|EFG]efg-", "-abcEFGefg-"},
	{"-abc[123|]efg-", "-abc123[efg|]-", "-abc-"},

	{"-abc123[efg|EFG]-", "-abc[123|d]efg-", "-abcdEFG-"},
	{"-abc123[|EFG]efg-", "-abc[123|d]efg-", "-abcdEFGefg-"},
	{"-abc123[efg|]-", "-abc[123|d]efg-", "-abcd-"},

	// sstart = 2/2, ostart = 0/2

	{"-abc123[efg|EFG]-", "-abc123[|d]efg-", "-abc123dEFG-"},
	{"-abc123[|EFG]efg-", "-abc123[|d]efg-", "-abc123EFGdefg-"},
	{"-abc123[efg|]-", "-abc123[|d]efg-", "-abc123d-"},

	{"-abc123[efg|EFG]-", "-abc[123|]efg-", "-abcEFG-"},
	{"-abc123[|EFG]efg-", "-abc[123|]efg-", "-abcEFGefg-"},
	{"-abc123[efg|]-", "-abc[123|]efg-", "-abc-"},

	// Series A.1 - conflicting splices
	{"-abc[123|d]4fgh-", "-abc12[34|e]fgh-", "-abcdefgh-"},
	{"-abc[123|]4fgh-", "-abc12[34|e]fgh-", "-abcefgh-"},

	{"-abc[123|d]4fgh-", "-abc[1234|e]fgh-", "-abcefgh-"},
	{"-abc[123|]4fgh-", "-abc[1234|e]fgh-", "-abcefgh-"},

	{"-abc12[34|d]fgh-", "-abc[123|e]4fgh-", "-abcedfgh-"},
	{"-abc12[34|]fgh-", "-abc[123|e]4fgh-", "-abcefgh-"},

	{"-abc[1234|d]fgh-", "-abc12[34|e]fgh-", "-abcdfgh-"},
	{"-abc12[34|e]fgh-", "-abc[1234|d]fgh-", "-abcdfgh-"},

	{"-abc[1234|d]fgh-", "-abc[12|e]34fgh-", "-abcdfgh-"},
	{"-abc[12|e]34fgh-", "-abc[1234|d]fgh-", "-abcdfgh-"},

	{"-abc[1234|d]fgh-", "-abc1[23|e]4fgh-", "-abcdfgh-"},
	{"-abc1[23|e]4fgh-", "-abc[1234|d]fgh-", "-abcdfgh-"},

	// Series B: non conflicting simple move move actions

	{"-abc[123]d|e456fg-", "-abc123d|e[456]fg-", "-abcd456123efg-"}, // unexpected but ok

	{"-abc[123]d|e456fg-", "-abc123d|e[456]fg-", "-abcd456123efg-"},
	{"-abc[123]d|456efg-", "-abc123d[456]e|fg-", "-abcd123e456fg-"},
	{"-ab|c[123]d456efg-", "-abc123|d[456]efg-", "-ab123c456defg-"},
	{"-ab|c[123]456def-", "-abc123[456]d|ef-", "-ab123cd456ef-"},

	// Series C: non conflicting move actions that span over each other

	{"-abc[123]de|456fg-", "-abc123d|e[456]fg-", "-abcd456e123fg-"},
	{"-abc[123]d|e456fg-", "-abc123|de[456]fg-", "-abc456d123efg-"},
	{"-abc[123]de|456fg-", "-abc123|de[456]fg-", "-abc456de123fg-"},
	{"-abc[123]de456|fg-", "-abc|123de[456]fg-", "-abc456de123fg-"},
	{"-abc[123]de456f|g-", "-ab|c123de[456]fg-", "-ab456cdef123g-"},
	{"-abc[123]de456f|g-", "-abc123de[456]fg|-", "-abcdef123g456-"},

	// conflicting move move tests
	{"-abc[123]de4|56fg-", "-abc1|23de[456]fg-", "-abc456de123fg-"},

	// Series D: non-conflicting splice vs move
	{"-abc[123|d]e456f-", "-abc123|e[456]f-", "-abcd456ef-"},
	{"-abc[123|d]e456f-", "-abc123e[456]f|-", "-abcdef456-"},
	{"-abc[123|d]e456f-", "-abc|123e[456]f-", "-abc456def-"},
	{"-abc[123|d]e456f-", "-ab|c123e[456]f-", "-ab456cdef-"},
	{"-abc[123|]e456f-", "-abc123|e[456]f-", "-abc456ef-"},
	{"-abc[123|]e456f-", "-abc123e[456]f|-", "-abcef456-"},
	{"-abc[123|]e456f-", "-abc|123e[456]f-", "-abc456ef-"},
	{"-abc[123|]e456f-", "-ab|c123e[456]f-", "-ab456cef-"},
	{"-abc123[|d]e456f-", "-abc123|e[456]f-", "-abc123456def-"}, // ok.  can also be abc1234d456f
	{"-abc123[|d]e456f-", "-abc123e[456]f|-", "-abc123def456-"},
	{"-abc123[|d]e456f-", "-abc|123e[456]f-", "-abc456123def-"},
	{"-abc123[|d]e456f-", "-ab|c123e[456]f-", "-ab456c123def-"},

	{"-abc[123|d]456f-", "-abc123[456]f|-", "-abcdf456-"},
	{"-abc[123|d]456f-", "-abc|123[456]f-", "-abc456df-"},
	{"-abc[123|d]456f-", "-ab|c123[456]f-", "-ab456cdf-"},
	{"-abc[123|]456f-", "-abc123[456]f|-", "-abcf456-"},
	{"-abc[123|]456f-", "-abc|123[456]f-", "-abc456f-"},
	{"-abc[123|]456f-", "-ab|c123[456]f-", "-ab456cf-"},
	{"-abc123[|d]456f-", "-abc123[456]f|-", "-abc123df456-"},
	{"-abc123[|d]456f-", "-abc|123[456]f-", "-abc456123df-"},
	{"-abc123[|d]456f-", "-ab|c123[456]f-", "-ab456c123df-"},

	// Series E: non-conflicting move vs splice
	{"-ab|c[123]de456gh-", "-abc123de[456|f]gh-", "-ab123cdefgh-"},
	{"-abc[123]d|e456gh-", "-abc123de[456|f]gh-", "-abcd123efgh-"},
	{"-abc[123]de|456gh-", "-abc123de[456|f]gh-", "-abcde123fgh-"},
	{"-abc[123]de456|gh-", "-abc123de[456|f]gh-", "-abcdef123gh-"},
	{"-abc[123]de456g|h-", "-abc123de[456|f]gh-", "-abcdefg123h-"},

	{"-ab|c[123]de456gh-", "-abc123de[456|]gh-", "-ab123cdegh-"},
	{"-abc[123]d|e456gh-", "-abc123de[456|]gh-", "-abcd123egh-"},
	{"-abc[123]de|456gh-", "-abc123de[456|]gh-", "-abcde123gh-"},
	{"-abc[123]de456|gh-", "-abc123de[456|]gh-", "-abcde123gh-"},
	{"-abc[123]de456g|h-", "-abc123de[456|]gh-", "-abcdeg123h-"},

	{"-ab|c[123]degh-", "-abc123de[|f]gh-", "-ab123cdefgh-"},
	{"-abc[123]d|egh-", "-abc123de[|f]gh-", "-abcd123efgh-"},
	{"-abc[123]de|gh-", "-abc123de[|f]gh-", "-abcde123fgh-"},
	{"-abc[123]deg|h-", "-abc123de[|f]gh-", "-abcdefg123h-"},

	{"-ab|c[123]456ef-", "-abc123[456|d]ef-", "-ab123cdef-"},
	{"-abc[123]456|ef-", "-abc123[456|d]ef-", "-abcd123ef-"},
	{"-abc[123]456e|f-", "-abc123[456|d]ef-", "-abcde123f-"},

	{"-ab|c[123]456ef-", "-abc123[456|]ef-", "-ab123cef-"},
	{"-abc[123]456|ef-", "-abc123[456|]ef-", "-abc123ef-"},
	{"-abc[123]456e|f-", "-abc123[456|]ef-", "-abce123f-"},

	{"-ab|c[123]ef-", "-abc123[|d]ef-", "-ab123cdef-"},
	{"-abc[123]e|f-", "-abc123[|d]ef-", "-abcde123f-"},

	// Series F: Splices conflicting with moves.
	// Note a very large number of cases here are testable but the actual
	// results of these operations are not interesting so long as there is "convergence"
	// so the only real test being done is when one range fully contains the other.
	{"-abc[1234|d]ef-", "-abc[12]3|4ef-", "-abcdef-"},
	{"-abc[1234|d]ef-", "-ab|c[1234]ef-", "-abdcef-"},
	{"-abc[1234|]ef-", "-ab|c[1234]ef-", "-abcef-"},
	{"-abc[|d]ef-", "-a|b[ce]f-", "-acebdf-"}, // this is odd but ok
	{"-abc[|d]ef-", "-ab|c[]ef-", "-abcdef-"},
	{"-ab|c[1234]ef-", "-abc1[23|d]4ef-", "-ab1d4cef-"},

	// Range tests are a bit more involved.
	// 1. Each string is considered as an array (since range operations require elements to
	//    apply changes to)
	// 2. Each array is of the form of a set of three fields: char, strike and underline
	// 3. Char holds the byte code of the character. The other two hold formatting info.
	// 4. Range operations look like splice in the examples here but the replacement text
	//    is one of "strike", "unstrike", "underline", "ununderline".  These get translated
	//    to a range operation that uses the corresponding set mutation on all the elements
	// 5. Textual representation of strike through and underline is done via unicode.
	//    Use "go run cmd/underline/underline.go brick" to get the formatted version of "brick"
	// 6. The code has a sloppy way of detecting if the input is meant to be a regular splice or a
	//    range version.  It does this by detecting if any of the strings have strike through or
	//    underline formatting.  So, range tests sometimes have an artificial strike through in
	//    the leading character to trigger this.

	// Note that VI edits these characters very well.  Emacs is not quite as nice.

	// Range/Splice no conflict
	{"-yellow [brick|strike] road-", "-[yellow|YELLOW] brick road-", "-YELLOW b̶r̶i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-[yellow |YELLOW ]brick road-", "-YELLOW b̶r̶i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow brick[ road| ROAD]-", "-yellow b̶r̶i̶c̶k̶ ROAD-"},
	{"-yellow [brick|strike] road-", "-yellow brick [road|ROAD]-", "-yellow b̶r̶i̶c̶k̶ ROAD-"},

	// Range covers splice
	{"-yellow [brick|strike] road-", "-yellow [bri|BRI]ck road-", "-yellow B̶R̶I̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow b[ric|RIC]k road-", "-yellow b̶R̶I̶C̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow br[ick|ICK] road-", "-yellow b̶r̶I̶C̶K̶ road-"},

	// splice covers range.  Note that the first character in these has a strike through
	// to force going through the formatted parsing code
	{"-S̶yellow [brick|strike] road-", "-S̶[yellow brick road|YELLOW BRICK ROAD]-", "-S̶YELLOW BRICK ROAD-"},
	{"-S̶yellow [brick|strike] road-", "-S̶yellow [brick road|BRICK ROAD]-", "-S̶yellow BRICK ROAD-"},
	{"-S̶yellow [brick|strike] road-", "-S̶yellow [brick|BRICK] road-", "-S̶yellow BRICK road-"},
	{"-S̶yellow [brick|strike] road-", "-S̶[yellow brick|YELLOW BRICK] road-", "-S̶YELLOW BRICK road-"},

	// Range intersects splice
	{"-yellow [brick|strike] road-", "-yello[w br|W BR]ick road-", "-yelloW BRi̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow br[ick r|ICK R]oad-", "-yellow b̶r̶ICK Road-"},

	// Range/Range tests
	{"-yellow [brick|strike] road-", "-yellow [brick|underline] road-", "-yellow b̶͟r̶͟i̶͟c̶͟k̶͟ road-"},
	{"-[yellow|strike] brick road-", "-yellow brick [road|underline]-", "-y̶e̶l̶l̶o̶w̶ brick r͟o͟a͟d͟-"},
	{"-yellow [bric|strike]k road-", "-yellow b[rick|underline] road-", "-yellow b̶r̶͟i̶͟c̶͟k͟ road-"},
	{"-yellow [brick|strike] road-", "-yellow b[ric|underline]k road-", "-yellow b̶r̶͟i̶͟c̶͟k̶ road-"},

	// Range/Move tests
	// No conflicts
	{"-yellow [brick|strike] road-", "-[yellow] |brick road-", "- yellowb̶r̶i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "|-[yellow ]brick road-", "yellow -b̶r̶i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-[yellow] brick| road-", "- b̶r̶i̶c̶k̶yellow road-"},
	{"-yellow [brick|strike] road-", "-yellow brick| [road]-", "-yellow b̶r̶i̶c̶k̶road -"},
	{"-yellow [brick|strike] road-", "-yellow brick[ road]-|", "-yellow b̶r̶i̶c̶k̶- road"},
	{"-yellow [brick|strike] road-", "-yellow |brick[ road]-", "-yellow  roadb̶r̶i̶c̶k̶-"},

	// No conflicts but insertion
	{"-yellow [brick|strike] road-", "-[yellow] br|ick road-", "- b̶r̶yellowi̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-[yellow ]br|ick road-", "-b̶r̶yellow i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow br|ick[ road]-", "-yellow b̶r̶ roadi̶c̶k̶-"},
	{"-yellow [brick|strike] road-", "-yellow br|ick [road]-", "-yellow b̶r̶roadi̶c̶k̶ -"},

	// Move covers range
	{"-yellow [brick|strike] road-", "-yellow[ brick] |road-", "-yellow  b̶r̶i̶c̶k̶road-"},
	{"-yellow [brick|strike] road-", "-yellow| [brick ]road-", "-yellowb̶r̶i̶c̶k̶  road-"},
	{"-yellow [brick|strike] road-", "-|yellow[ brick ]road-", "- b̶r̶i̶c̶k̶ yellowroad-"},

	// Range covers move
	{"-yellow [brick|strike] road-", "-yellow b|r[ick] road-", "-yellow b̶i̶c̶k̶r̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow| [b]rick road-", "-yellowb̶ r̶i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow| [brick] road-", "-yellowb̶r̶i̶c̶k̶  road-"},

	// Range intersects move
	{"-yellow [brick|strike] road-", "-yellow b|r[ick ]road-", "-yellow b̶i̶c̶k̶ r̶road-"},
	{"-yellow [brick|strike] road-", "-yellow| br[ick ]road-", "-yellowi̶c̶k̶  b̶r̶road-"},
	{"-yellow [brick|strike] road-", "-yellow br[ick ]road|-", "-yellow b̶r̶roadi̶c̶k̶ -"},
	{"-yellow [brick|strike] road-", "-yellow[ b]r|ick road-", "-yellowr̶ b̶i̶c̶k̶ road-"},
	{"-yellow [brick|strike] road-", "-yellow[ b]rick |road-", "-yellowr̶i̶c̶k̶  b̶road-"},
	{"-yellow [brick|strike] road-", "-|yellow[ b]rick road-", "- b̶yellowr̶i̶c̶k̶ road-"},
}
