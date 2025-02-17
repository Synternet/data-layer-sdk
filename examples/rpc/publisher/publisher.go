package main

import (
	"context"

	v1 "github.com/synternet/data-layer-sdk/examples/rpc/types/example/v1"
	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/rpc"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

const (
	SourceParam = "src"
)

type Publisher struct {
	*service.Service
	*rpc.ServiceRegistrar
	userService UserService
}

func New(o ...options.Option) (*Publisher, error) {
	ret := &Publisher{
		Service: &service.Service{},
	}

	ret.Service.Configure(o...)

	return ret, nil
}

func (p *Publisher) Start() context.Context {
	v1.RegisterUserServiceServer(p.ServiceRegistrar, &p.userService)
	return p.Service.Start()
}
