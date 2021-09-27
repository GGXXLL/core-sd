package zk

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
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
	Options *InstancerOption
}

// InstancerOption wraps args of zk.NewInstancer func.
type InstancerOption struct {
	Path string
}

func provideInstancer(in instancerIn) (sd.Instancer, error) {
	return zk.NewInstancer(in.Client, in.Options.Path, in.Logger)
}
