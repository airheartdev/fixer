package fixer

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client for the Foreign exchange rates and currency conversion API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string
}

// NewClient creates a Client
func NewClient(options ...func(*Client)) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "api.fixer.io",
		},
		userAgent: "fixer/client.go (https://github.com/peterhellberg/fixer)",
	}

	for _, f := range options {
		f(c)
	}

	return c
}

// HTTPClient changes the HTTP client to the provided *http.Client
func HTTPClient(hc *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// BaseURL changes the base URL to the provided rawurl
func BaseURL(rawurl string) func(*Client) {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// UserAgent changes the User-Agent used by the client
func UserAgent(ua string) func(*Client) {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// Base sets the base query variable based on a Currency
func Base(c Currency) url.Values {
	return url.Values{"base": {string(c)}}
}

// Symbols sets the symbols query variable based on the provided currencies
func Symbols(currencies ...Currency) url.Values {
	symbols := []string{}

	for _, c := range currencies {
		symbols = append(symbols, string(c))
	}

	return url.Values{"symbols": {strings.Join(symbols, ",")}}
}

// Latest foreign exchange reference rates
func (c *Client) Latest(ctx context.Context, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, "/latest", c.query(attributes))
}

// Date returns historical rates for any day since 1999
func (c *Client) Date(ctx context.Context, t time.Time, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, "/"+c.date(t), c.query(attributes))
}

func (c *Client) date(t time.Time) string {
	return t.Format("2006-01-02")
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (*Response, error) {
	req, err := c.request(ctx, path, query)
	if err != nil {
		return nil, err
	}

	r, err := c.do(req)
	if err != nil {
		return nil, err
	}

	r.Links = Links{
		"base": c.baseURL.String(),
		"self": req.URL.String(),
	}

	return r, nil
}

func (c *Client) query(attributes []url.Values) url.Values {
	v := url.Values{}

	for _, a := range attributes {
		if base := a.Get("base"); base != "" {
			v.Set("base", base)
		}

		if symbols := a.Get("symbols"); symbols != "" {
			v.Set("symbols", symbols)
		}
	}

	return v
}

func (c *Client) request(ctx context.Context, path string, query url.Values) (*http.Request, error) {
	rawurl := path

	if len(query) > 0 {
		rawurl += "?" + query.Encode()
	}

	rel, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.baseURL.ResolveReference(rel).String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) do(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, 64)
		_ = resp.Body.Close()
	}()

	if err := responseError(resp); err != nil {
		return nil, err
	}

	var r Response

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}