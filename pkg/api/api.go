// Copyright (c) 2024 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"fmt"
	"net"

	"github.com/containerd/ttrpc"
	cdh "github.com/slackhq/simple-kubernetes-webhook/pkg/api/cdh"
)

const (
	defaultMinTimeout = 60
	CDHSocket         = "/run/confidential-containers/cdh.sock"
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
