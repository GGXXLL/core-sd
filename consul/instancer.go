package consul

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
)

type instancerIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  consul.Client
	Options *InstancerOption
}

type InstancerOption struct {
	Service     string
	Tags        []string
	PassingOnly bool
}

func provideInstancer(in instancerIn) (sd.Instancer, error) {
	return consul.NewInstancer(in.Client, in.Logger, in.Options.Service, in.Options.Tags, in.Options.PassingOnly), nil

}