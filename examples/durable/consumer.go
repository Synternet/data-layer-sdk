package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/synternet/data-layer-sdk/pkg/options"
	"github.com/synternet/data-layer-sdk/pkg/service"
)

const (
	SourceParam = "src"
)

type Publisher struct {
	*service.Service
}

type MyMessage struct {
	RequestWas []byte `json:"request"`
}

func New(o ...options.Option) (*Publisher, error) {
	ret := &Publisher{
		Service: &service.Service{},
	}

	err := ret.Service.Configure(o...)
	if err != nil {
		return nil, err
	}

	err = ret.AddStream(10, 10240, time.Hour*24, options.Param(ret.Options, SourceParam, ""))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *Publisher) Start() context.Context {
	err := p.subscribe()
	if err != nil {
		p.Fail(err)
		return p.Context
	}

	return p.Service.Start()
}

func (p *Publisher) subscribe() error {
	src := options.Param(p.Options, SourceParam, "")
	if src == "" {
		return errors.New("source subject must not be empty")
	}
	if _, err := p.SubscribeTo(p.handleQuery, src); err != nil {
		return err
	}

	return nil
}

var cnt = 0

func (p *Publisher) handleQuery(nmsg service.Message) {
	fmt.Println("---", nmsg.Subject())
	for k, v := range nmsg.Header() {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
	}
	fmt.Println()

	fmt.Println(string(nmsg.Data()))
	fmt.Println()
	fmt.Println()
}
