// Code generated by protoc-gen-go-ttrpc. DO NOT EDIT.
// source: cdh.proto
package __

import (
	context "context"
	ttrpc "github.com/containerd/ttrpc"
)

type SealedSecretServiceService interface {
	UnsealSecret(context.Context, *UnsealSecretInput) (*UnsealSecretOutput, error)
}

func RegisterSealedSecretServiceService(srv *ttrpc.Server, svc SealedSecretServiceService) {
	srv.RegisterService("cdh.SealedSecretService", &ttrpc.ServiceDesc{
		Methods: map[string]ttrpc.Method{
			"UnsealSecret": func(ctx context.Context, unmarshal func(interface{}) error) (interface{}, error) {
				var req UnsealSecretInput
				if err := unmarshal(&req); err != nil {
					return nil, err
				}
				return svc.UnsealSecret(ctx, &req)
			},
		},
	})
}

type sealedsecretserviceClient struct {
	client *ttrpc.Client
}

func NewSealedSecretServiceClient(client *ttrpc.Client) SealedSecretServiceService {
	return &sealedsecretserviceClient{
		client: client,
	}
}

func (c *sealedsecretserviceClient) UnsealSecret(ctx context.Context, req *UnsealSecretInput) (*UnsealSecretOutput, error) {
	var resp UnsealSecretOutput
	if err := c.client.Call(ctx, "cdh.SealedSecretService", "UnsealSecret", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetResourceServiceService interface {
	GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error)
}

func RegisterGetResourceServiceService(srv *ttrpc.Server, svc GetResourceServiceService) {
	srv.RegisterService("cdh.GetResourceService", &ttrpc.ServiceDesc{
		Methods: map[string]ttrpc.Method{
			"GetResource": func(ctx context.Context, unmarshal func(interface{}) error) (interface{}, error) {
				var req GetResourceRequest
				if err := unmarshal(&req); err != nil {
					return nil, err
				}
				return svc.GetResource(ctx, &req)
			},
		},
	})
}

type getresourceserviceClient struct {
	client *ttrpc.Client
}

func NewGetResourceServiceClient(client *ttrpc.Client) GetResourceServiceService {
	return &getresourceserviceClient{
		client: client,
	}
}

func (c *getresourceserviceClient) GetResource(ctx context.Context, req *GetResourceRequest) (*GetResourceResponse, error) {
	var resp GetResourceResponse
	if err := c.client.Call(ctx, "cdh.GetResourceService", "GetResource", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type SecureMountServiceService interface {
	SecureMount(context.Context, *SecureMountRequest) (*SecureMountResponse, error)
}

func RegisterSecureMountServiceService(srv *ttrpc.Server, svc SecureMountServiceService) {
	srv.RegisterService("cdh.SecureMountService", &ttrpc.ServiceDesc{
		Methods: map[string]ttrpc.Method{
			"SecureMount": func(ctx context.Context, unmarshal func(interface{}) error) (interface{}, error) {
				var req SecureMountRequest
				if err := unmarshal(&req); err != nil {
					return nil, err
				}
				return svc.SecureMount(ctx, &req)
			},
		},
	})
}

type securemountserviceClient struct {
	client *ttrpc.Client
}

func NewSecureMountServiceClient(client *ttrpc.Client) SecureMountServiceService {
	return &securemountserviceClient{
		client: client,
	}
}

func (c *securemountserviceClient) SecureMount(ctx context.Context, req *SecureMountRequest) (*SecureMountResponse, error) {
	var resp SecureMountResponse
	if err := c.client.Call(ctx, "cdh.SecureMountService", "SecureMount", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}