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
	// Version here is 1 + actual version to make zero value of
	// session state refer to non prior state
	Version int
	Pending []ops.Op
}

// Reconnnect creates a new session from this state
func (ss SessionState) Reconnect(serverUrl string, numClients int, wg *sync.WaitGroup) *Session {
	session, s, _ := dot.Reconnect("http://localhost:8083/stress/", ss.Version-1, ss.Pending)
	stateStream := &StateStream{Stream: s, Value: ss.State}
	countStream := stateStream.Count()
	result := &Session{stateStream, session}

	last := int32(countStream.Value) / int32(numClients)
	var l sync.Mutex
	s.Nextf(session, func() {
		l.Lock()
		defer l.Unlock()
		countStream = countStream.Latest()
		current := int32(countStream.Value) / int32(numClients)
		if current > last {
			wg.Add(int(last - current))
		}
		last = current
	})
	return result

}

// Session represents a single session
type Session struct {
	*StateStream
	*dot.Session
}

// NewSession creates a new session
func NewSession(serverUrl string, numClients int, wg *sync.WaitGroup) *Session {
	return (SessionState{}).Reconnect(serverUrl, numClients, wg)
}

// Close releases all resources
func (s *Session) Close() SessionState {
	var ss SessionState
	ss.Version, ss.Pending = s.Session.Close()
	s.StateStream.Stream.Nextf(s.Session, nil)
	// Verision is 1 + actual version so that zero value
	// corresponds to no previous state
	ss.Version++
	ss.State = s.StateStream.Latest().Value
	return ss
}

// MakeSomeRandomChanges does exactly that but also increments the count
func (s *Session) MakeSomeRandomChanges(iterations int) {
	go func() {
		stream := s.StateStream.Latest().Text()
		defer s.StateStream.Count().Increment(1)

		for kk := 0; kk < iterations; kk++ {
			l := len(stream.Value)
			insert := s.randString(3)

			if l == 0 {
				stream = stream.Splice(0, 0, insert)
			} else {
				var offset, count int
				if l > 0 {
					offset = rand.Intn(l)
				}
				if l-offset > 0 {
					count = rand.Intn(l - offset)
				}
				stream = stream.Splice(offset, count, insert)
			}
		}
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
