// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

//+build stress

package stress

//go:generate go run codegen.go

import (
	"math/rand"
	"sync"

	"github.com/dotchain/dot"
	"github.com/dotchain/dot/ops"
)

// SessionState is the state associated with a previous session
type SessionState struct {
	State State
	*dot.Session
}

// Reconnnect creates a new session from this state
func (ss SessionState) Reconnect(serverUrl string, numClients int, wg *sync.WaitGroup) *Session {
	if ss.Session == nil {
		ss.Session = dot.NewSession()
	}
	s, store := ss.Session.Stream("http://localhost:8083/stress/", nil)
	stateStream := &StateStream{Stream: s, Value: ss.State}
	countStream := stateStream.Count()
	last := int32(countStream.Value) / int32(numClients)
	return &Session{stateStream, ss.Session, store, numClients, wg, last}
}

// Session represents a single session
type Session struct {
	*StateStream
	*dot.Session
	ops.Store
	numClients int
	wg         *sync.WaitGroup
	last       int32
}

// NewSession creates a new session
func NewSession(serverUrl string, numClients int, wg *sync.WaitGroup) *Session {
	return (SessionState{}).Reconnect(serverUrl, numClients, wg)
}

// Close releases all resources
func (s *Session) Close() SessionState {
	s.Store.Close()
	return SessionState{s.StateStream.Latest().Value, s.Session}
}

// MakeSomeRandomChanges does exactly that but also increments the count
func (s *Session) MakeSomeRandomChanges(iterations int) {
	go func() {
		stream := s.StateStream.Latest()
		for kk := 0; kk < iterations; kk++ {
			l := len(stream.Value.Text)
			insert := s.randString(3)

			if l == 0 {
				stream.Text().Splice(0, 0, insert)
			} else {
				var offset, count int
				if l > 0 {
					offset = rand.Intn(l)
				}
				if l-offset > 0 {
					count = rand.Intn(l - offset)
				}
				stream.Text().Splice(offset, count, insert)
			}
		}

		stream.Count().Increment(1)
		if err := stream.Stream.Push(); err != nil {
			panic(err)
		}
		stream = stream.Latest()
		current := int32(stream.Value.Count) / int32(s.numClients)
		for current == s.last {
			if err := stream.Stream.Pull(); err != nil {
				panic(err)
			}

			stream = stream.Latest()
			current = int32(stream.Value.Count) / int32(s.numClients)
		}
		s.wg.Add(int(s.last - current))
		s.last = current
	}()
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *Session) randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
