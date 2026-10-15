package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/NovikovRoman/dadata/v2/client/transport/encoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type capturedRequest struct {
	method      string
	path        string
	contentType string
	body        []byte
}

func newTestServer(t *testing.T, status int, response string) (*url.URL, *capturedRequest) {
	t.Helper()

	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.method = r.Method
		captured.path = r.URL.Path
		captured.contentType = r.Header.Get("Content-Type")
		captured.body, _ = io.ReadAll(r.Body)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL + "/")
	require.NoError(t, err)
	return u, captured
}

func TestClient_GetSendsNoBody(t *testing.T) {
	u, captured := newTestServer(t, http.StatusOK, `{"balance": 1.5}`)

	var result struct {
		Balance float64 `json:"balance"`
	}
	err := NewClient(u).Get(context.Background(), "profile/balance", nil, &result)

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, captured.method)
	assert.Empty(t, captured.body)
	assert.Empty(t, captured.contentType)
	assert.Equal(t, 1.5, result.Balance)
}

func TestClient_GetWithRawEncoder(t *testing.T) {
	u, _ := newTestServer(t, http.StatusOK, `{"balance": 1.5}`)

	cli := NewClient(u,
		WithEncoderFactory(encoder.RawEncoderFactory()),
		WithDecoderFactory(encoder.RawDecoderFactory()),
	)
	response := &bytes.Buffer{}
	err := cli.Get(context.Background(), "profile/balance", nil, response)

	require.NoError(t, err)
	assert.Equal(t, `{"balance": 1.5}`, response.String())
}

func TestClient_PostSendsJSONBody(t *testing.T) {
	u, captured := newTestServer(t, http.StatusOK, `{}`)

	var result map[string]interface{}
	err := NewClient(u).Post(context.Background(), "suggest/address", map[string]string{"query": "москва"}, &result)

	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, captured.method)
	assert.Equal(t, "/suggest/address", captured.path)
	assert.Equal(t, "application/json", captured.contentType)
	assert.JSONEq(t, `{"query": "москва"}`, string(captured.body))
}

func TestClient_ResponseErrorContainsBody(t *testing.T) {
	u, _ := newTestServer(t, http.StatusForbidden, `{"detail": "Invalid API key"}`)

	var result map[string]interface{}
	err := NewClient(u).Get(context.Background(), "profile/balance", nil, &result)

	var respErr *ResponseError
	require.True(t, errors.As(err, &respErr))
	assert.Equal(t, http.StatusForbidden, respErr.StatusCode)
	assert.Equal(t, `{"detail": "Invalid API key"}`, string(respErr.Body))
	assert.Contains(t, err.Error(), "Invalid API key")
}

func TestClient_ResponseErrorBodyIsTruncated(t *testing.T) {
	u, _ := newTestServer(t, http.StatusInternalServerError, strings.Repeat("x", maxErrorBodySize*2))

	var result map[string]interface{}
	err := NewClient(u).Get(context.Background(), "profile/balance", nil, &result)

	var respErr *ResponseError
	require.True(t, errors.As(err, &respErr))
	assert.Len(t, respErr.Body, maxErrorBodySize)
}
