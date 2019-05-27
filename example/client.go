// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build ignore

package main

import (
	"fmt"
	"time"

	"github.com/dotchain/dot/example"
)

func main() {
	var stream *example.TodoListStream

	count := 1
	go func() {
		for {
			time.Sleep(time.Second * 2)
			todo := example.Todo{Description: fmt.Sprintf("Heyaaa %d", count)}
			if stream != nil {
				example.Lock.Lock()
				example.AddTodo(stream.Latest(), todo)
				stream = stream.Latest()
				example.Lock.Unlock()
				count++
			}
		}
	}()
	example.Client(nil, func(s *example.TodoListStream) {
		stream = s
		for kk, todo := range s.Value {
			fmt.Printf("%d %s\n", kk, todo.Description)
		}
	})
}
