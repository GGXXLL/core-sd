package etcd

import (
	"fmt"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
)

type instancerIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  etcdv3.Client
	Options *InstancerOption
}

type InstancerOption struct {
	Prefix string
}

func provideInstancer(in instancerIn) (sd.Instancer, error) {
	if in.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	return etcdv3.NewInstancer(in.Client, in.Options.Prefix, in.Logger)
}
