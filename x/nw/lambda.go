// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"context"
	"errors"
	"io"
	"time"
)

// HandleLambda supports a AWS-lambda like invocation
func (h *Handler) HandleLambda(ctx context.Context, ct string, httpError func(status string, code int), r io.Reader, w io.Writer) {
	codecs := h.Codecs
	if codecs == nil {
		codecs = DefaultCodecs
	}

	codec := codecs[ct]
	if codec == nil {
		httpError("Invalid content-type", 400)
		return
	}

	var req request
	err := codec.Decode(&req, r)
	if err != nil {
		httpError(err.Error(), 400)
		return
	}

	duration := 30 * time.Second
	if req.Duration != 0 {
		duration = req.Duration
	}
	ctx, done := context.WithTimeout(ctx, duration)
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

	if err := codec.Encode(&res, w); err != nil {
		httpError(err.Error(), 400)
	}
}
