// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"errors"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/x/idgen"
	"testing"
)

func TestAppendFailurePanic(t *testing.T) {
	expectedError := errors.New("my error")
	defer func() {
		if r := recover(); r != expectedError {
			t.Fatal("Failed to panic or panic'ed differently", r)
		}
	}()

	store := &failAppend{err: expectedError}
	client := streams.New()
	s := ops.NewSync(store, -1, client, idgen.New)
	defer s.Close()
	client.Append(changes.Move{1, 2, 3})
}

type failAppend struct {
	err error
	ops.Store
}

func (f *failAppend) Append(ctx context.Context, opx []ops.Op) error {
	return f.err
}
