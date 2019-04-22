// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
)

// Handler implements ServerHTTP using the provided store and codecs
// map. If no codecs map is provided, DefaultCodecs is used instead.
type Handler struct {
	ops.Store
	Codecs map[string]Codec
	log.Log

	once sync.Once
}

// ServeHTTP uses the code to unmarshal a request, apply it and then
// encode back the response
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		if h.Log == nil {
			h.Log = log.Default()
		}
	})

	defer func() {
		h.report("unexpected client close", r.Body.Close())
	}()

	ct := r.Header.Get("Content-Type")

	codecs := h.Codecs
	if codecs == nil {
		codecs = DefaultCodecs
	}

	codec := codecs[ct]
	if codec == nil {
		h.Log.Println("Client used an unknown type", ct)
		http.Error(w, "Invalid content-type", 400)
		return
	}

	var req request
	err := codec.Decode(&req, r.Body)
	if err != nil {
		http.Error(w, h.codecError(err).Error(), 400)
		return
	}

	duration := 30 * time.Second
	if req.Duration != 0 {
		duration = req.Duration
	}

	ctx, done := context.WithTimeout(r.Context(), duration)
	defer done()

	var res response
	res.Error = errors.New("unknown error")
	switch req.Name {
	case "Append":
		res.Error = h.Append(ctx, req.Ops)
	case "GetSince":
		res.Ops, res.Error = h.GetSince(ctx, req.Version, req.Limit)
	}

	// do this hack since we can't be sure what error types are possible
	h.patchResponseError(ctx, &res)

	var buf bytes.Buffer
	if err := codec.Encode(res, &buf); err != nil {
		http.Error(w, h.codecError(err).Error(), 400)
		return
	}

	w.Header().Add("Content-Type", ct)
	_, err = w.Write(buf.Bytes())
	h.report("Unexpected write error", err)
}

func (h *Handler) codecError(err error) error {
	h.Log.Println("Codec error (see https://github.com/dotchain/dot/wiki/Gob-error)")
	h.Log.Println(err)
	return err
}

func (h *Handler) report(msg string, err error) {
	if err != nil {
		h.Log.Println(msg, err)
	}
}

func (h *Handler) patchResponseError(ctx context.Context, res *response) {
	// error types are often not registered with encoding/gob and
	// so fail to get encoded. simply convert them to string errors
	if res.Error != nil {
		if res.Error != ctx.Err() {
			h.Log.Println("failed", res.Error)
		}
		res.Error = strError(res.Error.Error())
	}
}
