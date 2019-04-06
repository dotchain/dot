// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package nw

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/dotchain/dot/log"
)

// Client implements the ops.Store interface by making network calls
// to the provided Url.  All other fields of the Client are optional.
type Client struct {
	URL         string
	Header      map[string]string
	ContentType string
	Codecs      map[string]Codec
	log.Log

	*http.Client
}

func (c *Client) call(ctx context.Context, r *request, ct string, body []byte) (io.ReadCloser, error) {
	client := c.Client
	if client == nil {
		client = &http.Client{}
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", ct)
	for key, value := range c.Header {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req.WithContext(ctx))
	var resultbody io.ReadCloser
	if resp != nil {
		resultbody = resp.Body
	}

	if err == nil && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		err = httpStatusError{resp.StatusCode}
	}

	return resultbody, err
}
