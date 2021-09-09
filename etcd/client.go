package etcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
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

func provideClient(in clientIn) (etcdv3.Client, error) {
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

	for i, end := range in.Options.Endpoints {
		if strings.HasPrefix(end, "http") {
			continue
		}
		in.Options.Endpoints[i] = "http://" + end
	}

	return etcdv3.NewClient(context.Background(), in.Options.Endpoints, in.Options.ClientOptions)
}
