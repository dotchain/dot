// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// The tests in the package uses the following convention to
// represent sequence mutations:
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
// Sequence merge tests then are simply a table of three strings:
// The first string is a left "mutation", the second is a right
// "mutation" and the last is the expected output.
//
// So, all the tests are descrived via a collection of 3-tuples
package dot_test

import (
	"fmt"
	"github.com/dotchain/dot"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"testing"
)

type one []string

var table = []one{
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

	// Series A.1 - conflicting splices
	{"-abc[123|d]4fgh-", "-abc12[34|e]fgh-", "-abcdefgh-"},
	{"-abc[123|]4fgh-", "-abc12[34|e]fgh-", "-abcefgh-"},
	{"-abc[123|d]4fgh-", "-abc12[34|]fgh-", "-abcdfgh-"},
	{"-abc[1234|d]fgh-", "-abc12[34|e]fgh-", "-abcdfgh-"},
	{"-abc[1234|d]fgh-", "-abc[12|e]34fgh-", "-abcdfgh-"},
	{"-abc[1234|d]fgh-", "-abc1[23|e]4fgh-", "-abcdfgh-"},

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

// Runs through the sequences table and validates all the sequences pass
func TestSimpleSequences(t *testing.T) {
	asArray := false
	for _, single := range table {
		testMutation(t, single[0], single[1], single[2], asArray)
	}
}

// Runs through the sequences table but converts them to "arrays" instead and tests that
func TestSimpleSequencesAsArrays(t *testing.T) {
	asArray := true
	for _, single := range table {
		t.Run(single[0]+" vs "+single[1], func(t *testing.T) {
			testMutation(t, single[0], single[1], single[2], asArray)
		})
	}
}

// tests a splice with an operation having a conflicting path
func TestSpliceConflictingPath(t *testing.T) {
	inner := "hello world"
	input := []interface{}{inner, inner, inner, inner}
	insert := []interface{}{"yo"}
	emptyArray := []interface{}{}
	emptyPath := []string{}

	// splice inserts at offset = 1
	t.Run("Insert", func(t *testing.T) {
		splice := []dot.Change{
			{Path: emptyPath, Splice: &dot.SpliceInfo{Offset: 1, Before: emptyArray, After: insert}},
		}
		innerSpliceInfo := &dot.SpliceInfo{Offset: 2, Before: "llo", After: "LLO"}
		for offset := range input {
			t.Run(fmt.Sprintf("Offset %v", offset), func(t *testing.T) {
				path := []string{fmt.Sprintf("%v", offset)}
				innerSplice := []dot.Change{{Path: path, Splice: innerSpliceInfo}}
				modified := []interface{}{inner, inner, inner, inner}
				modified[offset] = "heLLO world"
				output := concat(modified[:1], insert, modified[1:])
				testOps(t, input, output, splice, innerSplice)
			})
		}
	})

	// splice deletes at offset = 1 and 2
	t.Run("Delete", func(t *testing.T) {
		splice := []dot.Change{
			{Path: emptyPath, Splice: &dot.SpliceInfo{Offset: 1, Before: []interface{}{inner, inner}, After: emptyArray}},
		}
		innerSpliceInfo := &dot.SpliceInfo{Offset: 2, Before: "llo", After: "LLO"}
		for offset := range input {
			t.Run(fmt.Sprintf("Offset %v", offset), func(t *testing.T) {
				path := []string{fmt.Sprintf("%v", offset)}
				innerSplice := []dot.Change{{Path: path, Splice: innerSpliceInfo}}
				modified := []interface{}{inner, inner, inner, inner}
				modified[offset] = "heLLO world"
				output := concat(modified[:1], modified[3:])
				testOps(t, input, output, splice, innerSplice)
			})
		}
	})

	// splice replaces 1-2 with insert
	t.Run("Replace", func(t *testing.T) {
		splice := []dot.Change{
			{Path: emptyPath, Splice: &dot.SpliceInfo{Offset: 1, Before: []interface{}{inner, inner}, After: insert}},
		}
		innerSpliceInfo := &dot.SpliceInfo{Offset: 2, Before: "llo", After: "LLO"}
		for offset := range input {
			t.Run(fmt.Sprintf("Offset %v", offset), func(t *testing.T) {
				path := []string{fmt.Sprintf("%v", offset)}
				innerSplice := []dot.Change{{Path: path, Splice: innerSpliceInfo}}
				modified := []interface{}{inner, inner, inner, inner}
				modified[offset] = "heLLO world"
				output := concat(modified[:1], insert, modified[3:])
				testOps(t, input, output, splice, innerSplice)
			})
		}
	})
}

// tests a move with an operation having a conflicting path
func TestMoveConflictingPath(t *testing.T) {
	inner := "hello world"
	input := []interface{}{inner, inner, inner, inner, inner}

	emptyPath := []string{}

	// move 2-3 to left by one
	t.Run("Left", func(t *testing.T) {
		move := []dot.Change{
			{Path: emptyPath, Move: &dot.MoveInfo{Offset: 2, Count: 2, Distance: -1}},
		}
		innerSpliceInfo := &dot.SpliceInfo{Offset: 2, Before: "llo", After: "LLO"}
		for offset := range input {
			t.Run(fmt.Sprintf("Offset %v", offset), func(t *testing.T) {
				path := []string{fmt.Sprintf("%v", offset)}
				innerSplice := []dot.Change{{Path: path, Splice: innerSpliceInfo}}
				modified := []interface{}{inner, inner, inner, inner, inner}
				modified[offset] = "heLLO world"
				output := []interface{}{modified[0], modified[2], modified[3], modified[1], modified[4]}
				testOps(t, input, output, move, innerSplice)
			})
		}
	})

	// move 1-3 to right by one
	t.Run("Left", func(t *testing.T) {
		move := []dot.Change{
			{Path: emptyPath, Move: &dot.MoveInfo{Offset: 1, Count: 2, Distance: 1}},
		}
		innerSpliceInfo := &dot.SpliceInfo{Offset: 2, Before: "llo", After: "LLO"}
		for offset := range input {
			t.Run(fmt.Sprintf("Offset %v", offset), func(t *testing.T) {
				path := []string{fmt.Sprintf("%v", offset)}
				innerSplice := []dot.Change{{Path: path, Splice: innerSpliceInfo}}
				modified := []interface{}{inner, inner, inner, inner, inner}
				modified[offset] = "heLLO world"
				output := []interface{}{modified[0], modified[3], modified[1], modified[2], modified[4]}
				testOps(t, input, output, move, innerSplice)
			})
		}
	})
}

//
//  Helper routines
//

var mutationRE = regexp.MustCompile(`\[.*\]`)

func toFormattedArray(s string) []interface{} {
	var it norm.Iter
	result := []interface{}{}
	it.InitString(norm.NFKD, s)
	for !it.Done() {
		n := it.Next()
		elt := map[string]interface{}{
			"char":      string([]byte{n[0]}),
			"strike":    false,
			"underline": false,
		}
		for n = n[1:]; len(n) >= 2; n = n[2:] {
			// hardcode it
			if n[0] == 204 && n[1] == 182 {
				elt["strike"] = true
			} else if n[0] == 205 && n[1] == 159 {
				elt["underline"] = true
			}
		}
		result = append(result, elt)
	}
	return result
}

func toFormattedString(formatted []interface{}) string {
	result := ""
	for _, elt := range formatted {
		elt := elt.(map[string]interface{})
		strike := elt["strike"] == true
		underline := elt["underline"] == true
		ch, _ := elt["char"].(string)
		bytes := []byte(ch)
		if strike {
			bytes = append(bytes, 204, 182)
		}
		if underline {
			bytes = append(bytes, 205, 159)
		}
		result += string(bytes)
	}
	return result
}

func allStyles(formatted []interface{}) (strike, underline bool) {
	for _, elt := range formatted {
		elt := elt.(map[string]interface{})
		strike = strike || elt["strike"] == true
		underline = underline || elt["underline"] == true
	}
	return strike, underline
}

func isFormatted(left, right, expected string) bool {
	strike, underline := allStyles(toFormattedArray(left + right + expected))
	return strike || underline
}

func parseMutationAsFormattedArray(m string) (interface{}, interface{}, dot.Change) {
	input, output, mutation := parseMutation(m, false)
	inx, outx := toFormattedArray(input.(string)), toFormattedArray(output.(string))

	if mutation.Move != nil {
		offset, count, distance := mutation.Move.Offset, mutation.Move.Count, mutation.Move.Distance
		instr := input.(string)
		offset = len(toFormattedArray(instr[:offset]))
		count = len(toFormattedArray(instr[offset : offset+count]))
		if distance > 0 {
			distance = len(toFormattedArray(instr[offset+count : offset+count+distance]))
		} else {
			distance = -len(toFormattedArray(instr[offset+distance : offset]))
		}
		mutation.Move.Offset = offset
		mutation.Move.Count = count
		mutation.Move.Distance = distance
		return inx, outx, mutation
	}

	before, after := mutation.Splice.Before.(string), mutation.Splice.After.(string)
	inner := dot.Change{}
	change := dot.Change{}

	switch after {
	case "strike":
		inner.Set = &dot.SetInfo{Key: "strike", Before: false, After: true}
	case "unstrike":
		inner.Set = &dot.SetInfo{Key: "strike", Before: true, After: false}
	case "underline":
		inner.Set = &dot.SetInfo{Key: "underline", Before: false, After: true}
	case "ununderline":
		inner.Set = &dot.SetInfo{Key: "underline", Before: true, After: false}
	}

	if inner.Set != nil {
		outx = toFormattedArray(input.(string))
		change.Range = &dot.RangeInfo{
			Offset:  len(toFormattedArray(m[:mutation.Splice.Offset])),
			Count:   len(toFormattedArray(before)),
			Changes: []dot.Change{inner},
		}
		for index := change.Range.Offset; index < change.Range.Offset+change.Range.Count; index++ {
			elt := outx[index].(map[string]interface{})
			elt[inner.Set.Key] = inner.Set.After
		}
	} else {
		change.Splice = &dot.SpliceInfo{
			Offset: len(toFormattedArray(m[:mutation.Splice.Offset])),
			Before: toFormattedArray(before),
			After:  toFormattedArray(after),
		}
	}

	return inx, outx, change
}

func parseMutation(m string, asArray bool) (interface{}, interface{}, dot.Change) {
	loc := mutationRE.FindStringIndex(m)
	if len(loc) != 2 {
		panic("Invalid example: " + m)
	}
	inner := m[loc[0]+1 : loc[1]-1]
	left := m[:loc[0]]
	right := m[loc[1]:]
	if index := strings.IndexAny(inner, "|"); index > -1 {
		// it is a splice operation
		before := inner[:index]
		after := inner[index+1:]
		input := left + before + right
		output := left + after + right
		splice := &dot.SpliceInfo{Offset: len(left), Before: makeArray(before, asArray), After: makeArray(after, asArray)}
		return makeArray(input, asArray), makeArray(output, asArray), dot.Change{Splice: splice}
	} else if index := strings.IndexAny(left, "|"); index > -1 {
		input := left[:index] + left[index+1:] + inner + right
		output := left[:index] + inner + left[index+1:] + right
		move := &dot.MoveInfo{Offset: len(left) - 1, Count: len(inner), Distance: index - len(left) + 1}
		return makeArray(input, asArray), makeArray(output, asArray), dot.Change{Move: move}
	} else if index := strings.IndexAny(right, "|"); index > -1 {
		input := left + inner + right[:index] + right[index+1:]
		output := left + right[:index] + inner + right[index+1:]
		move := &dot.MoveInfo{Offset: len(left), Count: len(inner), Distance: index}
		return makeArray(input, asArray), makeArray(output, asArray), dot.Change{Move: move}
	} else {
		panic(fmt.Sprintf("Could not parse: %#v", m))
	}
}

func testMutation(t *testing.T, left, right, expected string, asArray bool) {
	linput, loutput, lop := parseMutation(left, asArray)
	rinput, routput, rop := parseMutation(right, asArray)
	match := makeArray(expected, asArray)
	formatted := false

	stringify := func(input interface{}) string {
		if formatted {
			return toFormattedString(input.([]interface{}))
		}
		if asArray {
			arr := input.([]interface{})
			b := []byte{}
			for _, bb := range arr {
				switch bb := bb.(type) {
				case float64:
					b = append(b, byte(bb))
				case uint8:
					b = append(b, bb)
				default:
					panic("Unknown type")
				}
			}
			return string(b)
		}
		return input.(string)
	}

	if isFormatted(left, right, expected) {
		if !asArray {
			return
		}
		formatted = true
		linput, loutput, lop = parseMutationAsFormattedArray(left)
		rinput, routput, rop = parseMutationAsFormattedArray(right)
		match = toFormattedArray(expected)
	}

	if !dot.Utils(x).AreSame(linput, rinput) {
		t.Errorf("Inputs do not match: %v %v (%v != %v)\n", left, right, stringify(linput), stringify(rinput))
		return
	}
	loutputActual := applyMany(linput, []dot.Change{lop})
	if !dot.Utils(x).AreSame(loutputActual, loutput) {
		t.Errorf("Output of %v is %v.  Expected %v", left, stringify(loutputActual), stringify(loutput))
		return
	}

	routputActual := applyMany(rinput, []dot.Change{rop})
	if !dot.Utils(x).AreSame(routputActual, routput) {
		t.Errorf("Output of %v is %v.  Expected %v", right, stringify(routputActual), stringify(routput))
		return
	}

	x := dot.Transformer{}
	left1, right1 := x.MergeChanges([]dot.Change{lop}, []dot.Change{rop})
	allLeft := append([]dot.Change{lop}, left1...)
	allRight := append([]dot.Change{rop}, right1...)

	resultLeft := applyMany(linput, allLeft)
	resultRight := applyMany(rinput, allRight)
	if !dot.Utils(x).AreSame(resultLeft, resultRight) {
		t.Errorf("Merge of %v and %v resulted in %v and %v resp.", left, right, stringify(resultLeft), stringify(resultRight))
	} else if !dot.Utils(x).AreSame(resultLeft, match) {
		t.Errorf("Merge of %v and %v = %v.  But expected %v", left, right, stringify(resultLeft), expected)
	}
}
