package etcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcd"
)

type ClientOptions struct {
	Name          string
	Endpoints     []string
	ClientOptions etcd.ClientOptions
}

type clientIn struct {
	di.In

	Conf   contract.ConfigAccessor
	Logger log.Logger

	Options *ClientOptions `optional:"true"`
}

func provideClient(in clientIn) (etcd.Client, error) {
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
	in.Logger.Log("endpoints", in.Options.Endpoints)

	return etcd.NewClient(context.Background(), in.Options.Endpoints, in.Options.ClientOptions)
}
