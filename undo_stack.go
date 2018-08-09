// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

const (
	opUndo   = "undo"
	opRedo   = "redo"
	opLocal  = "local"
	opRemote = "remote"
)

// UndoStack implements an undo stack that deals with intervening remote
// operations. This is an immutable type and all operations on it
// simply return a new Log without mutating the current log
//
// Concurrency
//
// This type is not safe for concurrent access.  While all the
// operations work in an immutable fashion and so no panics will
// happen if used concurrently, the correctness of the
// algorithm requires serialized access.
type UndoStack struct {
	Transformer
	Operations []Operation
	Types      []string
}

// Push appends an operation.  It ignores operations that have already
// been pushed.
func (u *UndoStack) Push(op Operation, isLocal bool) *UndoStack {
	for _, opx := range u.Operations {
		if opx.ID == op.ID {
			return u
		}
	}

	if isLocal {
		return u.push(op, opLocal)
	}
	return u.push(op, opRemote)
}

func (u *UndoStack) push(op Operation, opType string) *UndoStack {
	result := &UndoStack{Transformer: u.Transformer}
	result.Operations = append(u.Operations, op)
	result.Types = append(u.Types, opType)
	return result
}

// GetUndo returns an operation which when applied will have
// the effect of undoing the last operation on the client operation
// stack. The provided ID is used for the operation but other metadata
// associated with the operation (such as Parents) should be filled in
// by the caller.
//
// GetUndo can return nil if there are no operations to undo.
func (u *UndoStack) GetUndo(id string) (*Operation, *UndoStack) {
	skipCount := 0
	l := len(u.Operations) - 1
	for kk := range u.Operations {
		switch u.Types[l-kk] {
		case opRedo, opLocal:
			if skipCount == 0 {
				op := u.undo(l - kk)
				op.ID = id
				return &op, u.push(op, opUndo)
			}
			skipCount--
		case opUndo:
			skipCount++
		}
	}
	return nil, u
}

// GetRedo returns an operation which reverses the effect of the last
// undo on the stack.
//
// GetRedo returns nil if there is no undo operation to redo
func (u *UndoStack) GetRedo(id string) (*Operation, *UndoStack) {
	skipCount := 0
	l := len(u.Operations) - 1
	for kk := range u.Operations {
		switch u.Types[l-kk] {
		case opUndo:
			if skipCount == 0 {
				op := u.undo(l - kk)
				op.ID = id
				return &op, u.push(op, opRedo)
			}
			skipCount--
		case opRedo:
			skipCount++
		case opLocal:
			return nil, u
		}
	}
	return nil, u
}

// gets the result of undoing a single operation in the stack by
// merging the undo of the operation with the rest of the operations
// that follow
func (u *UndoStack) undo(i int) Operation {
	rest := u.simplify(u.Operations[i+1:], u.Types[i+1:])
	_, result := u.MergeOperations([]Operation{u.Operations[i].Undo()}, rest)
	return result[0]
}

// simplify removes all consecutive <undo/redo> pairs and <local/undo> pairs
// the simplification process guarantees that the result of simplification
// would have the same effect as the original chain but with the cumbersome
// op/undo pairs removed (since op/undo pairs often cause odd merge behavior)a
func (u *UndoStack) simplify(ops []Operation, types []string) []Operation {
	var result []Operation
	var resultTypes []string
	for kk, opType := range types {
		l := len(result)
		if l > 0 {
			lastOpType := resultTypes[l-1]
			cancel1 := (lastOpType == opLocal || lastOpType == opRedo) && opType == opUndo
			cancel2 := lastOpType == opUndo && opType == opRedo
			if cancel1 || cancel2 {
				result = result[0 : l-1]
				resultTypes = resultTypes[0 : l-1]
				continue
			}
		}
		resultTypes = append(resultTypes, opType)
		result = append(result, ops[kk])
	}
	return result
}
