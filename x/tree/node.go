// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package tree

import "github.com/dotchain/dot/changes"

// Node uses a map to hold the attributes as well as the children.
// The "Children" key is used to hold the children which should be of
// type Nodes.
type Node map[string]interface{}

// WithKey sets the key for a node
func (n Node) WithKey(key interface{}) Node {
	return *(&n).update("Key", key)
}

// Key gets  the key for a node
func (n Node) Key() interface{} {
	key := n["Key"]
	return key
}

// Children returns all the child nodes
func (n Node) Children() Nodes {
	nn, _ := n["Children"].(Nodes)
	return nn
}

// Apply implements changes.Value
func (n *Node) Apply(c changes.Change) changes.Value {
	switch c := c.(type) {
	case nil:
		return n
	case changes.Replace:
		return c.After
	case changes.PathChange:
		if len(c.Path) == 0 {
			return n.Apply(c.Change)
		}
		if c.Path[0] == "Children" {
			return n.applyChildren(c.Path[1:], c.Change)
		}
		return n.applyAttribute(c.Path, c.Change)
	case changes.Custom:
		return c.ApplyTo(n)
	}
	panic("Unexpected change")
}

func (n *Node) applyAttribute(path []interface{}, c changes.Change) changes.Value {
	if len(path) != 1 {
		panic("Unexpected attribute path")
	}
	if v, ok := (*n)[path[0].(string)]; ok {
		return n.update(path[0].(string), (changes.Atomic{v}).Apply(c))
	}
	return n.update(path[0].(string), changes.Nil.Apply(c))
}

func (n *Node) applyChildren(path []interface{}, c changes.Change) changes.Value {
	c = changes.PathChange{path, c}
	if v, ok := (*n)["Children"].(changes.Value); ok {
		return n.update("Children", v.Apply(c))
	}
	return n.update("Children", changes.Nil.Apply(c))
}

func (n *Node) update(k string, v interface{}) *Node {
	result := Node{}
	for key, val := range *n {
		result[key] = val
	}
	if v == changes.Nil {
		delete(result, k)
	} else if atomic, ok := v.(changes.Atomic); ok {
		result[k] = atomic.Value
	} else {
		result[k] = v
	}
	return &result
}
