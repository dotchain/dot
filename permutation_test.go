// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot_test

import "github.com/dotchain/dot"

type permutation struct {
	journal  []dot.Operation
	children map[int]map[int]bool
	parents  map[int]map[int]bool

	// if permutaiton = [5, 2, 1], seen[5] = 0, seen[2] = 1 etc
	seen map[int]int

	// indices with zero, one and two parents
	zeros, ones, twos map[int]bool
}

func (p *permutation) permute() [][]dot.Operation {
	var results [][]dot.Operation
	if len(p.seen) == len(p.journal) {
		result := make([]dot.Operation, len(p.journal))
		for k, v := range p.seen {
			result[v] = p.journal[k]
		}
		return [][]dot.Operation{result}
	}

	choices := []int{}
	for k, v := range p.zeros {
		if v {
			choices = append(choices, k)
		}
	}

	for _, next := range choices {
		p.zeros[next] = false
		p.seen[next] = len(p.seen)

		for child := range p.children[next] {
			if p.ones[child] {
				p.ones[child] = false
				p.zeros[child] = true
			} else {
				p.twos[child] = false
				p.ones[child] = true
			}
		}
		results = append(results, p.permute()...)

		// backtrack
		p.zeros[next] = true
		delete(p.seen, next)
		for child := range p.children[next] {
			if p.zeros[child] {
				p.ones[child] = true
				p.zeros[child] = false
			} else {
				p.twos[child] = true
				p.ones[child] = false
			}
		}
	}
	return results
}

// GetPermutations returns all permutations of the operations
// that are valid. A valid permutation will have all operations
// following their basis + parent
//
// The algorithm is bruteforce with backtracking which is very
// slow for large datasets but it will suffice here for the
// size of the journal we use in tests..
func GetPermutations(journal []dot.Operation) [][]dot.Operation {
	p := &permutation{
		journal:  journal,
		seen:     map[int]int{},
		parents:  map[int]map[int]bool{},
		children: map[int]map[int]bool{},
		zeros:    map[int]bool{},
		ones:     map[int]bool{},
		twos:     map[int]bool{},
	}
	idx := map[string]int{}
	for kk, op := range journal {
		idx[op.ID] = kk
		p.parents[kk] = map[int]bool{}
		p.children[kk] = map[int]bool{}
	}

	for kk, op := range journal {
		count := 0
		parentID, basisID := op.ParentID(), op.BasisID()
		if basisID != "" && parentID != basisID {
			count++
			p.parents[kk][idx[basisID]] = true
			p.children[idx[basisID]][kk] = true
		}
		if parentID != "" {
			count++
			p.parents[kk][idx[parentID]] = true
			p.children[idx[parentID]][kk] = true
		}
		if count == 0 {
			p.zeros[kk] = true
		} else if count == 1 {
			p.ones[kk] = true
		} else {
			p.twos[kk] = true
		}
	}
	return p.permute()
}
