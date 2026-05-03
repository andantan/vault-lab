package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type Client struct {
	url        string
	httpClient *http.Client
	nextID     uint64
}

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type Response[T any] struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      uint64    `json:"id"`
	Result  T         `json:"result"`
	Error   *RPCError `json:"error,omitempty"`
}

type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	if len(e.Data) == 0 {
		return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
	}

	return fmt.Sprintf("rpc error %d: %s: %s", e.Code, e.Message, string(e.Data))
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
}

func (c *Client) Call(ctx context.Context, method string, params any, result any) error {
	id := atomic.AddUint64(&c.nextID, 1)
	reqBody := Request{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send rpc request: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read rpc response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return fmt.Errorf("http status %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var rpcResp Response[json.RawMessage]
	if err = json.Unmarshal(respBytes, &rpcResp); err != nil {
		return fmt.Errorf("decode rpc response: %w: %s", err, string(respBytes))
	}

	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	if result == nil {
		return nil
	}

	if err = json.Unmarshal(rpcResp.Result, result); err != nil {
		return fmt.Errorf("decode rpc result for %s: %w: %s", method, err, string(rpcResp.Result))
	}

	return nil
}
