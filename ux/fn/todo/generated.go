// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.
//

package todo

import (
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/ux/core"
	"github.com/dotchain/dot/ux/fn"
	uxstreams "github.com/dotchain/dot/ux/streams"
)

// TaskStream is a stream of Task values.
type TaskStream struct {
	// Notifier provides On/Off/Notify support. New instances of
	// TaskStream created via the AppendLocal or AppendRemote
	// share the same Notifier value.
	*streams.Notifier

	// Value holds the current value. The latest value may be
	// fetched via the Latest() method.
	Value Task

	// Change tracks the change that leads to the next value.
	Change changes.Change

	// Next tracks the next value in the stream.
	Next *TaskStream
}

// NewTaskStream creates a new Task stream
func NewTaskStream(value Task) *TaskStream {
	return &TaskStream{&streams.Notifier{}, value, nil, nil}
}

// Latest returns the latest value in the stream
func (s *TaskStream) Latest() *TaskStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Append appends a local change. isLocal identifies if the caller is
// local or remote. It returns the updated stream whose value matches
// the provided value and whose Latest() converges to the latest of
// the stream.
func (s *TaskStream) Append(c changes.Change, value Task, isLocal bool) *TaskStream {
	if c == nil {
		c = changes.Replace{Before: s.wrapValue(s.Value), After: s.wrapValue(value)}
	}

	// return value: after is correctly set to provided value
	result := &TaskStream{Notifier: s.Notifier, Value: value}

	// before tracks s, after tracks result, v tracks latest value
	// of after chain
	before := s
	var v changes.Value = changes.Atomic{value}

	// walk the chain of Next and find corresponding values to
	// add to after so that both s annd after converge
	after := result
	for ; before.Next != nil; before = before.Next {
		var afterChange changes.Change

		if isLocal {
			c, afterChange = before.Change.Merge(c)
		} else {
			afterChange, c = c.Merge(before.Change)
		}

		if c == nil {
			// the convergence point is before.Next
			after.Change, after.Next = afterChange, before.Next
			return result
		}

		if afterChange == nil {
			continue
		}

		// append this to after and continue with that
		v = v.Apply(nil, afterChange)
		after.Change = afterChange
		after.Next = &TaskStream{Notifier: s.Notifier, Value: s.unwrapValue(v)}
		after = after.Next
	}

	// append the residual change (c) to converge to wherever
	// after has landed. Notify since s.Latest() has now changed
	before.Change, before.Next = c, after
	s.Notify()
	return result
}

func (s *TaskStream) wrapValue(i interface{}) changes.Value {
	if x, ok := i.(changes.Value); ok {
		return x
	}
	return changes.Atomic{i}
}

func (s *TaskStream) unwrapValue(v changes.Value) Task {
	if x, ok := v.(interface{}).(Task); ok {
		return x
	}
	return v.(changes.Atomic).Value.(Task)
}

// SetDone updates the field with a new value
func (s *TaskStream) SetDone(v bool) *TaskStream {
	c := changes.Replace{s.wrapValue(s.Value.Done), s.wrapValue(v)}
	value := s.Value
	value.Done = v
	key := []interface{}{"Done"}
	return s.Append(changes.PathChange{key, c}, value, true)
}

// DoneSubstream returns a stream for Done that is automatically
// connected to the current TaskStream instance.  Updates to one
// automatically update the other.
func (s *TaskStream) DoneSubstream(cache streams.Cache) (field *uxstreams.BoolStream) {
	n := s.Notifier
	handler := &streams.Handler{nil}
	if f, h, ok := cache.GetSubstream(n, "Done"); ok {
		field, handler = f.(*uxstreams.BoolStream), h
	} else {
		field = uxstreams.NewBoolStream(s.Value.Done)
		parent, merging, path := s, false, []interface{}{"Done"}
		handler.Handle = func() {
			if merging {
				return
			}

			merging = true
			for ; field.Next != nil; field = field.Next {
				v := parent.Value
				v.Done = field.Next.Value
				c := changes.PathChange{path, field.Change}
				parent = parent.Append(c, v, true)
			}

			for ; parent.Next != nil; parent = parent.Next {
				result := refs.Merge(path, parent.Change)
				if result == nil {
					field = field.Append(nil, parent.Next.Value.Done, true)
				} else {
					field = field.Append(result.Affected, parent.Next.Value.Done, true)
				}
			}
			merging = false
		}
		field.On(handler)
		parent.On(handler)
	}

	handler.Handle()
	field = field.Latest()
	n2 := field.Notifier
	close := func() { n.Off(handler); n2.Off(handler) }
	cache.SetSubstream(n, "Done", field, handler, close)
	return field
}

// SetDescription updates the field with a new value
func (s *TaskStream) SetDescription(v string) *TaskStream {
	c := changes.Replace{s.wrapValue(s.Value.Description), s.wrapValue(v)}
	value := s.Value
	value.Description = v
	key := []interface{}{"Description"}
	return s.Append(changes.PathChange{key, c}, value, true)
}

// DescriptionSubstream returns a stream for Description that is automatically
// connected to the current TaskStream instance.  Updates to one
// automatically update the other.
func (s *TaskStream) DescriptionSubstream(cache streams.Cache) (field *uxstreams.TextStream) {
	n := s.Notifier
	handler := &streams.Handler{nil}
	if f, h, ok := cache.GetSubstream(n, "Description"); ok {
		field, handler = f.(*uxstreams.TextStream), h
	} else {
		field = uxstreams.NewTextStream(s.Value.Description)
		parent, merging, path := s, false, []interface{}{"Description"}
		handler.Handle = func() {
			if merging {
				return
			}

			merging = true
			for ; field.Next != nil; field = field.Next {
				v := parent.Value
				v.Description = field.Next.Value
				c := changes.PathChange{path, field.Change}
				parent = parent.Append(c, v, true)
			}

			for ; parent.Next != nil; parent = parent.Next {
				result := refs.Merge(path, parent.Change)
				if result == nil {
					field = field.Append(nil, parent.Next.Value.Description, true)
				} else {
					field = field.Append(result.Affected, parent.Next.Value.Description, true)
				}
			}
			merging = false
		}
		field.On(handler)
		parent.On(handler)
	}

	handler.Handle()
	field = field.Latest()
	n2 := field.Notifier
	close := func() { n.Off(handler); n2.Off(handler) }
	cache.SetSubstream(n, "Description", field, handler, close)
	return field
}

// TasksStream is a stream of Tasks values.
type TasksStream struct {
	// Notifier provides On/Off/Notify support. New instances of
	// TasksStream created via the AppendLocal or AppendRemote
	// share the same Notifier value.
	*streams.Notifier

	// Value holds the current value. The latest value may be
	// fetched via the Latest() method.
	Value Tasks

	// Change tracks the change that leads to the next value.
	Change changes.Change

	// Next tracks the next value in the stream.
	Next *TasksStream
}

// NewTasksStream creates a new Tasks stream
func NewTasksStream(value Tasks) *TasksStream {
	return &TasksStream{&streams.Notifier{}, value, nil, nil}
}

// Latest returns the latest value in the stream
func (s *TasksStream) Latest() *TasksStream {
	for s.Next != nil {
		s = s.Next
	}
	return s
}

// Append appends a local change. isLocal identifies if the caller is
// local or remote. It returns the updated stream whose value matches
// the provided value and whose Latest() converges to the latest of
// the stream.
func (s *TasksStream) Append(c changes.Change, value Tasks, isLocal bool) *TasksStream {
	if c == nil {
		c = changes.Replace{Before: s.wrapValue(s.Value), After: s.wrapValue(value)}
	}

	// return value: after is correctly set to provided value
	result := &TasksStream{Notifier: s.Notifier, Value: value}

	// before tracks s, after tracks result, v tracks latest value
	// of after chain
	before := s
	var v changes.Value = changes.Atomic{value}

	// walk the chain of Next and find corresponding values to
	// add to after so that both s annd after converge
	after := result
	for ; before.Next != nil; before = before.Next {
		var afterChange changes.Change

		if isLocal {
			c, afterChange = before.Change.Merge(c)
		} else {
			afterChange, c = c.Merge(before.Change)
		}

		if c == nil {
			// the convergence point is before.Next
			after.Change, after.Next = afterChange, before.Next
			return result
		}

		if afterChange == nil {
			continue
		}

		// append this to after and continue with that
		v = v.Apply(nil, afterChange)
		after.Change = afterChange
		after.Next = &TasksStream{Notifier: s.Notifier, Value: s.unwrapValue(v)}
		after = after.Next
	}

	// append the residual change (c) to converge to wherever
	// after has landed. Notify since s.Latest() has now changed
	before.Change, before.Next = c, after
	s.Notify()
	return result
}

func (s *TasksStream) wrapValue(i interface{}) changes.Value {
	if x, ok := i.(changes.Value); ok {
		return x
	}
	return changes.Atomic{i}
}

func (s *TasksStream) unwrapValue(v changes.Value) Tasks {
	if x, ok := v.(interface{}).(Tasks); ok {
		return x
	}
	return v.(changes.Atomic).Value.(Tasks)
}

// Substream returns a stream for the specified index that is
// automatically connected to the current TasksStream instance.  Updates to
// one automatically update the other.
func (s *TasksStream) Substream(cache streams.Cache, index int) (entry *TaskStream) {
	n := s.Notifier
	handler := &streams.Handler{nil}
	if f, h, ok := cache.GetSubstream(n, index); ok {
		entry, handler = f.(*TaskStream), h
	} else {
		entry = NewTaskStream(s.Value[index])
		parent, merging, path := s, false, []interface{}{index}
		handler.Handle = func() {
			if merging {
				return
			}

			merging = true
			for ; entry.Next != nil; entry = entry.Next {
				v := append(Tasks(nil), parent.Value...)
				v[index] = entry.Next.Value
				c := changes.PathChange{path, entry.Change}
				parent = parent.Append(c, v, true)
			}

			for ; parent.Next != nil; parent = parent.Next {
				result := refs.Merge(path, parent.Change)
				var c changes.Change
				if result != nil {
					index = result.P[0].(int)
					// TODO: if the index changed fix up
					// the key in the cache
					c = result.Affected
				}
				entry = entry.Append(c, parent.Next.Value[index], true)
			}
			merging = false
		}
		entry.On(handler)
		parent.On(handler)
	}

	handler.Handle()
	entry = entry.Latest()
	n2 := entry.Notifier
	close := func() { n.Off(handler); n2.Off(handler) }
	cache.SetSubstream(n, index, entry, handler, close)
	return entry
}

type taskEditCtx struct {
	streams.Cache

	initialized bool

	fn struct {
		fn.CheckboxCache
		fn.ElementCache
		fn.TextEditCache
	}
	memoized struct {
		result1 core.Element
		styles  core.Styles
		task    *TaskStream
	}
}

func (ctx *taskEditCtx) areArgsSame(styles core.Styles, task *TaskStream) bool {

	if styles != ctx.memoized.styles {
		return false
	}

	return task == ctx.memoized.task

}

func (ctx *taskEditCtx) refreshIfNeeded(styles core.Styles, task *TaskStream) (result1 core.Element) {
	if !ctx.initialized || !ctx.areArgsSame(styles, task) {
		return ctx.refresh(styles, task)
	}
	return ctx.memoized.result1
}

func (ctx *taskEditCtx) refresh(styles core.Styles, task *TaskStream) (result1 core.Element) {
	ctx.initialized = true
	ctx.memoized.styles, ctx.memoized.task = styles, task

	ctx.Cache.Begin()
	defer ctx.Cache.End()
	ctx.fn.CheckboxCache.Begin()
	defer ctx.fn.CheckboxCache.End()

	ctx.fn.ElementCache.Begin()
	defer ctx.fn.ElementCache.End()

	ctx.fn.TextEditCache.Begin()
	defer ctx.fn.TextEditCache.End()
	ctx.memoized.result1 = TaskEdit(ctx, styles, task)
	return ctx.memoized.result1
}

func (ctx *taskEditCtx) close() {
	ctx.Cache.Begin()
	defer ctx.Cache.End()
	ctx.fn.CheckboxCache.Begin()
	defer ctx.fn.CheckboxCache.End()

	ctx.fn.ElementCache.Begin()
	defer ctx.fn.ElementCache.End()

	ctx.fn.TextEditCache.Begin()
	defer ctx.fn.TextEditCache.End()
}

// TaskEditCache is good
type TaskEditCache struct {
	old, current map[interface{}]*taskEditCtx
}

// Begin starts a round
func (c *TaskEditCache) Begin() {
	c.old, c.current = c.current, map[interface{}]*taskEditCtx{}
}

// End finishes the round cleaning up any unused components
func (c *TaskEditCache) End() {
	for _, ctx := range c.old {
		ctx.close()
	}
	c.old = nil
}

// TaskEdit is good
func (c *TaskEditCache) TaskEdit(ctxKey interface{}, styles core.Styles, task *TaskStream) (result1 core.Element) {
	ctxOld, ok := c.old[ctxKey]
	if ok {
		delete(c.old, ctxKey)
	} else {
		ctxOld = &taskEditCtx{}
	}
	c.current[ctxKey] = ctxOld
	return ctxOld.refreshIfNeeded(styles, task)
}

type tasksViewCtx struct {
	streams.Cache

	TaskEditCache
	initialized bool

	fn struct {
		fn.ElementCache
	}
	memoized struct {
		result1     core.Element
		showDone    *uxstreams.BoolStream
		showNotDone *uxstreams.BoolStream
		styles      core.Styles
		tasks       *TasksStream
	}
}

func (ctx *tasksViewCtx) areArgsSame(styles core.Styles, showDone *uxstreams.BoolStream, showNotDone *uxstreams.BoolStream, tasks *TasksStream) bool {

	if styles != ctx.memoized.styles {
		return false
	}

	if showDone != ctx.memoized.showDone {
		return false
	}

	if showNotDone != ctx.memoized.showNotDone {
		return false
	}

	return tasks == ctx.memoized.tasks

}

func (ctx *tasksViewCtx) refreshIfNeeded(styles core.Styles, showDone *uxstreams.BoolStream, showNotDone *uxstreams.BoolStream, tasks *TasksStream) (result1 core.Element) {
	if !ctx.initialized || !ctx.areArgsSame(styles, showDone, showNotDone, tasks) {
		return ctx.refresh(styles, showDone, showNotDone, tasks)
	}
	return ctx.memoized.result1
}

func (ctx *tasksViewCtx) refresh(styles core.Styles, showDone *uxstreams.BoolStream, showNotDone *uxstreams.BoolStream, tasks *TasksStream) (result1 core.Element) {
	ctx.initialized = true
	ctx.memoized.styles, ctx.memoized.showDone, ctx.memoized.showNotDone, ctx.memoized.tasks = styles, showDone, showNotDone, tasks

	ctx.Cache.Begin()
	defer ctx.Cache.End()
	ctx.TaskEditCache.Begin()
	defer ctx.TaskEditCache.End()

	ctx.fn.ElementCache.Begin()
	defer ctx.fn.ElementCache.End()
	ctx.memoized.result1 = TasksView(ctx, styles, showDone, showNotDone, tasks)
	return ctx.memoized.result1
}

func (ctx *tasksViewCtx) close() {
	ctx.Cache.Begin()
	defer ctx.Cache.End()
	ctx.TaskEditCache.Begin()
	defer ctx.TaskEditCache.End()

	ctx.fn.ElementCache.Begin()
	defer ctx.fn.ElementCache.End()
}

// TasksViewCache is good
type TasksViewCache struct {
	old, current map[interface{}]*tasksViewCtx
}

// Begin starts a round
func (c *TasksViewCache) Begin() {
	c.old, c.current = c.current, map[interface{}]*tasksViewCtx{}
}

// End finishes the round cleaning up any unused components
func (c *TasksViewCache) End() {
	for _, ctx := range c.old {
		ctx.close()
	}
	c.old = nil
}

// TasksView is good
func (c *TasksViewCache) TasksView(ctxKey interface{}, styles core.Styles, showDone *uxstreams.BoolStream, showNotDone *uxstreams.BoolStream, tasks *TasksStream) (result1 core.Element) {
	ctxOld, ok := c.old[ctxKey]
	if ok {
		delete(c.old, ctxKey)
	} else {
		ctxOld = &tasksViewCtx{}
	}
	c.current[ctxKey] = ctxOld
	return ctxOld.refreshIfNeeded(styles, showDone, showNotDone, tasks)
}

type appCtx struct {
	streams.Cache

	TasksViewCache
	initialized bool

	fn struct {
		fn.CheckboxCache
		fn.ElementCache
	}
	memoized struct {
		result1 core.Element
		styles  core.Styles
		tasks   *TasksStream
	}
}

func (ctx *appCtx) areArgsSame(styles core.Styles, tasks *TasksStream) bool {

	if styles != ctx.memoized.styles {
		return false
	}

	return tasks == ctx.memoized.tasks

}

func (ctx *appCtx) refreshIfNeeded(styles core.Styles, tasks *TasksStream) (result1 core.Element) {
	if !ctx.initialized || !ctx.areArgsSame(styles, tasks) {
		return ctx.refresh(styles, tasks)
	}
	return ctx.memoized.result1
}

func (ctx *appCtx) refresh(styles core.Styles, tasks *TasksStream) (result1 core.Element) {
	ctx.initialized = true
	ctx.memoized.styles, ctx.memoized.tasks = styles, tasks

	ctx.Cache.Begin()
	defer ctx.Cache.End()
	ctx.TasksViewCache.Begin()
	defer ctx.TasksViewCache.End()

	ctx.fn.ElementCache.Begin()
	defer ctx.fn.ElementCache.End()

	ctx.fn.CheckboxCache.Begin()
	defer ctx.fn.CheckboxCache.End()
	ctx.memoized.result1 = App(ctx, styles, tasks)
	return ctx.memoized.result1
}

func (ctx *appCtx) close() {
	ctx.Cache.Begin()
	defer ctx.Cache.End()
	ctx.TasksViewCache.Begin()
	defer ctx.TasksViewCache.End()

	ctx.fn.ElementCache.Begin()
	defer ctx.fn.ElementCache.End()

	ctx.fn.CheckboxCache.Begin()
	defer ctx.fn.CheckboxCache.End()
}

// AppCache is good
type AppCache struct {
	old, current map[interface{}]*appCtx
}

// Begin starts a round
func (c *AppCache) Begin() {
	c.old, c.current = c.current, map[interface{}]*appCtx{}
}

// End finishes the round cleaning up any unused components
func (c *AppCache) End() {
	for _, ctx := range c.old {
		ctx.close()
	}
	c.old = nil
}

// App is good
func (c *AppCache) App(ctxKey interface{}, styles core.Styles, tasks *TasksStream) (result1 core.Element) {
	ctxOld, ok := c.old[ctxKey]
	if ok {
		delete(c.old, ctxKey)
	} else {
		ctxOld = &appCtx{}
	}
	c.current[ctxKey] = ctxOld
	return ctxOld.refreshIfNeeded(styles, tasks)
}
