// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Command dotls lists the operations
//
// The argument can be a file name or a url
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/bolt"
	"github.com/dotchain/dot/ops/nw"
	"github.com/dotchain/dot/ops/pg"
)

var version = flag.Int("version", 0, "starting version")
var raw = flag.Bool("raw", false, "do not transform the operations")

func main() {
	flag.Parse()
	name := flag.Arg(0)
	if name == "" {
		log.Fatal("Missing file name or url argument")
	}

	var store ops.Store
	_, err := url.ParseRequestURI(name)

	switch {
	case err == nil:
		store = &nw.Client{URL: name}
	case strings.HasSuffix(strings.ToLower(name), ".bolt"):
		store, err = bolt.New(name, "dot_root", nil)
	default:
		store, err = pg.New(name, "dot_root", nil)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	if !*raw {
		store = ops.Transformed(store, &cache{map[int]ops.Op{}, map[int][]ops.Op{}})
	}

	ver := *version
	versions := map[interface{}]int{}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		result, err := store.GetSince(ctx, ver, 1000)
		cancel()

		if err != nil {
			log.Fatal(err)
		}
		if len(result) == 0 {
			break
		}
		for _, op := range result {
			versions[op.ID()] = op.Version()
			print(op, versions)
		}
		ver = result[len(result)-1].Version() + 1
	}
}

func print(op ops.Op, version map[interface{}]int) {
	c := formatChanges(nil, op.Changes(), "")
	if p := op.Parent(); p != nil {
		log.Printf("%d (%d %d) %s\n", op.Version(), op.Basis(), version[p], c)
	} else {
		log.Printf("%d (%d) %s\n", op.Version(), op.Basis(), c)
	}
}

func formatChanges(path []interface{}, c changes.Change, prefix string) string {
	switch c := c.(type) {
	case nil:
		return "<nil>"
	case changes.PathChange:
		path = append(append([]interface{}(nil), path...), c.Path...)
		return formatChanges(path, c.Change, prefix+"  ")
	case changes.ChangeSet:
		return formatChangeSet(path, c, prefix)
	case changes.Splice:
		return formatSplice(path, c, prefix)
	case changes.Replace:
		return formatReplace(path, c, prefix)
	}
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	return formatPath(path, string(b))
}

func formatChangeSet(path []interface{}, c changes.ChangeSet, prefix string) string {
	result := ""
	for _, cx := range c {
		if cx != nil {
			if result != "" {
				result += ",\n" + prefix
			}
			result += formatChanges(path, cx, prefix+"  ")
		}
	}
	return result
}

func formatSplice(path []interface{}, c changes.Splice, prefix string) string {
	if b, ok := c.Before.(types.Counter); ok {
		a, _ := c.After.(types.Counter)
		return formatPath(path, fmt.Sprintf("%d", a-b))
	}

	switch {
	case c.Before.Count() == 0 && c.After.Count() == 0:
		return formatPath(path, "<empty splice>")
	case c.Before.Count() == 0:
		return formatPath(path, fmt.Sprintf("insert %v at %d", c.After, c.Offset))
	case c.After.Count() == 0:
		return formatPath(path, fmt.Sprintf("delete %v at %d", c.Before, c.Offset))
	}
	return formatPath(path, fmt.Sprintf("%v => %v at %d", c.Before, c.After, c.Offset))
}

func formatReplace(path []interface{}, c changes.Replace, prefix string) string {
	switch {
	case c.Before == changes.Nil:
		return formatPath(path, fmt.Sprintf("set %v", c.After))
	case c.After == changes.Nil:
		return formatPath(path, fmt.Sprintf("remove %v", c.Before))
	}
	return formatPath(path, fmt.Sprintf("%v => %v", c.Before, c.After))
}

func formatPath(path []interface{}, s string) string {
	if len(path) == 0 {
		return s
	}

	return fmt.Sprintf("%v: %s", path, s)
}

type cache struct {
	x     map[int]ops.Op
	merge map[int][]ops.Op
}

func (c cache) Load(ver int) (ops.Op, []ops.Op) {
	return c.x[ver], c.merge[ver]
}

func (c cache) Store(ver int, op ops.Op, merge []ops.Op) {
	c.x[ver] = op
	c.merge[ver] = merge
}
