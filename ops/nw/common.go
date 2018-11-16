// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"encoding/gob"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/changes/types"
	"io"
	"time"
)

// Codec is the interface codecs will have to implement to marshal and
// unmarshal requests and responses
type Codec interface {
	Encode(value interface{}, writer io.Writer) error
	Decode(value interface{}, reader io.Reader) error
}

// DefaultCodecs is the default codecs list which contains a map of
// content-type to codec.
var DefaultCodecs = map[string]Codec{
	"application/x-gob": gobCodec{},
}

type gobCodec struct{}

func (c gobCodec) Encode(value interface{}, writer io.Writer) error {
	return gob.NewEncoder(writer).Encode(value)
}

func (c gobCodec) Decode(value interface{}, reader io.Reader) error {
	return gob.NewDecoder(reader).Decode(value)
}

type request struct {
	Name           string
	Ops            []ops.Op
	Version, Limit int
	Duration       time.Duration
}

type response struct {
	Ops   []ops.Op
	Error error
}

var standardTypes = []interface{}{
	changes.Replace{},
	changes.Move{},
	changes.Splice{},
	changes.PathChange{},
	changes.ChangeSet{},
	changes.Atomic{},
	changes.Nil,
	types.A{},
	types.S8(""),
	types.S16(""),
	types.M{},
	ops.Operation{},
	refs.Update{},
	refs.Range{},
	refs.Path{},
	refs.Caret{},
}

func init() {
	for _, typ := range standardTypes {
		gob.Register(typ)
	}
	gob.Register(strError(""))
}

type strError string

func (s strError) Error() string {
	return string(s)
}
