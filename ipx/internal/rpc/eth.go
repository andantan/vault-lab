package rpc

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) ChainID(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, "eth_chainId", []any{}, &result)
	return result, err
}

func (c *Client) BlockNumber(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, "eth_blockNumber", []any{}, &result)
	return result, err
}

func (c *Client) GetBalance(ctx context.Context, address string, block string) (string, error) {
	if block == "" {
		block = "latest"
	}

	var result string
	err := c.Call(ctx, "eth_getBalance", []any{address, block}, &result)
	return result, err
}

func (c *Client) GetCode(ctx context.Context, address string, block string) (string, error) {
	if block == "" {
		block = "latest"
	}

	var result string
	err := c.Call(ctx, "eth_getCode", []any{address, block}, &result)
	return result, err
}

func (c *Client) CallContract(ctx context.Context, params any, block string) (string, error) {
	if block == "" {
		block = "latest"
	}

	var result string
	err := c.Call(ctx, "eth_call", []any{params, block}, &result)
	return result, err
}

func (c *Client) EstimateGas(ctx context.Context, params any) (string, error) {
	var result string
	err := c.Call(ctx, "eth_estimateGas", []any{params}, &result)
	return result, err
}

func (c *Client) GetTransactionCount(ctx context.Context, address string, block string) (string, error) {
	if block == "" {
		block = "pending"
	}

	var result string
	err := c.Call(ctx, "eth_getTransactionCount", []any{address, block}, &result)
	return result, err
}

func (c *Client) MaxPriorityFeePerGas(ctx context.Context) (string, error) {
	var result string
	err := c.Call(ctx, "eth_maxPriorityFeePerGas", []any{}, &result)
	return result, err
}

func (c *Client) SendRawTransaction(ctx context.Context, rawTx string) (string, error) {
	var result string
	err := c.Call(ctx, "eth_sendRawTransaction", []any{rawTx}, &result)
	return result, err
}

func (c *Client) BlockByNumber(ctx context.Context, block string) (map[string]any, error) {
	if block == "" {
		block = "latest"
	}
	var result map[string]any
	err := c.Call(ctx, "eth_getBlockByNumber", []any{block, false}, &result)
	return result, err
}

func (c *Client) TransactionReceipt(ctx context.Context, txHash string) (map[string]any, error) {
	var result map[string]any
	err := c.Call(ctx, "eth_getTransactionReceipt", []any{txHash}, &result)
	return result, err
}

func (c *Client) WaitForReceipt(ctx context.Context, txHash string, timeout time.Duration) (map[string]any, error) {
	deadline := time.Now().Add(timeout)

	for {
		receipt, err := c.TransactionReceipt(ctx, txHash)
		if err == nil && receipt != nil {
			return receipt, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for receipt: %s", txHash)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
