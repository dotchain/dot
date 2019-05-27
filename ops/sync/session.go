// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package sync

import (
	"context"
	"strconv"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
)

type session struct {
	config *Config
	stream streams.Stream
	out    []ops.Op
}

func (s *session) push() error {
	stream, c := streams.Latest(s.stream)
	s.stream = stream
	err := s.appendChange(c)

	if len(s.out) > 0 && err == nil {
		err = s.config.Store.Append(context.Background(), s.out)
		if err == nil {
			s.out = nil
		}
	}

	return err
}

func (s *session) appendChange(c changes.Change) error {
	if c == nil {
		return nil
	}

	cfg := s.config
	id, err := s.newID()
	if err == nil {
		op := ops.Operation{OpID: id, BasisID: cfg.Version, Change: c}
		if len(cfg.Pending) > 0 {
			op.ParentID = cfg.Pending[len(cfg.Pending)-1].ID()
		}
		cfg.Pending = append(cfg.Pending, op)
		cfg.MergeChain = append(cfg.MergeChain, op)
		s.out = append(s.out, op)
		cfg.Notify(cfg.Version, cfg.Pending, cfg.MergeChain)
	}
	return err
}

func (s *session) pull() error {
	cfg := s.config
	version := cfg.Version

	ops, err := cfg.Store.GetSince(context.Background(), version+1, 1000)
	if err != nil {
		return err
	}

	for _, op := range ops {
		if op.Version() != cfg.Version+1 {
			return verMismatchError{op.Version(), cfg.Version + 1}
		}

		if len(cfg.Pending) > 0 && cfg.Pending[0].ID() == op.ID() {
			cfg.Pending = cfg.Pending[1:]
			cfg.MergeChain = cfg.MergeChain[1:]
		} else {
			for idx, pending := range cfg.MergeChain {
				pc, oc := changes.Merge(op.Changes(), pending.Changes())
				cfg.MergeChain[idx] = pending.WithChanges(pc)
				op = op.WithChanges(oc)
			}
			s.stream = s.stream.ReverseAppend(op.Changes())
		}
		cfg.Version++
	}

	if cfg.Version > version {
		cfg.Notify(cfg.Version, cfg.Pending, cfg.MergeChain)
	}
	return nil
}

type verMismatchError struct {
	got, expected int
}

func (v verMismatchError) Error() string {
	exp, got := strconv.Itoa(v.expected), strconv.Itoa(v.got)
	return "version mismatched: " + got + " (expected " + exp + ")"
}
