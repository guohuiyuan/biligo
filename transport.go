package bilibili

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type responseEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Msg     string          `json:"msg"`
	TTL     int             `json:"ttl"`
	Data    json.RawMessage `json:"data"`
	Result  json.RawMessage `json:"result"`
}

// RequestBuilder encapsulates the transport pipeline.
type RequestBuilder struct {
	client   *Client
	endpoint endpoint
	params   url.Values
	form     url.Values
	headers  http.Header
	body     io.Reader
}

func newRequestBuilder(client *Client, endpoint endpoint) *RequestBuilder {
	return &RequestBuilder{
		client:   client,
		endpoint: endpoint,
		params:   make(url.Values),
		form:     make(url.Values),
		headers:  make(http.Header),
	}
}

func (r *RequestBuilder) Param(key, value string) *RequestBuilder {
	r.params.Set(key, value)
	return r
}

func (r *RequestBuilder) ParamInt(key string, value int64) *RequestBuilder {
	r.params.Set(key, strconv.FormatInt(value, 10))
	return r
}

func (r *RequestBuilder) ParamBool(key string, value bool) *RequestBuilder {
	if value {
		r.params.Set(key, "1")
	} else {
		r.params.Set(key, "0")
	}
	return r
}

func (r *RequestBuilder) Header(key, value string) *RequestBuilder {
	r.headers.Set(key, value)
	return r
}

func (r *RequestBuilder) Form(key, value string) *RequestBuilder {
	r.form.Set(key, value)
	return r
}

func (r *RequestBuilder) FormInt(key string, value int64) *RequestBuilder {
	r.form.Set(key, strconv.FormatInt(value, 10))
	return r
}

func (r *RequestBuilder) FormBool(key string, value bool) *RequestBuilder {
	if value {
		r.form.Set(key, "1")
	} else {
		r.form.Set(key, "0")
	}
	return r
}

func (r *RequestBuilder) JSONBody(v any) *RequestBuilder {
	body, _ := json.Marshal(v)
	r.body = bytes.NewReader(body)
	r.headers.Set("Content-Type", "application/json")
	return r
}

func (r *RequestBuilder) Do(ctx context.Context, out any) error {
	credential := r.client.credentialOrDefault()
	if r.endpoint.needLogin {
		if err := credential.EnsureSessData(); err != nil {
			return err
		}
	}
	if r.endpoint.needCSRF {
		if err := credential.EnsureBiliJct(); err != nil {
			return err
		}
		r.params.Set("csrf", credential.BiliJct)
		r.params.Set("csrf_token", credential.BiliJct)
	}
	if r.endpoint.withWBI {
		if err := r.client.wbi.sign(ctx, r.client, r.params); err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, r.endpoint.method, r.endpoint.url(r.client), r.body)
	if err != nil {
		return err
	}
	req.URL.RawQuery = r.params.Encode()
	if r.body == nil && len(r.form) > 0 {
		req.Body = io.NopCloser(strings.NewReader(r.form.Encode()))
		req.ContentLength = int64(len(r.form.Encode()))
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	for key, value := range defaultHeaders {
		req.Header.Set(key, value)
	}
	if r.client.config.UserAgent != "" {
		req.Header.Set("User-Agent", r.client.config.UserAgent)
	}
	for key, values := range r.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{StatusCode: resp.StatusCode, Method: req.Method, URL: req.URL.String()}
	}

	var envelope responseEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}
	if !r.endpoint.ignoreCode && envelope.Code != 0 {
		message := envelope.Message
		if message == "" {
			message = envelope.Msg
		}
		return &APIError{Code: envelope.Code, Message: message, Data: body}
	}

	if out == nil {
		return nil
	}

	payload := envelope.Data
	switch strings.ToLower(r.endpoint.dataField) {
	case "result":
		payload = envelope.Result
	case "", "data":
	}
	if len(payload) == 0 || string(payload) == "null" {
		return nil
	}
	return json.Unmarshal(payload, out)
}
