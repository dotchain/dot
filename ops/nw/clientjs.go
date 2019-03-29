// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build js,!jsreflect

package nw

import (
	"bytes"
	"context"
	"time"

	"github.com/dotchain/dot/ops"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/xhr"
)

// Client implements the ops.Store interface by making network calls
// to the provided Url.  All other fields of the Client are optional.
type Client struct {
	URL         string
	ContentType string
	Codecs      map[string]Codec
}

func (c *Client) call(ctx context.Context, r *request) ([]ops.Op, error) {
	if deadline, ok := ctx.Deadline(); ok {
		r.Duration = time.Until(deadline)
	}

	contentType := c.ContentType
	if contentType == "" {
		contentType = "application/x-gob"
	}

	codecs := c.Codecs
	if codecs == nil {
		codecs = DefaultCodecs
	}

	codec := codecs[contentType]
	var buf bytes.Buffer
	err := codec.Encode(r, &buf)
	if err != nil {
		return nil, err
	}

	req := xhr.NewRequest("POST", c.URL)
	req.Timeout = int(r.Duration / time.Millisecond)
	req.ResponseType = xhr.ArrayBuffer
	req.SetRequestHeader("Content-Type", contentType)
	if err := req.Send(buf.Bytes()); err != nil {
		return nil, err
	}

	body := js.Global.Get("Uint8Array").New(req.Response).Interface().([]byte)
	var res response
	if err := codec.Decode(&res, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	return res.Ops, res.Error
}

// Append proxies the Append call over to the url
func (c *Client) Append(ctx context.Context, o []ops.Op) error {
	_, err := c.call(ctx, &request{"Append", o, -1, -1, 0})
	return err
}

// GetSince proxies the GetSince call over to the url
func (c *Client) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	return c.call(ctx, &request{"GetSince", nil, version, limit, 0})
}

// Poll proxies the Poll call over to the url
func (c *Client) Poll(ctx context.Context, version int) error {
	_, err := c.call(ctx, &request{"Poll", nil, version, -1, 0})
	return err
}

// Close proxies the Close call over to the url
func (c *Client) Close() {
}
