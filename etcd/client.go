package etcd

import (
	"context"
	"fmt"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/log"
)

type ClientOptions struct {
	Name          string
	Endpoints     []string
	ClientOptions etcdv3.ClientOptions
}

type clientIn struct {
	di.In

	Conf   contract.ConfigAccessor
	Logger log.Logger

	Options *ClientOptions `optional:"true"`
}

type HTTPClient etcdv3.Client
type GRPCClient etcdv3.Client

func provideClient(in clientIn) (HTTPClient, GRPCClient, error) {
	if in.Options == nil {
		in.Options = &ClientOptions{
			Name: "default",
		}
	}

	if len(in.Options.Endpoints) == 0 {
		if in.Options.Name == "" {
			in.Options.Name = "default"
		}
		in.Options.Endpoints = in.Conf.Strings(fmt.Sprintf("etcd.%s.endpoints", in.Options.Name))
	}
	forHTTP, err := etcdv3.NewClient(context.Background(), in.Options.Endpoints, in.Options.ClientOptions)
	if err != nil {
		return nil, nil, err
	}
	forGRPC, err := etcdv3.NewClient(context.Background(), in.Options.Endpoints, in.Options.ClientOptions)
	if err != nil {
		return nil, nil, err
	}

	return HTTPClient(forHTTP), GRPCClient(forGRPC), nil
}
