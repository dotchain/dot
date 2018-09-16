// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"context"
	"errors"
	"github.com/dotchain/dot/ops"
	"net/http"
	"time"
)

// Handler implements ServerHTTP using the provided store and codecs
// map. If no codecs map is provided, DefaultCodecs is used instead.
type Handler struct {
	ops.Store
	Codecs map[string]Codec
}

// ServeHTTP uses the code to unmarshal a request, apply it and then
// encode back the response
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	codecs := h.Codecs
	if codecs == nil {
		codecs = DefaultCodecs
	}

	var req request
	ct := r.Header.Get("Content-Type")
	codec := codecs[ct]
	if codec == nil {
		http.Error(w, "Invalid content-type", 400)
		return
	}

	err := codec.Decode(&req, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	duration := req.Duration
	if duration == 0 {
		duration = 5 * time.Second
	}
	ctx, done := context.WithDeadline(context.Background(), time.Now().Add(duration))
	defer done()

	var res response
	res.Error = errors.New("Unknown error")
	switch req.Name {
	case "Append":
		res.Error = h.Append(ctx, req.Ops)
	case "GetSince":
		res.Ops, res.Error = h.GetSince(ctx, req.Version, req.Limit)
	case "Poll":
		res.Error = h.Poll(ctx, req.Version)
	}

	// do this hack since we can't be sure what error types are possible
	if res.Error != nil {
		res.Error = strError(res.Error.Error())
	}

	w.Header().Add("Content-Type", ct)
	if err := codec.Encode(&res, w); err != nil {
		http.Error(w, err.Error(), 400)
	}
}
