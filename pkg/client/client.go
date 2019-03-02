/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/gravitational/fakeiot/pkg/metric"
	"github.com/gravitational/roundtrip"
	"github.com/gravitational/trace"
)

// New returns a new client
func New(cfg Config) (*Client, error) {
	if err := cfg.Check(); err != nil {
		return nil, trace.Wrap(err)
	}
	var client *http.Client
	if cfg.CACert != nil {
		pool := x509.NewCertPool()
		pool.AddCert(cfg.CACert)
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: pool,
				},
			},
		}
	} else {
		client = &http.Client{}
	}
	clt, err := roundtrip.NewClient(
		cfg.URL.String(), "",
		roundtrip.HTTPClient(client),
		roundtrip.BearerAuth(cfg.BearerToken))
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return &Client{Config: cfg, Client: clt}, nil
}

// Config is a client config
type Config struct {
	// BearerToken is a bearer token
	// that is used to authenticate requests
	BearerToken string
	// URL is a URL to send the data to
	URL *url.URL
	// CACert is a certificate authority certificate
	CACert *x509.Certificate
}

// Check checks configuration
func (c *Config) Check() error {
	if c.URL == nil {
		return trace.BadParameter("missing parameter URL")
	}
	if c.URL.Scheme != "https" {
		return trace.BadParameter("only HTTPS scheme is supported")
	}
	// NOTE: client does not check bearer token
	// to make it possible to test bogus requests
	return nil
}

// Client is an HTTP client sending the data
// to the IOT metric server
type Client struct {
	// Config is a client configuration
	Config
	// Client is a roundtrip client
	Client *roundtrip.Client
}

// Send sends the metric to the IOT server
func (c *Client) Send(ctx context.Context, m metric.Metric) error {
	re, err := c.Client.PostJSON(ctx, c.Client.Endpoint("metrics"), m)
	if err != nil {
		return trace.ConvertSystemError(err)
	}
	contentType := re.Headers().Get(contentTypeHeader)
	if contentType == "" {
		return trace.BadParameter("server did not respond with %v header", contentTypeHeader)
	}
	if contentType != contentTypeJSON {
		return trace.BadParameter("received unexpected %v: %q, expected %q", contentTypeHeader, contentType, contentTypeJSON)
	}
	return trace.ReadError(re.Code(), re.Bytes())
}

// SendCorruptedData sends corrupted non-json data to the server
func (c *Client) SendCorruptedData(ctx context.Context) error {
	re, err := c.Client.PostForm(ctx, c.Client.Endpoint("metrics"), url.Values{"bad": []string{"format"}})
	if err != nil {
		return trace.ConvertSystemError(err)
	}
	contentType := re.Headers().Get(contentTypeHeader)
	if contentType == "" {
		return trace.BadParameter("server did not respond with %v header", contentTypeHeader)
	}
	if contentType != contentTypeJSON {
		return trace.BadParameter("received unexpected %v: %q, expected %q", contentTypeHeader, contentType, contentTypeJSON)
	}
	return trace.ReadError(re.Code(), re.Bytes())
}

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)
