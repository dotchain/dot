// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package nw

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"time"

	"github.com/dotchain/dot/log"
	"github.com/dotchain/dot/ops"
)

func (c *Client) request(ctx context.Context, r *request) ([]ops.Op, error) {
	if c.Log == nil {
		c.Log = log.Default()
	}

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
	err := codec.Encode(*r, &buf)
	if err != nil {
		return nil, c.codecError(err)
	}

	// do JS or non-JS specific network call
	body, err := c.call(ctx, r, contentType, buf.Bytes())
	if closer, ok := body.(io.Closer); ok {
		defer func() {
			c.must(closer.Close())
		}()
	}

	if err != nil {
		c.Log.Println(err)
		return nil, err
	}

	var res response
	err = codec.Decode(&res, body)
	if err != nil {
		return nil, c.codecError(err)
	}
	return res.Ops, res.Error
}

// Append proxies the Append call over to the url
func (c *Client) Append(ctx context.Context, o []ops.Op) error {
	_, err := c.request(ctx, &request{"Append", o, -1, -1, 0})
	return err
}

// GetSince proxies the GetSince call over to the url
func (c *Client) GetSince(ctx context.Context, version, limit int) ([]ops.Op, error) {
	return c.request(ctx, &request{"GetSince", nil, version, limit, 0})
}

// Close proxies the Close call over to the url
func (c *Client) Close() {
}

func (c *Client) codecError(err error) error {
	c.Log.Println("Codec error (see https://github.com/dotchain/dot/wiki/Gob-error)")
	c.must(err)
	return err
}

func (c *Client) must(err error) {
	if err != nil {
		c.Log.Fatal("client unexpected error", err)
	}
}

type httpStatusError struct {
	status int
}

func (h httpStatusError) Error() string {
	return "unexpected http status " + strconv.Itoa(h.status)
}
