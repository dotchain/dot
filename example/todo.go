// Code generated by github.com/tvastar/test/cmd/testmd/testmd.go. DO NOT EDIT.

package example

import (
	"encoding/gob"
	"math/rand"
	"net/http"

	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/ops/bolt"
	"github.com/dotchain/dot/ops/nw"
)

func Server() {
	// import net/http
	// import github.com/dotchain/dot/ops/nw
	// import github.com/dotchain/dot/ops/bolt

	// uses a local-file backed bolt DB backend
	store, _ := bolt.New("file.bolt", "instance", nil)
	store = nw.MemPoller(store)
	defer store.Close()
	http.Handle("/api/", &nw.Handler{Store: store})
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

// import github.com/dotchain/dot/ops/nw
// import github.com/dotchain/dot/ops
// import math/rand

func Client(stop chan struct{}, render func(*TodoListStream)) {
	version, pending, todos := SavedSession()

	store := &nw.Client{URL: "http://localhost:8080/api/"}
	defer store.Close()
	client := ops.NewConnector(version, pending, ops.Transformed(store), rand.Float64)
	stream := &TodoListStream{Stream: client.Stream, Value: todos}

	// start the network processing
	client.Connect()

	// save session before shutdown
	defer func() {
		SaveSession(client.Version, client.Pending, stream.Latest().Value)
	}()
	defer client.Disconnect()

	client.Stream.Nextf("key", func() {
		stream = stream.Latest()
		render(stream)
	})
	render(stream)
	defer func() {
		client.Stream.Nextf("key", nil)
	}()

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
