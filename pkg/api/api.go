// Copyright (c) 2024 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/containerd/ttrpc"
	cdh "github.com/slackhq/simple-kubernetes-webhook/pkg/api/cdh"
)

const (
	defaultMinTimeout  = 60
	CDHSocket          = "/run/confidential-containers/cdh.sock"
	sealedSecretPrefix = "sealed."
)

type cdhClient struct {
	conn               net.Conn
	sealedSecretClient cdh.SealedSecretServiceService
}

func CreateCDHClient() (*cdhClient, error) {
	conn, err := net.Dial("unix", CDHSocket)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cdh socket: %w", err)
	}

	ttrpcClient := ttrpc.NewClient(conn)
	client := cdh.NewSealedSecretServiceClient(ttrpcClient)

	c := &cdhClient{
		conn:               conn,
		sealedSecretClient: client,
	}
	return c, nil
}

func (c *cdhClient) Close() error {
	return c.conn.Close()
}

func (c *cdhClient) UnsealSecret(ctx context.Context, secret string) (string, error) {
	input := cdh.UnsealSecretInput{Secret: []byte(secret)}
	output, err := c.sealedSecretClient.UnsealSecret(ctx, &input)
	if err != nil {
		return "", fmt.Errorf("failed to unseal secret: %w", err)
	}

	return string(output.GetPlaintext()[:]), nil
}

func (c *cdhClient) UnsealEnv(ctx context.Context, env string) (string, error) {
	unsealedValue, err := c.UnsealSecret(ctx, env)
	if err != nil {
		return "", fmt.Errorf("failed to unseal secret from env: %w", err)
	}
	return unsealedValue, nil
}

func HasSealedSecretsPrefix(value string) bool {
	return strings.HasPrefix(value, sealedSecretPrefix)
}
