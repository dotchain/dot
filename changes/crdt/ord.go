// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package crdt

import (
	"math/big"
	"strconv"
	"strings"
)

// NextOrd returns the next ordinal.
//
// Less(key, next(key)) is guaranteed to always be true
func NextOrd(key string) string {
	return fromString(key).next().toString()
}

// PrevOrd returns the prev ordinal.
//
// LessOrd(prev(key), key) is guaranteed to always be true
func PrevOrd(key string) string {
	return fromString(key).prev().toString()
}

// LessOrd returns true if a < b
func LessOrd(a, b string) bool {
	return fromString(a).less(fromString(b))
}

// BetweenOrd returns keys between the two provided keys.
func BetweenOrd(a, b string, count int) []string {
	ak, bk := fromString(a), fromString(b)
	if ak.less(bk) {
		return ak.between(bk, count)
	}
	return bk.between(ak, count)
}

func fromString(s string) *ordkey {
	if s == "" {
		return &ordkey{}
	}

	parts := strings.Split(s, ",")
	var k ordkey
	k.r.Num().SetString(parts[1], 62)
	logd, _ := strconv.Atoi(parts[0])
	k.r.Denom().Lsh(one, uint(logd-1))
	return &k
}

type ordkey struct {
	r big.Rat
}

func (k ordkey) toString() string {
	n, d := k.r.Num(), k.r.Denom()
	if d.Cmp(one) == 0 && n.BitLen() == 0 {
		return ""
	}
	return strconv.Itoa(d.BitLen()) + "," + n.Text(62)
}

func (k *ordkey) less(o *ordkey) bool {
	return k.r.Cmp(&o.r) < 0
}

func (k *ordkey) prev() *ordkey {
	var result ordkey
	result.r.Num().Quo(k.r.Num(), k.r.Denom())
	result.r.Num().Sub(result.r.Num(), one)
	return &result
}

func (k *ordkey) next() *ordkey {
	var result ordkey
	result.r.Num().Quo(k.r.Num(), k.r.Denom())
	result.r.Num().Add(result.r.Num(), one)
	return &result
}

var one = big.NewInt(1)

func (k *ordkey) between(o *ordkey, count int) []string {
	result := []string{}
	inc := big.NewRat(int64(1), int64(count+1))
	var diff big.Rat
	diff.Sub(&k.r, &o.r).Abs(&diff).Mul(&diff, inc)

	var last ordkey
	last.r.Set(&k.r)

	for count > 0 {
		last.r.Add(&last.r, &diff)
		result = append(result, last.toString())
		count--
	}
	return result
}
