// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package vdom implement DOM reconciliation a la React
//
// The main export is a Reconciler which can be used to convert a
// MutableNode into the same shape as the provided  "virtual" Node.
package vdom

// Node is the interface for a virtual node. It is read-only and
// provides access to the Tag, attributes and child nodes.  There is
// no explicit node-type. It can either be encoded into the Key() or
// even the Tag() itself.
//
// The Key is expected to be unique for all nodes that share a
// parent. This is used to dedupe children.
type Node interface {
	Tag() string
	Key() interface{}
	ForEachAttribute(fn func(key, val string))
	ForEachNode(fn func(n Node))
}

// MutableNode represents an actual DOM node with mutable semantics.
type MutableNode interface {
	Node
	SetAttribute(key, val string)
	RemoveAttribute(key string)
	Children() MutableNodes
}

// MutableNodes represents the mutable children node list.  Next()
// returns the current head of the list and also advances.
//
// Insert inserts before the head of the list (leaving the head as
// is).
//
// Remove removes the current head of the list, advancing forward.
//
// Only one MutableNodes representation is in use at any given time
// for the children of a given node.
type MutableNodes interface {
	Next() MutableNode
	Insert(MutableNode)
	Remove() MutableNode
}

// Reconciler implements virtual dom reconciliation. This function is
// the constructor for new mutable nodes created.
type Reconciler func(tag string, key interface{}) MutableNode

// Reconcile remakses the mutable node in the shape of the provided
// virtual node. Note that if the root is itself modified (say,
// because the Key changed or some such reason), the function just
// returns the updated root node. The caller is expected to work with
// the parent of the current node and replace it.
func (r Reconciler) Reconcile(m MutableNode, n Node) MutableNode {
	if n == nil {
		return nil
	}

	if m == nil || m.Tag() != n.Tag() || m.Key() != n.Key() {
		return r.clone(n)
	}

	keys := map[string]bool{}
	n.ForEachAttribute(func(key, val string) {
		m.SetAttribute(key, val)
		keys[key] = true
	})
	deletions := []string(nil)
	m.ForEachAttribute(func(key, val string) {
		if _, ok := keys[key]; !ok {
			deletions = append(deletions, key)
		}
	})
	for _, key := range deletions {
		m.RemoveAttribute(key)
	}

	rx := &reconciler{Reconciler: r}
	rx.reconcileChildren(m, n)
	return m
}

func (r Reconciler) clone(n Node) MutableNode {
	result := r(n.Tag(), n.Key())
	n.ForEachAttribute(result.SetAttribute)
	children := result.Children()
	n.ForEachNode(func(child Node) {
		children.Insert(r.clone(child))
	})
	return result
}

// reconciler reconciles children
type reconciler struct {
	Reconciler
	before, after map[interface{}]bool
	stash         map[interface{}]Node
	keys          []interface{}
	nodes         MutableNodes
}

func (r *reconciler) reconcileChildren(m MutableNode, n Node) {
	r.stash = map[interface{}]Node{}
	r.before, r.after = r.toMap(m), r.toMap(n)
	r.keys = make([]interface{}, 0, len(r.before))
	m.ForEachNode(func(child Node) {
		r.keys = append(r.keys, child.Key())
	})
	r.nodes = m.Children()
	r.removeDeleted()
	n.ForEachNode(func(child Node) {
		for !r.handleChild(child) {
		}
	})
}

func (r *reconciler) handleChild(child Node) bool {
	defer r.removeDeleted()
	key := child.Key()

	if stashed, ok := r.stash[key]; ok {
		r.nodes.Insert(stashed.(MutableNode))
		r.Reconcile(stashed.(MutableNode), child)
	} else if _, ok := r.before[key]; !ok {
		r.nodes.Insert(r.clone(child))
	} else if len(r.keys) > 0 && r.keys[0] == key {
		r.keys = r.keys[1:]
		own := r.nodes.Next()
		r.Reconcile(own.(MutableNode), child)
	} else {
		node := r.nodes.Remove()
		r.keys = r.keys[1:]
		r.stash[node.Key()] = node
		return false
	}

	return true
}

func (r *reconciler) removeDeleted() {
	for len(r.keys) > 0 {
		if _, ok := r.after[r.keys[0]]; ok {
			return
		}
		r.keys = r.keys[1:]
		r.nodes.Remove()
	}
}

func (r *reconciler) toMap(n Node) map[interface{}]bool {
	result := map[interface{}]bool{}
	n.ForEachNode(func(child Node) {
		result[child.Key()] = true
	})
	return result
}
