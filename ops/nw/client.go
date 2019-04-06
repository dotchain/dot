// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package nw

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dotchain/dot/ops"
)

// Client implements the ops.Store interface by making network calls
// to the provided Url.  All other fields of the Client are optional.
type Client struct {
	URL string
	*http.Client
	http.Header
	ContentType string
	Codecs      map[string]Codec
}

func (c *Client) call(ctx context.Context, r *request) ([]ops.Op, error) {
	if deadline, ok := ctx.Deadline(); ok {
		r.Duration = time.Until(deadline)
	}

	client := c.Client
	if client == nil {
		client = &http.Client{}
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
		log.Println("Encoding error (see https://github.com/dotchain/dot/wiki/Gob-error)")
		log.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", c.URL, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	for key := range c.Header {
		req.Header.Add(key, c.Header.Get(key))
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err == nil {
		defer func() { must(resp.Body.Close()) }()
	}

	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err == nil {
			err = errors.New(resp.Status)
		}
		log.Println("client connect error", err)
		return nil, err
	}

	var res response
	err = codec.Decode(&res, resp.Body)
	if err != nil {
		log.Println("Decoding error (see https://github.com/dotchain/dot/wiki/Gob-error)")
		log.Println(err)
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

// Close proxies the Close call over to the url
func (c *Client) Close() {
}

func must(err error) {
	if err != nil {
		log.Println("client unexpected error", err)
	}
}
