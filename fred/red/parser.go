// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package red implements a parser for fred
package red

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dotchain/dot/fred"
)

type opInfo struct {
	New      func(d1, d2 fred.Def) fred.Def
	Priority int

	PrefixTerm func() fred.Def

	BeginGroup bool
	EndGroup   func(beginGroup *opInfo, d fred.Def, beginOffset, endOffset int) fred.Def
}

type parser struct {
	Operators map[string]*opInfo
	Error     func(offset int, message string)

	StringTerm  func(s string, first rune) fred.Def
	NumericTerm func(s string) fred.Def
	NameTerm    func(s string) fred.Def
	CallTerm    func(fn, args fred.Def) fred.Def
}

func (p *parser) parse(input string, offset int) (fred.Def, string) {
	lastWasTerm := false
	ops := []*opInfo{}
	terms := []fred.Def{}
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
				term = ops[l-1].New(ops[l-1].PrefixTerm(), term)
				ops = ops[:l-1]
			}
			terms = append(terms, term)
			lastWasTerm = true
		case op != nil:
			if !lastWasTerm {
				p.failIf(op.PrefixTerm == nil, offset, "unexpected op")
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

func (p *parser) parseGroup(begin *opInfo, input string, offset int) (fred.Def, int, string) {
	group, rest := p.parse(input, offset)
	offset += len(input) - len(rest)
	op, _, nextOffset, rest := p.scan(rest, &offset)

	p.failIf(op == nil || op.EndGroup == nil, nextOffset, "unexpected char")

	return op.EndGroup(begin, group, offset, nextOffset), nextOffset, rest
}

func (p *parser) merge(ops []*opInfo, terms []fred.Def, pri int) ([]*opInfo, []fred.Def) {
	for t, l := len(terms), len(ops); l > 0 && ops[l-1].Priority >= pri; t, l = t-1, l-1 {
		terms[t-2] = ops[l-1].New(terms[t-2], terms[t-1])
		ops, terms = ops[:l-1], terms[:t-1]
	}
	return ops, terms
}

func (p *parser) scan(s string, offset *int) (*opInfo, fred.Def, int, string) {
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

func (p *parser) scanTerm(s string, offset int) (fred.Def, string) {
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
