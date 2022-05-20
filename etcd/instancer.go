package etcd

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/ggxxll/core-sd/internal"
	"github.com/ggxxll/core-sd/util"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

type instancerIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  etcdv3.Client
	Options *InstancerOption `optional:"true"`
}

type AgentServiceRegistrationInterceptor func(*api.AgentServiceRegistration)

// InstancerOption wraps args of etcdv3.NewInstancer func.
type InstancerOption struct {
	HTTPPrefix string
	GRPCPrefix string
}

func provideInstancer(in instancerIn) (sd.Instancer, error) {
	if in.Options == nil {
		in.Options = &InstancerOption{
			HTTPPrefix: internal.ServiceKey(in.AppName, in.Env, "http"),
			GRPCPrefix: internal.ServiceKey(in.AppName, in.Env, "grpc"),
		}
	}
	ie := util.Instancer{}
	if in.Options.HTTPPrefix != "" {
		i, err := etcdv3.NewInstancer(in.Client, in.Options.HTTPPrefix, in.Logger)
		if err != nil {
			return nil, err
		}
		ie.HTTP = i
	}
	if in.Options.GRPCPrefix != "" {
		i, err := etcdv3.NewInstancer(in.Client, in.Options.GRPCPrefix, in.Logger)
		if err != nil {
			return nil, err
		}
		ie.GRPC = i
	}

	return &ie, nil
}
