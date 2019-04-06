// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build js,!jsreflect

package nw

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/dotchain/dot/log"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/xhr"
)

// Client implements the ops.Store interface by making network calls
// to the provided Url.  All other fields of the Client are optional.
type Client struct {
	URL         string
	ContentType string
	Header      map[string]string
	Codecs      map[string]Codec
	log.Log
}

func (c *Client) call(ctx context.Context, r *request, ct string, body []byte) (io.Reader, error) {
	req := xhr.NewRequest("POST", c.URL)
	req.Timeout = int(r.Duration / time.Millisecond)
	req.ResponseType = xhr.ArrayBuffer
	req.SetRequestHeader("Content-Type", ct)
	for key, value := range c.Header {
		req.SetRequestHeader(key, value)
	}

	if err := req.Send(body); err != nil {
		return nil, err
	}

	if req.Status < 200 || req.Status >= 300 {
		return nil, httpStatusError{req.Status}
	}

	body = js.Global.Get("Uint8Array").New(req.Response).Interface().([]byte)
	return bytes.NewReader(body), nil
}
