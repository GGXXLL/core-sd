package zk

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/ggxxll/core-sd/internal"
	"github.com/ggxxll/core-sd/util"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/zk"
	"github.com/go-kit/log"
)

type instancerIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  zk.Client
	Options *InstancerOption `optional:"true"`
}

// InstancerOption wraps args of zk.NewInstancer func.
type InstancerOption struct {
	HTTPPath string
	GRPCPath string
}

func provideInstancer(in instancerIn) (sd.Instancer, error) {
	if in.Options == nil {
		in.Options = &InstancerOption{
			HTTPPath: "/" + internal.ServiceKey(in.AppName, in.Env, "http"),
			GRPCPath: "/" + internal.ServiceKey(in.AppName, in.Env, "grpc"),
		}
	}
	ie := util.Instancer{}
	if in.Options.HTTPPath != "" {
		i, err := zk.NewInstancer(in.Client, in.Options.HTTPPath, in.Logger)
		if err != nil {
			return nil, err
		}
		ie.HTTP = i
	}
	if in.Options.GRPCPath != "" {
		i, err := zk.NewInstancer(in.Client, in.Options.GRPCPath, in.Logger)
		if err != nil {
			return nil, err
		}
		ie.GRPC = i
	}

	return &ie, nil
}
