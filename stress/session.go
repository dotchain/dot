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
	State   State
	Version int
	Pending []ops.Op
}

// Reconnnect creates a new session from this state
func (ss SessionState) Reconnect(serverUrl string, numClients int, wg *sync.WaitGroup) *Session {
	session, s := dot.Reconnect("http://localhost:8083/stress/", ss.Version, ss.Pending)
	result := &Session{
		Session:     session,
		StateStream: &StateStream{Stream: s, Value: ss.State},
		scheduler:   s.(scheduler),
	}
	last := int32(ss.State.Count) / int32(numClients)
	s.Nextf(session, func() {
		result.StateStream = result.StateStream.Latest()
		current := int32(result.StateStream.Value.Count) / int32(numClients)
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
	scheduler
}

// NewSession creates a new session
func NewSession(serverUrl string, numClients int, wg *sync.WaitGroup) *Session {
	session, s := dot.Connect("http://localhost:8083/stress/")
	result := &Session{
		Session:     session,
		StateStream: &StateStream{Stream: s},
		scheduler:   s.(scheduler),
	}
	var last int32
	s.Nextf(session, func() {
		result.StateStream = result.StateStream.Latest()
		current := int32(result.StateStream.Value.Count) / int32(numClients)
		if current > last {
			wg.Add(int(last - current))
		}
		last = current
	})
	return result
}

// Close releases all resources
func (s *Session) Close() SessionState {
	var ss SessionState
	closed := make(chan interface{}, 1)
	s.scheduler.Schedule(func() {
		s.StateStream.Stream.Nextf(s.Session, nil)
		ss.Version, ss.Pending = s.Session.Close()
		ss.State = s.StateStream.Latest().Value
		closed <- nil
	})
	<-closed
	return ss
}

// MakeSomeRandomChanges does exactly that but also increments the count
func (s *Session) MakeSomeRandomChanges(iterations int) {
	go s.scheduler.Schedule(func() {
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
	})
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *Session) randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type scheduler interface {
	Schedule(fn func())
}
