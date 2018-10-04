// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"context"
	"github.com/dotchain/dot/ops"
	"net/http"
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
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		w.Header().Add("Content-Type", ct)
	}
	httpError := func(status string, code int) {
		http.Error(w, status, code)
	}
	h.HandleLambda(context.Background(), ct, httpError, r.Body, w)
}
