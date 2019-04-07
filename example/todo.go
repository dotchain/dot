// Code generated by github.com/tvastar/test/cmd/testmd/testmd.go. DO NOT EDIT.

package example

import (
	"encoding/gob"
	"net/http"

	"github.com/dotchain/dot"
	"github.com/dotchain/dot/ops"
)

func Server() {
	// import net/http
	// import github.com/dotchain/dot

	// uses a local-file backed bolt DB backend
	http.Handle("/dot/", dot.BoltServer("file.bolt"))
	http.ListenAndServe(":8080", nil)
}

// Todo tracks a single todo item
type Todo struct {
	Complete    bool
	Description string
}

// TodoList tracks a collection of todo items
type TodoList []Todo

// import encoding/gob

func init() {
	gob.Register(Todo{})
	gob.Register(TodoList{})
}
func Toggle(t *TodoListStream, index int) {
	// TodoListStream.Item() is generated code. It returns
	// a stream of the n'th element of the slice so that
	// particular stream can be modified. When that stream is
	// modified, the effect is automatically merged into the
	// parent (and available via .Next of the parent stream)
	todoStream := t.Item(index)

	// TodoStream.Complete is generated code. It returns a stream
	// for the Todo.Complete field so that it can be modified. As
	// with slices above, mutations on the field's stream are
	// reflected on the struct stream (via .Next or .Latest())
	completeStream := todoStream.Complete()

	// completeStream is of type streams.Bool. All streams
	// implement the simple Update(newValue) method that replaces
	// the current value with a new value.
	completeStream.Update(!completeStream.Value)
}
func SpliceDescription(t *TodoListStream, index, offset, count int, replacement string) {
	// TodoListStream.Item() is generated code. It returns
	// a stream of the n'th element of the slice so that
	// particular stream can be modified. When that stream is
	// modified, the effect is automatically merged into the
	// parent (and available via .Next of the parent stream)
	todoStream := t.Item(index)

	// TodoStream.Description is generated code. It returns a
	// stream for the Todo.Description field so that it can be
	// modified. As with slices above, mutations on the field's
	// stream are reflected on the struct stream (via .Next or
	// .Latest())
	// TodoStream.Description() returns streams.S16 type
	descStream := todoStream.Description()

	// streams.S16 implements Splice(offset, removeCount, replacement)
	descStream.Splice(offset, count, replacement)
}
func AddTodo(t *TodoListStream, todo Todo) {
	// All slice streams implement Splice(offset, removeCount, replacement)
	t.Splice(len(t.Value), 0, todo)
}

// import github.com/dotchain/dot

func Client(stop chan struct{}, render func(*TodoListStream)) {
	url := "http://localhost:8080/dot/"
	version, pending, todos := SavedSession()

	session, s, _ := dot.Reconnect(url, version, pending)
	todosStream := &TodoListStream{Stream: s, Value: todos}

	// save session before shutdown
	defer func() {
		todosStream.Stream.Nextf("key", nil)
		version, pending = session.Close()
		todos = todosStream.Latest().Value
		SaveSession(version, pending, todos)
	}()

	render(todosStream)
	todosStream.Stream.Nextf("key", func() {
		todosStream = todosStream.Latest()
		render(todosStream)
	})
	<-stop
}

func SaveSession(version int, pending []ops.Op, todos TodoList) {
	// this is not yet implemented. if it were, then
	// this value should be persisted locally and returned
	// by the call to savedSession
}

func SavedSession() (version int, pending []ops.Op, todos TodoList) {
	// this is not yet implemented. return default values
	return -1, nil, nil
}
