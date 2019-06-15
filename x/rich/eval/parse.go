// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package eval

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/x/rich/data"
)

// Parse converts an expression into the corresponding AST
//
// The result is not yet evaluated, and needs a call to Eval to work.
func Parse(s Scope, code string) changes.Value {
	var err error
	p := parser{
		Operators: map[string]*opInfo{
			",": {Priority: 2, New: combineArgs},
			"+": {Priority: 11, Prefix: 0, New: simpleCall("+")},

			"(": {Priority: 20, BeginGroup: true},
			")": {
				Priority: 20,
				EndGroup: func(b *opInfo, term changes.Value, start, end int) changes.Value {
					// FIX: a.(x) will now incorrectly
					// get evaluated to a.x.
					// Fix is a bit intricate
					// without extra unnecessary nodes
					return term
				},
			},
			".": {
				Priority: 100,
				New: func(l, r changes.Value) changes.Value {
					if ref, ok := r.(*data.Ref); ok {
						r = ref.ID.(changes.Value)
					}
					return simpleCall(".")(l, r)
				},
			},
		},
		Error: func(offset int, message string) {
			err = fmt.Errorf("error at %d: %s", offset, message)
		},
		StringTerm: func(s string, first rune) changes.Value {
			return types.S16(s)
		},
		NumericTerm: func(s string) changes.Value {
			n, _ := strconv.Atoi(s)
			return changes.Atomic{Value: n}
		},
		NameTerm: func(s string) changes.Value {
			return &data.Ref{ID: types.S16(s)}
		},
		CallTerm: func(fn, args changes.Value) changes.Value {
			a, ok := args.(types.A)
			if !ok {
				a = types.A{args}
			}
			return &Call{A: append(types.A{fn}, a...)}
		},
	}

	v, rest := p.parse(code, 0)
	if rest != "" {
		err = fmt.Errorf("error at %d: unexpected charater", len(code)-len(rest))
	}
	if err != nil {
		panic(err)
	}
	return v
}

func simpleCall(op string) func(l, r changes.Value) changes.Value {
	fn := &data.Ref{ID: types.S16(op)}
	return func(l, r changes.Value) changes.Value {
		return &Call{A: append(types.A{fn}, combineArgs(l, r).(types.A)...)}
	}
}

func combineArgs(l, r changes.Value) changes.Value {
	if x, ok := l.(types.A); ok {
		return append(x, r)
	}
	return types.A{l, r}
}

type opInfo struct {
	New      func(d1, d2 changes.Value) changes.Value
	Priority int

	Prefix interface{}

	BeginGroup bool
	EndGroup   func(beginGroup *opInfo, d changes.Value, beginOffset, endOffset int) changes.Value
}

type parser struct {
	Operators map[string]*opInfo
	Error     func(offset int, message string)

	StringTerm  func(s string, first rune) changes.Value
	NumericTerm func(s string) changes.Value
	NameTerm    func(s string) changes.Value
	CallTerm    func(fn, args changes.Value) changes.Value
}

func (p *parser) parse(input string, offset int) (changes.Value, string) {
	lastWasTerm := false
	ops := []*opInfo{}
	terms := []changes.Value{}
	prefixes := 0

main:
	for input != "" {
		op, term, nextOffset, rest := p.scan(input, &offset)
		switch {
		case op != nil && op.EndGroup != nil:
			break main
		case op != nil && op.BeginGroup:
			term, nextOffset, rest = p.parseGroup(op, rest, nextOffset)
			if lastWasTerm {
				ops, terms = p.merge(ops, terms, op.Priority)
				terms[len(terms)-1] = p.CallTerm(terms[len(terms)-1], term)
				break
			}
			fallthrough
		case term != nil:
			p.failIf(lastWasTerm, offset, "unexpected term")
			for l := len(ops); prefixes > 0; prefixes, l = prefixes-1, l-1 {
				prefix := changes.Atomic{Value: ops[l-1].Prefix}
				term = ops[l-1].New(prefix, term)
				ops = ops[:l-1]
			}
			terms = append(terms, term)
			lastWasTerm = true
		case op != nil:
			if !lastWasTerm {
				p.failIf(op.Prefix == nil, offset, "unexpected op")
				prefixes++
			} else {
				ops, terms = p.merge(ops, terms, op.Priority)
				lastWasTerm = false
			}
			ops = append(ops, op)
		default:
			input = rest
			break main
		}
		offset, input = nextOffset, rest
	}

	if len(terms) == 0 && len(ops) == 0 {
		return nil, input
	}

	p.failIf(!lastWasTerm || prefixes > 0, offset, "incomplete")

	_, terms = p.merge(ops, terms, -1)
	return terms[0], input
}

func (p *parser) parseGroup(begin *opInfo, input string, offset int) (changes.Value, int, string) {
	group, rest := p.parse(input, offset)
	offset += len(input) - len(rest)
	op, _, nextOffset, rest := p.scan(rest, &offset)

	p.failIf(op == nil || op.EndGroup == nil, nextOffset, "unexpected char")

	return op.EndGroup(begin, group, offset, nextOffset), nextOffset, rest
}

func (p *parser) merge(ops []*opInfo, terms []changes.Value, pri int) ([]*opInfo, []changes.Value) {
	for t, l := len(terms), len(ops); l > 0 && ops[l-1].Priority >= pri; t, l = t-1, l-1 {
		terms[t-2] = ops[l-1].New(terms[t-2], terms[t-1])
		ops, terms = ops[:l-1], terms[:t-1]
	}
	return ops, terms
}

func (p *parser) scan(s string, offset *int) (*opInfo, changes.Value, int, string) {
	rest := strings.TrimLeftFunc(s, unicode.IsSpace)
	*offset += len(s) - len(rest)
	var op *opInfo
	var match string

	if rest == "" {
		return nil, nil, *offset, rest
	}

	for pattern, opx := range p.Operators {
		if len(pattern) <= len(match) || !strings.HasPrefix(rest, pattern) {
			continue
		}
		op, match = opx, pattern
	}

	if op != nil {
		return op, nil, *offset + len(match), rest[len(match):]
	}

	term, restx := p.scanTerm(rest, *offset)
	return nil, term, *offset + len(rest) - len(restx), restx
}

func (p *parser) scanTerm(s string, offset int) (changes.Value, string) {
	first, _ := utf8.DecodeRune([]byte(s))
	if first == '"' || first == '\'' {
		q, restx := p.scanQuote(s, offset)
		return p.StringTerm(q, first), restx
	}

	if unicode.IsDigit(first) {
		q, restx := p.scanNumeric(s)
		return p.NumericTerm(q), restx
	}

	if q, restx := p.scanID(s); q != "" {
		return p.NameTerm(q), restx
	}

	return nil, s
}

func (p *parser) scanQuote(s string, offset int) (quoted, rest string) {
	var first rune
	skip := false
	result := []rune{}
	for idx, r := range s {
		switch {
		case idx == 0:
			first = r
		case !skip && r == '\\':
			skip = true
		case !skip && r == first:
			return string(result), s[idx+utf8.RuneLen(r):]
		default:
			result = append(result, r)
			skip = false
		}
	}
	p.fail(offset, "incomplete quote")
	return "", ""
}

func (p *parser) scanNumeric(s string) (num, rest string) {
	var f float64
	_, err := fmt.Sscanf(s, "%f%s", &f, &rest)
	if err != io.EOF {
		idx := strings.Index(s, rest)
		return s[:idx], s[idx:]
	}
	return s, ""
}

func (p *parser) scanID(s string) (id, rest string) {
	for idx, r := range s {
		if !unicode.IsLetter(r) && (idx == 0 || !unicode.IsDigit(r)) {
			return s[:idx], s[idx:]
		}
	}
	return s, ""
}

func (p *parser) fail(offset int, message string) {
	p.Error(offset, message)
	panic(message)
}

func (p *parser) failIf(cond bool, offset int, message string) {
	if cond {
		p.Error(offset, message)
		panic(message)
	}
}
